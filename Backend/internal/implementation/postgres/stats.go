package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
)

type statImplementation struct {
	Db *sql.DB
}

func NewStatImplementation(db *sql.DB) repository.StatsRepository {
	return statImplementation{
		Db: db,
	}
}

func (s statImplementation) calculateFloorPrice(assetId string, period int) (float64, error) {
	// this query is supposed to go through the list of prices items in activities
	var floorPrice float64
	query := `
	SELECT COALESCE(MIN(price), 0)
	FROM listings
	WHERE asset_id = $1
	AND is_active = true
	AND is_cancelled = false
	AND created_at >= NOW() - INTERVAL '1 hour' * $2
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

	err := s.Db.QueryRowContext(ctx, query, assetId, period).Scan(&floorPrice)
	if err != nil {
		return 0, err
	}

	return floorPrice, nil
}

func (s statImplementation) TopAssets(filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error) {
	topassets := make([]models.AssetStat, 0)

	query := `
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
        act.created_at >= NOW() - INTERVAL '1 hour' * $3
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
	a.is_minted
FROM
    assets a
JOIN users AS u ON a.creator_id = u.id
LEFT JOIN
    asset_sales asl ON a.id = asl.asset_id
LEFT JOIN
    asset_holders ah ON a.id = ah.asset_id
Left JOIN categories AS c ON a.category_id = c.id
WHERE    
($4::int IS NULL OR a.category_id = $4::int)
    AND ($5::text IS NULL OR a.blockchain = $5::text)
ORDER BY
    a.id,asl.sale_activity_count DESC, asl.total_sales_volume DESC, ah.holder_count DESC
LIMIT $1 OFFSET $2;

    `

	var totalCount int
	// // Return the total count along with the trending assets
	// return topassets, totalCount, nil
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
		filter.Limit,
		filter.Offset(),
		filter.Range,
		categoryId,
		blockchain,
	)

	if err != nil {
		return nil, totalCount, err
	}

	defer rows.Close()

	totalCountQuery := `
	   SELECT COUNT(*)
    FROM assets a
    JOIN users u ON a.creator_id = u.id
    LEFT JOIN (
        SELECT
            act.asset_id
        FROM
            activities act
        WHERE
            act.action = 'sale' AND
            act.created_at >= NOW() - INTERVAL '1 hour' * $1
        GROUP BY
            act.asset_id
    ) act ON a.id = act.asset_id
    LEFT JOIN (
        SELECT
            uh.asset_id
        FROM
            asset_owners uh
        GROUP BY
            uh.asset_id
    ) uh ON a.id = uh.asset_id
    LEFT JOIN categories c ON a.category_id = c.id
    WHERE
        ($2::int IS NULL OR a.category_id = $2::int)
        AND ($3::text IS NULL OR a.blockchain = $3::text);
	`

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = s.Db.QueryRowContext(ctx, totalCountQuery, filter.Range, categoryId, blockchain).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		topasset := models.AssetStat{}
		urls := ""
		err = rows.Scan(
			&topasset.AssetId,
			&topasset.Name,
			&topasset.Symbol,
			&topasset.CategoryId,
			&topasset.CategoryName,
			&topasset.Blockchain,
			&urls,
			&topasset.Status,
			&topasset.CreatorId,
			&topasset.CreatorUsername,
			&topasset.TotalVolume,
			&topasset.ActivityCount,
			&topasset.Owners,
			&topasset.IsMinted,
		)
		if err != nil {
			return nil, totalCount, err
		}

		err := json.NewDecoder(strings.NewReader(urls)).Decode(&topasset.URLs)
		if err != nil {
			return nil, totalCount, err
		}

		topasset.FloorPrice, err = s.calculateFloorPrice(topasset.AssetId, int(filter.Range))
		if err != nil {
			return nil, totalCount, err
		}
		topassets = append(topassets, topasset)

	}
	return topassets, totalCount, nil
}

func (s statImplementation) TrendingAssets(filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error) {
	topassets := make([]models.AssetStat, 0)

	query := `
	WITH asset_sales AS (
    SELECT
        a.id AS asset_id,
        COUNT(act.id) AS sale_activity_count
    FROM
        assets a
    LEFT JOIN
        activities act ON a.id = act.asset_id AND act.action = 'sale' AND act.created_at >= NOW() - INTERVAL '1 hour' * $3
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
    c."name" AS category_name,
    a.blockchain,
    a.urls,
    a.status,
    a.creator_id,
    u.username,
    COALESCE(asl.sale_activity_count, 0) AS sale_activity_count,
    COALESCE(ah.holder_count, 0) AS holder_count,
    a.is_minted
FROM
    assets a
JOIN
    users u ON a.creator_id = u.id
LEFT JOIN
    asset_sales asl ON a.id = asl.asset_id
LEFT JOIN
    asset_holders ah ON a.id = ah.asset_id
LEFT JOIN
    categories c ON a.category_id = c.id
WHERE
    ($4::int IS NULL OR a.category_id = $4::int)
    AND ($5::text IS NULL OR a.blockchain = $5::text)
ORDER BY
    a.id,asl.sale_activity_count DESC, ah.holder_count DESC
LIMIT $1 OFFSET $2;


    `
	var totalCount int

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fmt.Println("OFFSET", filter.Offset())
	rows, err := s.Db.QueryContext(
		ctx,
		query,
		filter.Limit,
		filter.Offset(),
		filter.Range,
		categoryId,
		blockchain,
	)
	if err != nil {
		return nil, totalCount, err
	}

	defer rows.Close()

	totalCountQuery := `
	SELECT COUNT(*)
 FROM assets a
 JOIN users u ON a.creator_id = u.id
 LEFT JOIN (
	 SELECT
		 act.asset_id
	 FROM
		 activities act
	 WHERE
		 act.action = 'sale' AND
		 act.created_at >= NOW() - INTERVAL '1 hour' * $1
	 GROUP BY
		 act.asset_id
 ) act ON a.id = act.asset_id
 LEFT JOIN (
	 SELECT
		 uh.asset_id
	 FROM
		 asset_owners uh
	 GROUP BY
		 uh.asset_id
 ) uh ON a.id = uh.asset_id
 LEFT JOIN categories c ON a.category_id = c.id
 WHERE
	 ($2::int IS NULL OR a.category_id = $2::int)
	 AND ($3::text IS NULL OR a.blockchain = $3::text);
 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = s.Db.QueryRowContext(ctx, totalCountQuery, filter.Range, categoryId, blockchain).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		topasset := models.AssetStat{}
		urls := ""
		err = rows.Scan(
			&topasset.AssetId,
			&topasset.Name,
			&topasset.Symbol,
			&topasset.CategoryId,
			&topasset.CategoryName,
			&topasset.Blockchain,
			&urls,
			&topasset.Status,
			&topasset.CreatorId,
			&topasset.CreatorUsername,
			&topasset.ActivityCount,
			&topasset.Owners,
			&topasset.IsMinted,
		)
		if err != nil {
			return nil, totalCount, err
		}

		err := json.NewDecoder(strings.NewReader(urls)).Decode(&topasset.URLs)
		if err != nil {
			return nil, totalCount, err
		}

		topasset.FloorPrice, err = s.calculateFloorPrice(topasset.AssetId, int(filter.Range))
		if err != nil {
			return nil, totalCount, err
		}
		topassets = append(topassets, topasset)

	}
	return topassets, totalCount, nil
}

func (s statImplementation) TrendingAssetPerfs(filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error) {
	topassets := make([]models.AssetStat, 0)
	var totalCount int
	query := `
	WITH asset_sales AS (
    SELECT
        a.id AS asset_id,
        a.name,
        a.symbol,
        COUNT(act.id) AS sale_activity_count
    FROM
        assets a
    LEFT JOIN
        activities act ON a.id = act.asset_id
    WHERE
        act.action = 'sale' AND
        act.created_at >= NOW() - INTERVAL '1 hour' * $3
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
    COALESCE(asl.sale_activity_count, 0) AS sale_activity_count,
    COALESCE(ah.holder_count, 0) AS holder_count,
	a.is_minted
FROM
    assets a
JOIN users AS u ON a.creator_id = u.id
LEFT JOIN
    asset_sales asl ON a.id = asl.asset_id
LEFT JOIN
    asset_holders ah ON a.id = ah.asset_id
Left JOIN categories AS c ON a.category_id = c.id
WHERE    
 ($4::int IS NULL OR a.category_id = $4::int)
    AND ($5::text IS NULL OR a.blockchain = $5::text)
ORDER BY
     a.id,asl.sale_activity_count DESC, ah.holder_count DESC
LIMIT $1 OFFSET $2;

    `

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
		filter.Limit,
		filter.Offset(),
		filter.Range,
		categoryId,
		blockchain,
	)
	if err != nil {
		return nil, totalCount, err
	}

	defer rows.Close()

	totalCountQuery := `
	SELECT COUNT(*)
 FROM assets a
 JOIN users u ON a.creator_id = u.id
 LEFT JOIN (
	 SELECT
		 act.asset_id
	 FROM
		 activities act
	 WHERE
		 act.action = 'sale' AND
		 act.created_at >= NOW() - INTERVAL '1 hour' * $1
	 GROUP BY
		 act.asset_id
 ) act ON a.id = act.asset_id
 LEFT JOIN (
	 SELECT
		 uh.asset_id
	 FROM
		 asset_owners uh
	 GROUP BY
		 uh.asset_id
 ) uh ON a.id = uh.asset_id
 LEFT JOIN categories c ON a.category_id = c.id
 WHERE
	 ($2::int IS NULL OR a.category_id = $2::int)
	 AND ($3::text IS NULL OR a.blockchain = $3::text);
 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = s.Db.QueryRowContext(ctx, totalCountQuery, filter.Range, categoryId, blockchain).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}

	for rows.Next() {
		topasset := models.AssetStat{}
		urls := ""
		err = rows.Scan(
			&topasset.AssetId,
			&topasset.Name,
			&topasset.Symbol,
			&topasset.CategoryId,
			&topasset.CategoryName,
			&topasset.Blockchain,
			&urls,
			&topasset.Status,
			&topasset.CreatorId,
			&topasset.CreatorUsername,
			&topasset.ActivityCount,
			&topasset.Owners,
			&topasset.IsMinted,
		)
		if err != nil {
			return nil, totalCount, err
		}

		err := json.NewDecoder(strings.NewReader(urls)).Decode(&topasset.URLs)
		if err != nil {
			return nil, totalCount, err
		}

		topasset.FloorPrice, err = s.calculateFloorPrice(topasset.AssetId, int(filter.Range))
		if err != nil {
			return nil, totalCount, err
		}
		topassets = append(topassets, topasset)

	}
	return topassets, totalCount, nil
}
func (s statImplementation) TrendingUserAssetPerfs(userId string, filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error) {
	topassets := make([]models.AssetStat, 0)
	var totalCount int
	query := `
	WITH asset_sales AS (
    SELECT
        a.id AS asset_id,
        a.name,
        a.symbol,
        COUNT(act.id) AS sale_activity_count
    FROM
        assets a
    LEFT JOIN
        activities act ON a.id = act.asset_id
    WHERE
        act.action = 'sale' AND
        act.created_at >= NOW() - INTERVAL '1 hour' * $3
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
    COALESCE(asl.sale_activity_count, 0) AS sale_activity_count,
    COALESCE(ah.holder_count, 0) AS holder_count,
	a.is_minted
FROM
    assets a
JOIN users AS u ON a.creator_id = u.id
LEFT JOIN
    asset_sales asl ON a.id = asl.asset_id
LEFT JOIN
    asset_holders ah ON a.id = ah.asset_id
Left JOIN categories AS c ON a.category_id = c.id
WHERE    
a.creator_id=$6 AND
 ($4::int IS NULL OR a.category_id = $4::int)
    AND ($5::text IS NULL OR a.blockchain = $5::text)
ORDER BY
     a.id,asl.sale_activity_count DESC, ah.holder_count DESC
LIMIT $1 OFFSET $2;

    `

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
		filter.Limit,
		filter.Offset(),
		filter.Range,
		categoryId,
		blockchain,
		userId,
	)
	if err != nil {
		return nil, totalCount, err
	}

	defer rows.Close()

	totalCountQuery := `
	SELECT COUNT(*)
 FROM assets a
 JOIN users u ON a.creator_id = u.id
 LEFT JOIN (
	 SELECT
		 act.asset_id
	 FROM
		 activities act
	 WHERE
		 act.action = 'sale' AND
		 act.created_at >= NOW() - INTERVAL '1 hour' * $1
	 GROUP BY
		 act.asset_id
 ) act ON a.id = act.asset_id
 LEFT JOIN (
	 SELECT
		 uh.asset_id
	 FROM
		 asset_owners uh
	 GROUP BY
		 uh.asset_id
 ) uh ON a.id = uh.asset_id
 LEFT JOIN categories c ON a.category_id = c.id
 WHERE
    a.creator_id=$4 AND
	 ($2::int IS NULL OR a.category_id = $2::int)
	 AND ($3::text IS NULL OR a.blockchain = $3::text);
 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = s.Db.QueryRowContext(ctx, totalCountQuery, filter.Range, categoryId, blockchain, userId).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}

	for rows.Next() {
		topasset := models.AssetStat{}
		urls := ""
		err = rows.Scan(
			&topasset.AssetId,
			&topasset.Name,
			&topasset.Symbol,
			&topasset.CategoryId,
			&topasset.CategoryName,
			&topasset.Blockchain,
			&urls,
			&topasset.Status,
			&topasset.CreatorId,
			&topasset.CreatorUsername,
			&topasset.ActivityCount,
			&topasset.Owners,
			&topasset.IsMinted,
		)
		if err != nil {
			return nil, totalCount, err
		}

		err := json.NewDecoder(strings.NewReader(urls)).Decode(&topasset.URLs)
		if err != nil {
			return nil, totalCount, err
		}

		topasset.FloorPrice, err = s.calculateFloorPrice(topasset.AssetId, int(filter.Range))
		if err != nil {
			return nil, totalCount, err
		}
		topassets = append(topassets, topasset)

	}
	return topassets, totalCount, nil
}

func (s statImplementation) TopAssetPerfs(filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error) {
	topassets := make([]models.AssetStat, 0)
	var totalCount int
	query := `
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
		act.created_at >= NOW() - INTERVAL '1 hour' * $3
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
		a.is_minted
	FROM
		assets a
	JOIN users AS u ON a.creator_id = u.id
	LEFT JOIN
		asset_sales asl ON a.id = asl.asset_id
	LEFT JOIN
		asset_holders ah ON a.id = ah.asset_id
	Left JOIN categories AS c ON a.category_id = c.id
	WHERE    
	($4::int IS NULL OR a.category_id = $4::int)
	AND ($5::text IS NULL OR a.blockchain = $5::text)
ORDER BY
	 a.id,asl.sale_activity_count DESC, asl.total_sales_volume DESC, ah.holder_count DESC
LIMIT $1 OFFSET $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
		filter.Limit,
		filter.Offset(),
		filter.Range,
		categoryId,
		blockchain,
	)
	if err != nil {
		return nil, totalCount, err
	}

	defer rows.Close()

	totalCountQuery := `
	SELECT COUNT(*)
 FROM assets a
 JOIN users u ON a.creator_id = u.id
 LEFT JOIN (
	 SELECT
		 act.asset_id
	 FROM
		 activities act
	 WHERE
		 act.action = 'sale' AND
		 act.created_at >= NOW() - INTERVAL '1 hour' * $1
	 GROUP BY
		 act.asset_id
 ) act ON a.id = act.asset_id
 LEFT JOIN (
	 SELECT
		 uh.asset_id
	 FROM
		 asset_owners uh
	 GROUP BY
		 uh.asset_id
 ) uh ON a.id = uh.asset_id
 LEFT JOIN categories c ON a.category_id = c.id
 WHERE
	 ($2::int IS NULL OR a.category_id = $2::int)
	 AND ($3::text IS NULL OR a.blockchain = $3::text);
 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = s.Db.QueryRowContext(ctx, totalCountQuery, filter.Range, categoryId, blockchain).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		topasset := models.AssetStat{}
		urls := ""
		err = rows.Scan(
			&topasset.AssetId,
			&topasset.Name,
			&topasset.Symbol,
			&topasset.CategoryId,
			&topasset.CategoryName,
			&topasset.Blockchain,
			&urls,
			&topasset.Status,
			&topasset.CreatorId,
			&topasset.CreatorUsername,
			&topasset.TotalVolume,
			&topasset.ActivityCount,
			&topasset.Owners,
			&topasset.IsMinted,
		)
		if err != nil {
			return nil, totalCount, err
		}

		err := json.NewDecoder(strings.NewReader(urls)).Decode(&topasset.URLs)
		if err != nil {
			return nil, totalCount, err
		}

		topasset.FloorPrice, err = s.calculateFloorPrice(topasset.AssetId, int(filter.Range))
		if err != nil {
			return nil, totalCount, err
		}
		topassets = append(topassets, topasset)

	}
	return topassets, totalCount, nil
}
func (s statImplementation) TopUserAssetPerfs(userId string, filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error) {
	topassets := make([]models.AssetStat, 0)
	var totalCount int
	query := `
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
		act.created_at >= NOW() - INTERVAL '1 hour' * $3
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
		a.is_minted
	FROM
		assets a
	JOIN users AS u ON a.creator_id = u.id
	LEFT JOIN
		asset_sales asl ON a.id = asl.asset_id
	LEFT JOIN
		asset_holders ah ON a.id = ah.asset_id
	Left JOIN categories AS c ON a.category_id = c.id
	WHERE    
	a.creator_id=$6 AND
	($4::int IS NULL OR a.category_id = $4::int)
	AND ($5::text IS NULL OR a.blockchain = $5::text)
ORDER BY
	 a.id,asl.sale_activity_count DESC, asl.total_sales_volume DESC, ah.holder_count DESC
LIMIT $1 OFFSET $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
		filter.Limit,
		filter.Offset(),
		filter.Range,
		categoryId,
		blockchain,
		userId,
	)
	if err != nil {
		return nil, totalCount, err
	}

	defer rows.Close()

	totalCountQuery := `
	SELECT COUNT(*)
 FROM assets a
 JOIN users u ON a.creator_id = u.id
 LEFT JOIN (
	 SELECT
		 act.asset_id
	 FROM
		 activities act
	 WHERE
		 act.action = 'sale' AND
		 act.created_at >= NOW() - INTERVAL '1 hour' * $1
	 GROUP BY
		 act.asset_id
 ) act ON a.id = act.asset_id
 LEFT JOIN (
	 SELECT
		 uh.asset_id
	 FROM
		 asset_owners uh
	 GROUP BY
		 uh.asset_id
 ) uh ON a.id = uh.asset_id
 LEFT JOIN categories c ON a.category_id = c.id
 WHERE
 a.creator_id=$4 AND
	 ($2::int IS NULL OR a.category_id = $2::int)
	 AND ($3::text IS NULL OR a.blockchain = $3::text);
 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = s.Db.QueryRowContext(ctx, totalCountQuery, filter.Range, categoryId, blockchain, userId).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		topasset := models.AssetStat{}
		urls := ""
		err = rows.Scan(
			&topasset.AssetId,
			&topasset.Name,
			&topasset.Symbol,
			&topasset.CategoryId,
			&topasset.CategoryName,
			&topasset.Blockchain,
			&urls,
			&topasset.Status,
			&topasset.CreatorId,
			&topasset.CreatorUsername,
			&topasset.TotalVolume,
			&topasset.ActivityCount,
			&topasset.Owners,
			&topasset.IsMinted,
		)
		if err != nil {
			return nil, totalCount, err
		}

		err := json.NewDecoder(strings.NewReader(urls)).Decode(&topasset.URLs)
		if err != nil {
			return nil, totalCount, err
		}

		topasset.FloorPrice, err = s.calculateFloorPrice(topasset.AssetId, int(filter.Range))
		if err != nil {
			return nil, totalCount, err
		}
		topassets = append(topassets, topasset)

	}
	return topassets, totalCount, nil
}

func (s statImplementation) GetAssetPriceData(assetId string, filter models.Filter) ([]models.AssetPriceData, error) {
	query := `
		SELECT
		act.price,
		act.created_at AS timestamp
	FROM
		activities act
	WHERE
		act.action = 'sale' AND
		act.asset_id = $1
		AND act.created_at >= NOW() - INTERVAL '1 hour' * $2
	ORDER BY
		act.created_at DESC
	LIMIT 5;

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(ctx, query, assetId, filter.Range)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	priceData := make([]models.AssetPriceData, 0)
	for rows.Next() {
		price := models.AssetPriceData{}
		err = rows.Scan(
			&price.Price,
			&price.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		priceData = append(priceData, price)
	}

	return priceData, nil
}

func (s statImplementation) GetAssetCreationSurvey() ([]models.CreatorSurveyMonthlyStat, error) {
	monthlyStats := make([]models.CreatorSurveyMonthlyStat, 0)

	query := `
WITH months AS (
       SELECT
        generate_series(
            DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '11 months',
            DATE_TRUNC('month', CURRENT_DATE),
            INTERVAL '1 month'
        )::date AS month_start
),
user_creation_month AS (
    SELECT
        id,
        DATE_TRUNC('month', created_at) AS user_created_month
    FROM
        users
),
assets_created_month AS (
    SELECT
        creator_id,
        DATE_TRUNC('month', created_at) AS asset_created_month
    FROM
        assets
)
SELECT
    TO_CHAR(m.month_start, 'MM') AS month,
    COALESCE(COUNT(DISTINCT CASE WHEN u.user_created_month = a.asset_created_month THEN u.id END), 0) AS new_users,
    COALESCE(COUNT(DISTINCT CASE WHEN u.user_created_month < a.asset_created_month THEN u.id END), 0) AS old_users
FROM
    months m
LEFT JOIN
    assets_created_month a ON DATE_TRUNC('month', a.asset_created_month) = m.month_start
LEFT JOIN
    user_creation_month u ON a.creator_id = u.id
GROUP BY
    m.month_start
ORDER BY
    m.month_start;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		monthlyStat := models.CreatorSurveyMonthlyStat{}

		err = rows.Scan(
			&monthlyStat.Month,
			&monthlyStat.NewUsersCount,
			&monthlyStat.HeritageUsersCount,
		)
		if err != nil {
			return nil, err
		}

		monthlyStats = append(monthlyStats, monthlyStat)

	}
	return monthlyStats, nil
}
func (s statImplementation) GetAssetCategorySurvey() ([]models.AssetSurveyMonthlyStat, error) {
	monthlyStats := make([]models.AssetSurveyMonthlyStat, 0)

	query := `
WITH months AS (
    SELECT
           generate_series(
            DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '11 months',
            DATE_TRUNC('month', CURRENT_DATE),
            INTERVAL '1 month'
        )::date AS month_start
)
SELECT
    TO_CHAR(m.month_start, 'MM') AS month,
    COALESCE(SUM(CASE WHEN a.category_id = 1 THEN 1 ELSE 0 END), 0) AS category_1,
    COALESCE(SUM(CASE WHEN a.category_id = 2 THEN 1 ELSE 0 END), 0) AS category_2,
    COALESCE(SUM(CASE WHEN a.category_id = 3 THEN 1 ELSE 0 END), 0) AS category_3,
    -- Add more categories as needed
    COALESCE(SUM(CASE WHEN a.category_id = 4 THEN 1 ELSE 0 END), 0) AS category_4
    -- Add more categories as needed
FROM
    months m
LEFT JOIN
    assets a ON DATE_TRUNC('month', a.created_at) = m.month_start
GROUP BY
    m.month_start
ORDER BY
    m.month_start;

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		monthlyStat := models.AssetSurveyMonthlyStat{}

		err = rows.Scan(
			&monthlyStat.Month,
			&monthlyStat.Category1,
			&monthlyStat.Category2,
			&monthlyStat.Category3,
			&monthlyStat.Category4,
		)
		if err != nil {
			return nil, err
		}

		monthlyStats = append(monthlyStats, monthlyStat)

	}
	return monthlyStats, nil
}
func (s statImplementation) SalesTrend(userId string, filter models.Filter) ([]models.SalesTrend, error) {
	saleTrends := make([]models.SalesTrend, 0)

	query := `

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		saleTrend := models.SalesTrend{}

		err = rows.Scan(
			&saleTrend.Period,
			&saleTrend.Sales,
			&saleTrend.SalesValue,
		)
		if err != nil {
			return nil, err
		}

		saleTrends = append(saleTrends, saleTrend)

	}
	return saleTrends, nil
}

func (s statImplementation) UsersCount() (int, error) {

	var totalCount int

	totalCountQuery := `
SELECT COUNT(*) AS total_users
FROM users;
 `

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := s.Db.QueryRowContext(ctx, totalCountQuery).Scan(&totalCount)
	if err != nil {
		return totalCount, err
	}

	return totalCount, nil
}
func (s statImplementation) CreatorsCount() (int, error) {

	var totalCount int

	totalCountQuery := `
SELECT COUNT(DISTINCT creator_id) AS total_creators
FROM assets;
 `

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := s.Db.QueryRowContext(ctx, totalCountQuery).Scan(&totalCount)
	if err != nil {
		return totalCount, err
	}

	return totalCount, nil
}

func (s statImplementation) AssetsCount() (int, error) {

	var totalCount int

	totalCountQuery := `
SELECT COUNT(*) AS total_assets
FROM assets;
 `

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := s.Db.QueryRowContext(ctx, totalCountQuery).Scan(&totalCount)
	if err != nil {
		return totalCount, err
	}

	return totalCount, nil
}

func (s statImplementation) CreatorsWeekIncrement() (float64, error) {

	var totalCount float64

	totalCountQuery := `
WITH creators_this_week AS (
    SELECT COUNT(DISTINCT creator_id) AS num_creators
    FROM assets
    WHERE created_at >= NOW() - INTERVAL '1 week'
),
creators_last_week AS (
    SELECT COUNT(DISTINCT creator_id) AS num_creators
    FROM assets
    WHERE created_at >= NOW() - INTERVAL '2 weeks' AND created_at < NOW() - INTERVAL '1 week'
)
SELECT 
    CASE 
        WHEN clw.num_creators = 0 THEN 0  -- Avoid division by zero
        ELSE (ctw.num_creators - clw.num_creators) * 100.0 / clw.num_creators
    END AS percent_increase_creators
FROM 
    creators_this_week ctw,
    creators_last_week clw;

 `

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := s.Db.QueryRowContext(ctx, totalCountQuery).Scan(&totalCount)
	if err != nil {
		return totalCount, err
	}

	return totalCount, nil
}
func (s statImplementation) UsersWeekIncrement() (float64, error) {

	var totalCount float64

	totalCountQuery := `
WITH users_this_week AS (
    SELECT COUNT(*) AS num_users
    FROM users
    WHERE created_at >= NOW() - INTERVAL '1 week'
),
users_last_week AS (
    SELECT COUNT(*) AS num_users
    FROM users
    WHERE created_at >= NOW() - INTERVAL '2 weeks' AND created_at < NOW() - INTERVAL '1 week'
)
SELECT 
    CASE 
        WHEN ulw.num_users = 0 THEN 0  -- Avoid division by zero
        ELSE (utw.num_users - ulw.num_users) * 100.0 / ulw.num_users
    END AS percent_increase_users
FROM 
    users_this_week utw,
    users_last_week ulw;


 `

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := s.Db.QueryRowContext(ctx, totalCountQuery).Scan(&totalCount)
	if err != nil {
		return totalCount, err
	}

	return totalCount, nil
}

func (s statImplementation) TopCreators(filter models.Filter, categoryId *int, blockchain *string) ([]models.CreatorStats, int, error) {
	topCreators := make([]models.CreatorStats, 0)

	query := `
WITH creator_sales AS (
    SELECT
        a.creator_id,
        COUNT(act.id) AS sale_activity_count,
        COALESCE(SUM(act.price), 0) AS total_sales_volume
    FROM
        assets a
    LEFT JOIN
        activities act ON a.id = act.asset_id
    WHERE
        act.action = 'sale' AND
        act.created_at >= NOW() - INTERVAL '1 hour' * $3
    GROUP BY
        a.creator_id
),
creator_assets AS (
    SELECT
        a.creator_id,
        COUNT(a.id) AS asset_count
    FROM
        assets a
    GROUP BY
        a.creator_id
)
SELECT
    u.id AS creator_id,
    u.username,
    -- COALESCE(cs.total_sales_volume, 0) AS total_sales_volume,
    -- COALESCE(cs.sale_activity_count, 0) AS sale_activity_count,
    COALESCE(ca.asset_count, 0) AS asset_count
FROM
    users u
LEFT JOIN
    creator_sales cs ON u.id = cs.creator_id
LEFT JOIN
    creator_assets ca ON u.id = ca.creator_id
WHERE
    ($4::int IS NULL OR EXISTS (
        SELECT 1 
        FROM assets a 
        WHERE a.creator_id = u.id AND a.category_id = $4::int
    ))
    AND ($5::text IS NULL OR EXISTS (
        SELECT 1 
        FROM assets a 
        WHERE a.creator_id = u.id AND a.blockchain = $5::text
    ))
ORDER BY
    cs.total_sales_volume DESC, cs.sale_activity_count DESC, ca.asset_count DESC
LIMIT $1 OFFSET $2;


    `

	var totalCount int
	// // Return the total count along with the trending assets
	// return topassets, totalCount, nil
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
		filter.Limit,
		filter.Offset(),
		filter.Range,
		categoryId,
		blockchain,
	)

	if err != nil {
		return nil, totalCount, err
	}

	defer rows.Close()

	totalCountQuery := `
	WITH creator_sales AS (
    SELECT
        a.creator_id,
        COUNT(act.id) AS sale_activity_count,
        COALESCE(SUM(act.price), 0) AS total_sales_volume
    FROM
        assets a
    LEFT JOIN
        activities act ON a.id = act.asset_id
    WHERE
        act.action = 'sale' AND
        act.created_at >= NOW() - INTERVAL '1 hour' * $1
    GROUP BY
        a.creator_id
),
creator_assets AS (
    SELECT
        a.creator_id,
        COUNT(a.id) AS asset_count
    FROM
        assets a
    GROUP BY
        a.creator_id
)
	   SELECT COUNT(*)
	FROM
		users u
	LEFT JOIN
		creator_sales cs ON u.id = cs.creator_id
	LEFT JOIN
		creator_assets ca ON u.id = ca.creator_id
	WHERE
		($2::int IS NULL OR EXISTS (
			SELECT 1 
			FROM assets a 
			WHERE a.creator_id = u.id AND a.category_id = $2::int
		))
		AND ($3::text IS NULL OR EXISTS (
			SELECT 1 
			FROM assets a 
			WHERE a.creator_id = u.id AND a.blockchain = $3::text
		))
	`

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = s.Db.QueryRowContext(ctx, totalCountQuery,
		filter.Range,
		categoryId,
		blockchain,
	).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		topCreator := models.CreatorStats{}

		err = rows.Scan(
			&topCreator.CreatorId,
			&topCreator.CreatorUsername,
			&topCreator.AssetsCount,
		)
		if err != nil {
			return nil, totalCount, err
		}

		topCreators = append(topCreators, topCreator)

	}
	return topCreators, totalCount, nil
}

// 	// this query is to get the current price for which an item is sold in a asset.
// 	var currentPrice float64
// 	query := `
// 	SELECT price
// 	FROM activities a
//       JOIN items i ON i.token_id::text = a.token_id
// 	  JOIN assets c ON c.id = i.asset_id AND c.contract_address = a.contract_address
// 	WHERE a.contract_address = $1
// 	  AND c.id = $2
// 	ORDER BY a.occurred_at DESC
// 	LIMIT 1
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	err := s.Db.QueryRowContext(ctx, query, asset.AssetId).Scan(&currentPrice)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return currentPrice, nil
// }

// func (s statImplementation) getPreviousPrice(asset models.TopAsset, period uint) (float64, error) {
// 	// this query is to get the price for which an item is sold or bought at a particular time an
// 	var previoustPrice float64
// 	query := `
// 	SELECT COALESCE(
// 		(SELECT price
// 		 FROM activities a
// 		 JOIN items i ON i.token_id::text = a.token_id
// 		 JOIN assets c ON c.id = i.asset_id AND c.contract_address = a.contract_address
// 		 WHERE a.contract_address = $1
// 		   AND c.id = $2
// 		   AND AGE(a.occurred_at) >= INTERVAL '1 day' * $3
// 		 ORDER BY a.occurred_at DESC
// 		 LIMIT 1),
// 		(SELECT price
// 		 FROM activities a
// 		 JOIN items i ON i.token_id::text = a.token_id
// 		 JOIN assets c ON c.id = i.asset_id AND c.contract_address = a.contract_address
// 		 WHERE a.contract_address = $1
// 		   AND c.id = $2
// 		 ORDER BY a.occurred_at ASC
// 		 LIMIT 1)
// 	  ) AS current_price
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	err := s.Db.QueryRowContext(ctx, query, asset.AssetId, period).Scan(&previoustPrice)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return previoustPrice, nil
// }
