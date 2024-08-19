package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
)

type assetImpl struct {
	Db *sql.DB
}

func NewAssetImplementation(db *sql.DB) repository.AssetRepository {
	return assetImpl{
		Db: db,
	}
}

func (a assetImpl) Create(asset models.Asset) (models.Asset, error) {
	// generate an ID for the user if it doesn't have one
	if asset.Id == "" {
		asset.Id = uuid.NewString()
	}

	stmt := `
INSERT INTO assets
(
	id,
	token_id,
	"name",
	symbol,
	category_id, 
	blockchain, 
	creator_id, 
	signatories,
	metadata_url, 
	urls, 
	legal_docs, 
	issuance_docs, 
	description, 
	total_supply, 
	serial_number, 
	properties, 
	status, 
	executed_at, 
	expires_at
)
VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
	RETURNING 
	    id,
		status,
	   created_at, 
	   updated_at;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if asset.Properties == nil {
		asset.Properties = make(models.AssetProperties, 0)
	}

	properties, err := json.Marshal(&asset.Properties)
	if err != nil {
		return models.Asset{}, err
	}

	if asset.URLs == nil {
		asset.URLs = make([]string, 0)
	}
	urls, err := json.Marshal(&asset.URLs)
	if err != nil {
		return models.Asset{}, err
	}
	if asset.LegalDocumentURLs == nil {
		asset.LegalDocumentURLs = make([]string, 0)
	}
	legalDocs, err := json.Marshal(asset.LegalDocumentURLs)
	if err != nil {
		return models.Asset{}, err
	}

	if asset.IssuanceDocumentURLs == nil {
		asset.IssuanceDocumentURLs = make([]string, 0)
	}
	issuanceDocs, err := json.Marshal(&asset.IssuanceDocumentURLs)
	if err != nil {
		return models.Asset{}, err
	}
	if asset.Signatories == nil {
		asset.Signatories = make([]string, 0)
	}
	signatories, err := json.Marshal(&asset.Signatories)
	if err != nil {
		return models.Asset{}, err
	}
	newAsset := asset

	err = a.Db.QueryRowContext(ctx, stmt,
		asset.Id,
		asset.TokenId,
		asset.Name,
		asset.Symbol,
		asset.CategoryId,
		asset.Blockchain,
		asset.CreatorUserID,
		signatories,
		asset.MetadataUrl,
		string(urls),
		string(legalDocs),
		string(issuanceDocs),
		asset.Description,
		asset.TotalSupply,
		asset.SerialNumber,
		properties,
		models.UNVERIFIED_ASSET_STATUS,
		asset.ExecutionDate,
		asset.ExpirationDate,
	).Scan(
		&newAsset.Id,
		&newAsset.Status,
		&newAsset.CreatedAt,
		&newAsset.UpdatedAt,
	)
	if err != nil {
		// TODO: handle duplicate error if any field is duplicate
		switch {
		case strings.Contains(err.Error(), duplicateAddress):
			return models.Asset{}, repository.ErrDuplicateDetails

		default:
			return models.Asset{}, err
		}
	}

	return newAsset, err
}

func (a assetImpl) calculateFloorPrice(assetId string, period int) (float64, error) {
	// this query is supposed to go through the list of prices items in activities
	var floorPrice float64
	query := `
	SELECT COALESCE(MIN(price), 0)
	FROM listings
	WHERE asset_id = $1
	AND is_active = true
	AND is_cancelled = false
	AND created_at >= NOW() - INTERVAL '1 day' * $2
		LIMIT 1	
	`
	// 	query := `
	// SELECT COALESCE(MIN(price), 0)
	// 	FROM activities v

	// 	  JOIN assets a ON a.id = v.asset_id
	// 	WHERE  a.id =$1 AND AGE(v.occurred_at) <= INTERVAL '1 day' * $2
	// 	LIMIT 1
	// 	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := a.Db.QueryRowContext(ctx, query, assetId, period).Scan(&floorPrice)
	if err != nil {
		return 0, err
	}

	return floorPrice, nil
}

func (a assetImpl) GetById(id string) (models.Asset, error) {
	stmt := `
	SELECT 
		a.id,
	   token_id,
	   a."name", 
	   symbol,
	    category_id,
	   blockchain, 
	   creator_id,
	   metadata_url, 
	   urls, 
	   legal_docs, 
	   issuance_docs, 
	   signatories,
	   description, 
	   total_supply, 
	   serial_number,
	   properties,
	   status,
	   executed_at, 
	   expires_at, 
	   a.created_at, 
	   a.updated_at,
	   c."name" ,
	   a.is_verified,
		a.is_minted,
		a.is_rejected
FROM assets AS a
	   LEFT JOIN categories AS c ON a.category_id = c.id 
		JOIN users AS u ON a.creator_id = u.id
	WHERE a.id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	properties := ""
	urls := ""
	legalDocs := ""
	issuanceDocs := ""
	signatories := ""

	fetchedAssets := models.Asset{}

	err := a.Db.QueryRowContext(ctx, stmt, id).Scan(
		&fetchedAssets.Id,
		&fetchedAssets.TokenId,
		&fetchedAssets.Name,
		&fetchedAssets.Symbol,
		&fetchedAssets.CategoryId,
		&fetchedAssets.Blockchain,
		&fetchedAssets.CreatorUserID,
		&fetchedAssets.MetadataUrl,
		&urls,
		&legalDocs,
		&issuanceDocs,
		&signatories,

		&fetchedAssets.Description,
		&fetchedAssets.TotalSupply,
		&fetchedAssets.SerialNumber,
		&properties,
		&fetchedAssets.Status,
		&fetchedAssets.ExecutionDate,
		&fetchedAssets.ExpirationDate,

		&fetchedAssets.CreatedAt,
		&fetchedAssets.UpdatedAt,
		&fetchedAssets.CategoryName,
		&fetchedAssets.IsVerified,
		&fetchedAssets.IsMinted,
		&fetchedAssets.IsRejected,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Asset{}, repository.ErrRecordNotFound

		default:
			return models.Asset{}, err
		}
	}

	err = json.NewDecoder(strings.NewReader(properties)).Decode(&fetchedAssets.Properties)
	if err != nil {
		return models.Asset{}, err
	}
	err = json.NewDecoder(strings.NewReader(urls)).Decode(&fetchedAssets.URLs)
	if err != nil {
		return models.Asset{}, err
	}
	err = json.NewDecoder(strings.NewReader(legalDocs)).Decode(&fetchedAssets.LegalDocumentURLs)
	if err != nil {
		return models.Asset{}, err
	}
	err = json.NewDecoder(strings.NewReader(issuanceDocs)).Decode(&fetchedAssets.IssuanceDocumentURLs)
	if err != nil {
		return models.Asset{}, err
	}
	err = json.NewDecoder(strings.NewReader(signatories)).Decode(&fetchedAssets.Signatories)
	if err != nil {
		return models.Asset{}, err
	}

	favourites, err := a.CountFavourites(fetchedAssets.Id)
	if err != nil {
		return models.Asset{}, err
	}
	fetchedAssets.Favourites = favourites

	views, err := a.CountViews(fetchedAssets.Id)
	if err != nil {
		return models.Asset{}, err
	}
	fetchedAssets.Views = views

	return fetchedAssets, nil
}
func (a assetImpl) List(filter models.Filter) ([]models.Asset, error) {
	stmt := `
	SELECT 
		a.id,
	   a.token_id,
	   a."name", 
	   symbol,
	    a.category_id,
	    a.blockchain, 
	    a.creator_id,
	   a. metadata_url, 
	    a.urls, 
	    a.legal_docs, 
	    a.issuance_docs, 
	    a.signatories,
	    a.description, 
	    a.total_supply, 
	    a.serial_number,
	    a.properties,
	    a.status,
	    a.executed_at, 
	    a.expires_at, 
	   a.created_at, 
	   a.updated_at,
	   c."name" ,
	   a.is_verified,
		a.is_minted,
		a.is_rejected
  FROM assets AS a
	   LEFT JOIN categories AS c ON a.category_id = c.id 
		JOIN users AS u ON a.creator_id = u.id
	WHERE 
	  	-- search by item name
	  	a.name ILIKE '%' || $1 || '%'

	-- reverse chronological order
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	assets := []models.Asset{}

	rows, err := a.Db.QueryContext(
		ctx,
		stmt,
		filter.Search,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return []models.Asset{}, err
	}
	defer rows.Close()

	for rows.Next() {
		fetchedAsset := models.Asset{}

		properties := ""
		urls := ""
		legalDocs := ""
		issuanceDocs := ""
		signatories := ""

		rows.Scan(
			&fetchedAsset.Id,
			&fetchedAsset.TokenId,
			&fetchedAsset.Name,
			&fetchedAsset.Symbol,
			&fetchedAsset.CategoryId,
			&fetchedAsset.Blockchain,
			&fetchedAsset.CreatorUserID,
			&fetchedAsset.MetadataUrl,
			&urls,
			&legalDocs,
			&issuanceDocs,
			&signatories,
			&fetchedAsset.Description,
			&fetchedAsset.TotalSupply,
			&fetchedAsset.SerialNumber,
			&properties,
			&fetchedAsset.Status,
			&fetchedAsset.ExecutionDate,
			&fetchedAsset.ExpirationDate,

			&fetchedAsset.CreatedAt,
			&fetchedAsset.UpdatedAt,
			&fetchedAsset.CategoryName,
			&fetchedAsset.IsVerified,
			&fetchedAsset.IsMinted,
			&fetchedAsset.IsRejected,
		)

		err = json.NewDecoder(strings.NewReader(properties)).Decode(&fetchedAsset.Properties)
		if err != nil {
			return []models.Asset{}, err
		}

		err = json.NewDecoder(strings.NewReader(urls)).Decode(&fetchedAsset.URLs)
		if err != nil {
			return []models.Asset{}, err
		}
		err = json.NewDecoder(strings.NewReader(legalDocs)).Decode(&fetchedAsset.LegalDocumentURLs)
		if err != nil {
			return []models.Asset{}, err
		}
		err = json.NewDecoder(strings.NewReader(issuanceDocs)).Decode(&fetchedAsset.IssuanceDocumentURLs)
		if err != nil {
			return []models.Asset{}, err
		}
		err = json.NewDecoder(strings.NewReader(signatories)).Decode(&fetchedAsset.Signatories)
		if err != nil {
			return []models.Asset{}, err
		}

		assets = append(assets, fetchedAsset)
	}

	if err = rows.Err(); err != nil {
		return []models.Asset{}, err
	}

	return assets, nil
}
func (a assetImpl) ListRecommended(asset models.Asset, filter models.Filter) ([]models.AssetStat, int, error) {
	stmt := `
	WITH asset_sales AS (
    SELECT
        a.id AS asset_id,
        a.name,
        a.symbol,
        COUNT(act.id) AS sale_activity_count,
        COALESCE(SUM(act.price), 0) AS total_sales_volume
    FROM
        assets a
    LEFT JOIN
        activities act ON a.id = act.asset_id
    WHERE
        act.action = 'sale' AND
        act.created_at >= NOW() - INTERVAL '1 day' * $3
    GROUP BY
        a.id
),
asset_holders AS (
    SELECT
        uh.asset_id,
        COUNT(DISTINCT uh.user_id) AS holder_count
    FROM
        asset_owners uh
    GROUP BY
        uh.asset_id
)
SELECT
    a.id,
    a.name,
    a.symbol,
    a.category_id,
    c."name",
    a.blockchain,
    a.urls,
    a.status,
	a.creator_id,
	u.username,
    COALESCE(asl.total_sales_volume, 0) AS total_sales_volume,
    COALESCE(asl.sale_activity_count, 0) AS sale_activity_count,
    COALESCE(ah.holder_count, 0) AS holder_count,
	(a.is_verified AND a.is_minted) AS is_verified_and_minted
FROM
    assets a
JOIN users AS u ON a.creator_id = u.id
LEFT JOIN
    asset_sales asl ON a.id = asl.asset_id
LEFT JOIN
    asset_holders ah ON a.id = ah.asset_id
Left JOIN categories AS c ON a.category_id = c.id
    WHERE ((a.creator_id  = $1 and a.blockchain=$2 and ($3::int IS NULL OR a.category_id = $3::int))
OR (a.blockchain=$2 and ($3::int IS NULL OR a.category_id = $3::int))
OR ( $3::int IS NULL OR a.category_id = $3::int))
AND a.id <> $4
LIMIT $5 OFFSET $6;
	`

	var totalCount int
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	assets := []models.AssetStat{}

	rows, err := a.Db.QueryContext(
		ctx,
		stmt,
		asset.CreatorUserID,
		asset.Blockchain,
		asset.CategoryId,
		asset.Id,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return []models.AssetStat{}, totalCount, err
	}
	defer rows.Close()

	totalCountQuery := `
	WITH asset_sales AS (
    SELECT
        a.id AS asset_id,
        a.name,
        a.symbol,
        COUNT(act.id) AS sale_activity_count,
        COALESCE(SUM(act.price), 0) AS total_sales_volume
    FROM
        assets a
    LEFT JOIN
        activities act ON a.id = act.asset_id
    WHERE
        act.action = 'sale' AND
        act.created_at >= NOW() - INTERVAL '1 day' * $3
    GROUP BY
        a.id
),
asset_holders AS (
    SELECT
        uh.asset_id,
        COUNT(DISTINCT uh.user_id) AS holder_count
    FROM
        asset_owners uh
    GROUP BY
        uh.asset_id
)
SELECT COUNT(*)
FROM
    assets a
JOIN users AS u ON a.creator_id = u.id
LEFT JOIN
    asset_sales asl ON a.id = asl.asset_id
LEFT JOIN
    asset_holders ah ON a.id = ah.asset_id
Left JOIN categories AS c ON a.category_id = c.id
    WHERE ((a.creator_id  = $1 and a.blockchain=$2 and ($3::int IS NULL OR a.category_id = $3::int))
OR (a.blockchain=$2 and ($3::int IS NULL OR a.category_id = $3::int))
OR ( $3::int IS NULL OR a.category_id = $3::int))
AND a.id <> $4

 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = a.Db.QueryRowContext(ctx, totalCountQuery,
		asset.CreatorUserID,
		asset.Blockchain,
		asset.CategoryId,
		asset.Id,
	).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		fetchedAsset := models.AssetStat{}

		urls := ""

		err := rows.Scan(
			&fetchedAsset.AssetId,
			&fetchedAsset.Name,
			&fetchedAsset.Symbol,
			&fetchedAsset.CategoryId,
			&fetchedAsset.CategoryName,
			&fetchedAsset.Blockchain,
			&urls,
			&fetchedAsset.Status,
			&fetchedAsset.CreatorId,
			&fetchedAsset.CreatorUsername,
			&fetchedAsset.TotalVolume,
			&fetchedAsset.ActivityCount,
			&fetchedAsset.Owners,
			&fetchedAsset.IsVerified,
		)

		if err != nil {
			log.Printf("err, retrieving asset")
			continue
		}

		err = json.NewDecoder(strings.NewReader(urls)).Decode(&fetchedAsset.URLs)
		if err != nil {
			return nil, totalCount, err
		}

		fetchedAsset.FloorPrice, err = a.calculateFloorPrice(fetchedAsset.AssetId, int(filter.Range))
		if err != nil {
			return nil, totalCount, err
		}

		assets = append(assets, fetchedAsset)
	}

	if err = rows.Err(); err != nil {
		return []models.AssetStat{}, totalCount, err
	}

	return assets, totalCount, nil
}
func (a assetImpl) Update(asset models.Asset) (models.Asset, error) {
	stmt := `
	UPDATE assets
	 	SET    
			properties = $2,
			token_id = $3,
			urls = $4,
			legal_docs = $5,
			issuance_docs = $6,
			signatories = $7,
			category_id = $8,
			description = $9,
			metadata_url = $10,
			"name" = $11,
			blockchain = $12,
			total_supply = $13,
			updated_at = now()
			WHERE id = $1
	RETURNING 
		id, 
	   updated_at;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	updatedAsset := asset
	properties, err := json.Marshal(asset.Properties)
	if err != nil {
		return updatedAsset, err
	}
	urls, err := json.Marshal(asset.URLs)
	if err != nil {
		return updatedAsset, err
	}
	legalDocs, err := json.Marshal(asset.LegalDocumentURLs)
	if err != nil {
		return updatedAsset, err
	}
	issuanceDocs, err := json.Marshal(asset.LegalDocumentURLs)
	if err != nil {
		return updatedAsset, err
	}
	signatories, err := json.Marshal(asset.LegalDocumentURLs)
	if err != nil {
		return updatedAsset, err
	}

	err = a.Db.QueryRowContext(ctx, stmt,
		asset.Id,
		properties,
		asset.TokenId,
		urls,
		legalDocs,
		issuanceDocs,
		signatories,
		asset.CategoryId,
		asset.Description,
		asset.MetadataUrl,
		asset.Name,
		asset.Blockchain,
		asset.TotalSupply,
	).Scan(
		&updatedAsset.Id,
		&updatedAsset.UpdatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), duplicateAddress):
			return models.Asset{}, repository.ErrDuplicateDetails

		default:
			return models.Asset{}, err
		}
	}

	return updatedAsset, nil

}

func (a assetImpl) GetListings(assetId string, filter models.Filter) ([]models.Listing, error) {
	stmt := `
	SELECT 
	    l.id,
        l."type",
        l.user_id, 
        l.asset_id,
		a."name",
        a.symbol ,
        l.price,
        l.min_invest_amount, 
        l.max_invest_amount, 
        l.min__raise_amount, 
        l.max_raise_amount,
        l.currency, 
        l.quantity, 
        l.start_date, 
        l.end_date, 
        l.is_active,
        l.is_cancelled, 
        l.created_at, 
        l.updated_at     
	FROM listings AS l 
	JOIN assets a ON a.id = l.asset_id 
	WHERE 
	    l.asset_id = $1 AND 
		(
			(l.is_active = true) 
			OR (l.start_date > NOW() AND l.is_cancelled = false)  -- Adjusted condition
		)
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3;

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	listings := []models.Listing{}

	rows, _ := a.Db.QueryContext(ctx, stmt, assetId, filter.Limit, filter.Offset())
	for rows.Next() {
		currency := ""
		listing := models.Listing{}

		rows.Scan(
			&listing.Id,
			&listing.Type,
			&listing.UserId,
			&listing.AssetId,
			&listing.AssetName,
			&listing.AssetSymbol,
			&listing.PriceUSD,
			&listing.MinInvestAmount,
			&listing.MaxInvestAmount,
			&listing.MinToRaise,
			&listing.MaxToRaise,
			&currency,
			&listing.Quantity,
			&listing.StartDate,
			&listing.EndDate,
			&listing.IsActive,
			&listing.IsCancelled,
			&listing.CreatedAt,
			&listing.UpdatedAt,
		)

		err := json.NewDecoder(strings.NewReader(currency)).Decode(&listing.Currency)
		if err != nil {
			return nil, err
		}

		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listings, nil
}

func (a assetImpl) AddFavourite(userId string, assetId string) error {
	stmt := `
	INSERT INTO favourites(user_id, asset_id) 
	VALUES ($1, $2);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := a.Db.ExecContext(
		ctx,
		stmt,
		userId,
		assetId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a assetImpl) RemoveFavourite(userId string, assetId string) error {
	stmt := `
	DELETE FROM favourites
	WHERE 
	    user_id = $1
	    AND asset_id = $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := a.Db.ExecContext(
		ctx,
		stmt,
		userId,
		assetId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a assetImpl) IsFavourite(userId string, assetId string) (bool, error) {
	stmt := `
	SELECT EXISTS(
	    SELECT true
	
		FROM favourites
		
		WHERE 
			user_id = $1
			AND asset_id = $2
	);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var isFavourite bool
	err := a.Db.QueryRowContext(
		ctx,
		stmt,
		userId,
		assetId,
	).Scan(&isFavourite)

	if err != nil {
		return false, err
	}

	return isFavourite, nil
}
func (a assetImpl) IsListed(userId string, assetId string) (bool, error) {
	stmt := `
	SELECT EXISTS(
    SELECT 1
    FROM listings
    WHERE 
        user_id = $1
        AND asset_id = $2
        AND is_active = TRUE
        AND is_cancelled = FALSE
);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var isListed bool
	err := a.Db.QueryRowContext(
		ctx,
		stmt,
		userId,
		assetId,
	).Scan(&isListed)

	if err != nil {
		return false, err
	}

	return isListed, nil
}

func (a assetImpl) CountFavourites(assetId string) (int64, error) {
	stmt := `
	SELECT count(*)

	FROM favourites
	
	WHERE asset_id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var favourites int64
	err := a.Db.QueryRowContext(
		ctx,
		stmt,
		assetId,
	).Scan(
		&favourites,
	)

	if err != nil {
		return 0, err
	}

	return favourites, nil
}

func (a assetImpl) AddView(userId string, assetId string) error {
	stmt := `
	INSERT INTO asset_views(user_id, asset_id) 
	VALUES ($1, $2);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := a.Db.ExecContext(
		ctx,
		stmt,
		userId,
		assetId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a assetImpl) IsViewed(userId string, assetId string) (bool, error) {
	stmt := `
	SELECT EXISTS(
	    SELECT true
	
		FROM asset_views
		
		WHERE 
			user_id = $1
			AND asset_id = $2
	);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var isViewed bool
	err := a.Db.QueryRowContext(
		ctx,
		stmt,
		userId,
		assetId,
	).Scan(&isViewed)

	if err != nil {
		return false, err
	}

	return isViewed, nil
}

func (a assetImpl) CountViews(assetId string) (int64, error) {
	stmt := `
	SELECT count(*)

	FROM asset_views
	
	WHERE asset_id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var views int64
	err := a.Db.QueryRowContext(
		ctx,
		stmt,
		assetId,
	).Scan(
		&views,
	)

	if err != nil {
		return 0, err
	}

	return views, nil
}

func (a assetImpl) GetApprovedUnmintedAssets(filter models.Filter) ([]models.Asset, error) {
	stmt := `
		SELECT 
		a.id,
	   a.token_id,
	   a."name", 
	   symbol,
	    a.category_id,
	    a.blockchain, 
	    a.creator_id,
	   a. metadata_url, 
	    a.urls, 
	    a.legal_docs, 
	    a.issuance_docs, 
	    a.signatories,
	    a.description, 
	    a.total_supply, 
	    a.serial_number,
	    a.properties,
	    a.status,
	    a.executed_at, 
	    a.expires_at, 
	   a.created_at, 
	   a.updated_at,
	   c."name" 
  FROM assets AS a
	   LEFT JOIN categories AS c ON a.category_id = c.id 
		JOIN users AS u ON a.creator_id = u.id
	WHERE 
		  a.status=$1 AND a.is_minted=false
	-- reverse chronological order
	ORDER BY created_at asc
	LIMIT $2 OFFSET $3;
	
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	assets := []models.Asset{}

	rows, err := a.Db.QueryContext(
		ctx,
		stmt,
		models.VERIFIED_ASSET_STATUS,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return []models.Asset{}, err
	}
	defer rows.Close()

	for rows.Next() {
		fetchedAsset := models.Asset{}

		properties := ""
		urls := ""
		legalDocs := ""
		issuanceDocs := ""
		signatories := ""

		rows.Scan(
			&fetchedAsset.Id,
			&fetchedAsset.TokenId,
			&fetchedAsset.Name,
			&fetchedAsset.Symbol,
			&fetchedAsset.CategoryId,
			&fetchedAsset.Blockchain,
			&fetchedAsset.CreatorUserID,
			&fetchedAsset.MetadataUrl,
			&urls,
			&legalDocs,
			&issuanceDocs,
			&signatories,
			&fetchedAsset.Description,
			&fetchedAsset.TotalSupply,
			&fetchedAsset.SerialNumber,
			&properties,
			&fetchedAsset.Status,
			&fetchedAsset.ExecutionDate,
			&fetchedAsset.ExpirationDate,

			&fetchedAsset.CreatedAt,
			&fetchedAsset.UpdatedAt,
			&fetchedAsset.CategoryName,
		)

		err = json.NewDecoder(strings.NewReader(properties)).Decode(&fetchedAsset.Properties)
		if err != nil {
			return []models.Asset{}, err
		}

		err = json.NewDecoder(strings.NewReader(urls)).Decode(&fetchedAsset.URLs)
		if err != nil {
			return []models.Asset{}, err
		}
		err = json.NewDecoder(strings.NewReader(legalDocs)).Decode(&fetchedAsset.LegalDocumentURLs)
		if err != nil {
			return []models.Asset{}, err
		}
		err = json.NewDecoder(strings.NewReader(issuanceDocs)).Decode(&fetchedAsset.IssuanceDocumentURLs)
		if err != nil {
			return []models.Asset{}, err
		}
		err = json.NewDecoder(strings.NewReader(signatories)).Decode(&fetchedAsset.Signatories)
		if err != nil {
			return []models.Asset{}, err
		}

		assets = append(assets, fetchedAsset)
	}

	if err = rows.Err(); err != nil {
		return []models.Asset{}, err
	}

	return assets, nil
}
func (a assetImpl) UpdateMintStatus(assetId string) error {
	stmt := `
	UPDATE assets
	SET  is_minted=true
	WHERE id=$1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := a.Db.ExecContext(
		ctx,
		stmt,
		assetId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a assetImpl) GetOrders(assetId string, status *models.OrderStatus, orderType *models.OrderType, filter models.Filter) ([]models.Order, int, error) {

	var totalCount int

	stmt := `
SELECT 
id, 
user_id, 
asset_id, 
"type", 
kind, 
status, 
price, 
quantity,
inital_quantity,
created_at, 
updated_at
FROM orders
WHERE asset_id=$1 and
($2::text IS NULL OR status = $2::text)
    AND ($3::text IS NULL OR "type" = $3::text)
 LIMIT $4 OFFSET $5;
		`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	orders := []models.Order{}

	rows, err := a.Db.QueryContext(
		ctx,
		stmt,
		assetId,
		status,
		orderType,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return []models.Order{}, totalCount, err
	}
	defer rows.Close()

	totalCountQuery := `
	   SELECT COUNT(*)
		FROM orders
	WHERE asset_id=$1 and
	($2::text IS NULL OR status = $2::text)
		AND ($3::text IS NULL OR "type" = $3::text);
	`

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = a.Db.QueryRowContext(ctx, totalCountQuery,
		assetId,
		status,
		orderType,
	).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		order := models.Order{}

		rows.Scan(
			&order.Id,
			&order.UserId,
			&order.AssetId,
			&order.Type,
			&order.Kind,
			&order.Status,
			&order.Price,
			&order.Quantity,
			&order.InitialQty,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return []models.Order{}, totalCount, err
	}

	return orders, totalCount, nil
}

// func (a assetImpl) Approve(assetId string) error {
// 	stmt := `
// 	UPDATE assets
// 	SET status = true
// 	WHERE id = $1;
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	_, err := a.Db.ExecContext(
// 		ctx,
// 		stmt,
// 		assetId,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
// func (a assetImpl) Reject(assetId string) error {
// 	stmt := `
// 	UPDATE assets
// 	SET is_rejected = true
// 	WHERE id = $1;
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	_, err := a.Db.ExecContext(
// 		ctx,
// 		stmt,
// 		assetId,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (a assetImpl) UpdateStatus(assetId string, status models.AssetStatus) error {
	stmt := `
	UPDATE assets
	SET status = $2
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := a.Db.ExecContext(
		ctx,
		stmt,
		assetId,
		status,
	)
	if err != nil {
		return err
	}

	return nil
}
