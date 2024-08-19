package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
)

type userImplementation struct {
	Db *sql.DB
}

func NewUserImplementation(db *sql.DB) repository.UserRepository {
	return userImplementation{Db: db}
}

// users table constraints
const (
	duplicateAddress   = "users_address_key"
	duplicateAccountId = "users_account_id_key"
	duplicateUsername  = "users_username_key"
	duplicateEmail     = "users_email_key"
)

func (u userImplementation) Create(user models.User) (models.User, error) {
	// generate an ID for the user if it doesn't have one
	if user.Id == "" {
		user.Id = uuid.NewString()
	}

	stmt := `
	INSERT INTO users(id, address, account_id,user_type)
	VALUES ($1, $2, $3,$4)
	RETURNING 
	    id,
	    address,
	    account_id, 
	    created_at, 
	    updated_at;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	newUser := models.User{}
	err := u.Db.QueryRowContext(ctx, stmt, user.Id, user.Address, user.AccountID, user.UserType).Scan(
		&newUser.Id,
		&newUser.Address,
		&newUser.AccountID,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), duplicateAddress),
			strings.Contains(err.Error(), duplicateAccountId),
			strings.Contains(err.Error(), duplicateUsername),
			strings.Contains(err.Error(), duplicateEmail):
			return models.User{}, repository.ErrDuplicateDetails

		default:
			return models.User{}, err
		}
	}

	return newUser, err
}

// GetByAddress returns the user with the given address.
// repository.ErrRecordNotFound is returned if no row with the matching address is found.
func (u userImplementation) GetByAddress(address string) (models.User, error) {
	stmt := `
SELECT  
   id, 
   address,
   account_id,
   email, 
   username,
   first_name, 
   last_name, 
   bio, 
   website, 
   twitter, 
   discord, 
   avatar,
   public_key,
   encrypted_private_key,
   is_verified, 
   is_active, 
   kyc_registered_date,
   created_at,
   updated_at,
   user_type
FROM users
WHERE address = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := models.User{}
	err := u.Db.QueryRowContext(ctx, stmt, address).Scan(
		&user.Id,
		&user.Address,
		&user.AccountID,
		&user.Email,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Bio,
		&user.Website,
		&user.Twitter,
		&user.Discord,
		&user.Avatar,
		&user.PublicKey,
		&user.EncryptedPrivateKey,
		&user.IsVerified,
		&user.IsActive,
		&user.KYCRegisteredDate,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.UserType,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repository.ErrRecordNotFound

		default:
			return models.User{}, err
		}
	}

	return user, nil
}

func (u userImplementation) GetById(id string) (models.User, error) {
	stmt := `
	SELECT  
   id, 
   address,
   account_id,
   email, 
   username,
   first_name, 
   last_name, 
   bio, 
   website, 
   twitter, 
   discord, 
   avatar,
   public_key,
   encrypted_private_key,
   is_verified, 
   is_active, 
   kyc_registered_date,
   created_at,
   updated_at,
   user_type
FROM users
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := models.User{}
	err := u.Db.QueryRowContext(ctx, stmt, id).Scan(
		&user.Id,
		&user.Address,
		&user.AccountID,
		&user.Email,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Bio,
		&user.Website,
		&user.Twitter,
		&user.Discord,
		&user.Avatar,
		&user.PublicKey,
		&user.EncryptedPrivateKey,
		&user.IsVerified,
		&user.IsActive,
		&user.KYCRegisteredDate,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.UserType,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repository.ErrRecordNotFound

		default:
			return models.User{}, err
		}
	}

	return user, nil
}

func (u userImplementation) SetPrivateKey(userId string, encryptedPrivateKey []byte) error {

	stmt := `
	UPDATE users 
	SET    
		encrypted_private_key = $2
	WHERE id = $1
	RETURNING 
	    id, 
   address,
   account_id,
   email, 
   username,
   first_name, 
   last_name, 
   bio, 
   website, 
   twitter, 
   discord, 
   avatar,
   public_key,
   encrypted_private_key,
   is_verified, 
   is_active, 
   kyc_registered_date,
   created_at,
   updated_at,
   user_type
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	updatedUser := models.User{}
	err := u.Db.QueryRowContext(
		ctx,
		stmt,
		userId,
		encryptedPrivateKey,
	).Scan(
		&updatedUser.Id,
		&updatedUser.Address,
		&updatedUser.AccountID,
		&updatedUser.Email,
		&updatedUser.Username,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Bio,
		&updatedUser.Website,
		&updatedUser.Twitter,
		&updatedUser.Discord,
		&updatedUser.Avatar,
		&updatedUser.PublicKey,
		&updatedUser.EncryptedPrivateKey,
		&updatedUser.IsVerified,
		&updatedUser.IsActive,
		&updatedUser.KYCRegisteredDate,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
		&updatedUser.UserType,
	)

	if err != nil {
		switch {
		case
			strings.Contains(err.Error(), duplicateAccountId),
			strings.Contains(err.Error(), duplicateEmail):
			return repository.ErrDuplicateDetails

		default:
			return err
		}
	}

	return nil
}

func (u userImplementation) Update(user models.User) (models.User, error) {

	stmt := `
	UPDATE users 
	SET    
	    email = $2,
	    first_name = $3,
	    last_name = $4,
	    account_id = $5,
		username = $6,
		avatar = $7,
		kyc_registered_date = $8,
		public_key = $9,
		bio = $10,
		website = $11,
		twitter = $12,
		discord = $13,
	    updated_at = now()
	WHERE id = $1
	RETURNING 
	    updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	updatedUser := user
	err := u.Db.QueryRowContext(
		ctx,
		stmt,
		user.Id,
		user.Email,
		user.FirstName,
		user.LastName,
		user.AccountID,
		user.Username,
		user.Avatar,
		user.KYCRegisteredDate,
		user.PublicKey,
		user.Bio,
		user.Website,
		user.Twitter,
		user.Discord,
	).Scan(
		&updatedUser.UpdatedAt,
	)

	if err != nil {

		switch {
		case
			strings.Contains(err.Error(), duplicateAccountId),
			strings.Contains(err.Error(), duplicateUsername),
			strings.Contains(err.Error(), duplicateEmail):
			return models.User{}, repository.ErrDuplicateDetails

		default:
			return models.User{}, err
		}
	}

	return updatedUser, nil
}

func (u userImplementation) Activate(id string) error {
	stmt := `
	UPDATE users
	SET is_active = $2
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := u.Db.ExecContext(ctx, stmt, id, true)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrRecordNotFound

		default:
			return err
		}
	}

	return nil
}
func (u userImplementation) Verify(id string) error {
	stmt := `
	UPDATE users
	SET is_verified = $2
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := u.Db.ExecContext(ctx, stmt, id, true)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrRecordNotFound

		default:
			return err
		}
	}

	return nil
}

func (u userImplementation) ListCreatedAssets(id string, filter models.Filter) ([]models.Asset, error) {
	stmt := `
	SELECT 
		a.id,
	   a.token_id,
	   a."name", 
	    a.symbol,
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
	   c."name",
	   a.is_verified,
		a.is_minted,
		a.is_rejected 
  FROM assets AS a
	   LEFT JOIN categories AS c ON a.category_id = c.id 
		JOIN users AS u ON a.creator_id = u.id
	WHERE 
		u.id = $1
	  	-- search by item name
	  	AND a.name ILIKE '%' || $2 || '%'

	-- reverse chronological order
	ORDER BY created_at DESC
	LIMIT $3 OFFSET $4;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	assets := []models.Asset{}

	rows, err := u.Db.QueryContext(
		ctx,
		stmt,
		id,
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
func (u userImplementation) ListOwnedAssets(id string, filter models.Filter) ([]models.Asset, error) {
	stmt := `
	SELECT 
    a.id,
    a.token_id,
    a."name", 
    a.symbol,
    a.category_id,
    a.blockchain, 
    a.creator_id,
    a.metadata_url, 
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
FROM 
    assets AS a
    LEFT JOIN categories AS c ON a.category_id = c.id 
    JOIN asset_owners AS ao ON a.id = ao.asset_id
    JOIN users AS u ON ao.user_id = u.id
WHERE 
    u.id = $1
    -- search by item name
    AND a.name ILIKE '%' || $2 || '%'
-- reverse chronological order
ORDER BY a.created_at DESC
LIMIT $3 OFFSET $4;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	assets := []models.Asset{}

	rows, err := u.Db.QueryContext(
		ctx,
		stmt,
		id,
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

func (u userImplementation) GetListedAssets(id string, filter models.Filter) ([]models.Asset, error) {
	stmt := `
	SELECT
	    a.id,
	   a.token_id,
	   a."name", 
	    a.symbol,
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
	    -- join users to get creator address
		JOIN users AS u ON u.id = a.creator_id
	    -- join categories to get contract address
		LEFT JOIN categories AS c ON a.category_id = c.id
		-- join listings to get price
		JOIN listings AS l ON l.asset_id = a.id

	WHERE u.id = $1
		-- get only active listings
		AND l.is_active = true
		-- search by asset name
	  	AND a.name ILIKE '%' || $2 || '%'

	-- reverse chronological order
	ORDER BY created_at DESC
	LIMIT $3 OFFSET $4;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	assets := []models.Asset{}

	rows, err := u.Db.QueryContext(
		ctx,
		stmt,
		id,
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
func (u userImplementation) GetUnListedAssets(id string, filter models.Filter) ([]models.Asset, error) {
	stmt := `
	SELECT
    a.id,
    a.token_id,
    a."name", 
    a.symbol,
    a.category_id,
    a.blockchain, 
    a.creator_id,
    a.metadata_url, 
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
    c."name",
	a.is_verified,
		a.is_minted,
		a.is_rejected 
FROM assets AS a
    -- join users to get creator address
    JOIN users AS u ON u.id = a.creator_id
    -- join categories to get contract address
    LEFT JOIN categories AS c ON a.category_id = c.id
    -- left join listings to check listing status
    LEFT JOIN listings AS l ON l.asset_id = a.id
WHERE u.id = $1
    AND (
        l.is_active = false OR 
        l.id IS NULL
    )
    -- search by asset name
    AND a.name ILIKE '%' || $2 || '%'
-- reverse chronological order
ORDER BY a.created_at DESC
LIMIT $3 OFFSET $4;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	assets := []models.Asset{}

	rows, err := u.Db.QueryContext(
		ctx,
		stmt,
		id,
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

func (u userImplementation) GetOrders(userId string, status *models.OrderStatus, orderType *models.OrderType, filter models.Filter) ([]models.Order, error) {
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
WHERE user_id=$1 and
($2::text IS NULL OR status = $2::text)
    AND ($3::text IS NULL OR "type" = $3::text)
 LIMIT $4 OFFSET $5;
		`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	orders := []models.Order{}

	rows, err := u.Db.QueryContext(
		ctx,
		stmt,
		userId,
		status,
		orderType,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return []models.Order{}, err
	}
	defer rows.Close()

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
		return []models.Order{}, err
	}

	return orders, nil
}

func (u userImplementation) HasOpenSellOrder(userId, assetId string) (bool, error) {
	stmt := `
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = $1
		AND asset_id = $2
		AND "type" = $3
		AND status IN ($4, $5)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var count int
	err := u.Db.QueryRowContext(ctx, stmt, userId, assetId, models.ORDER_SELL_TYPE, models.ORDER_OPEN_STATUS, models.ORDER_PARTIALLY_FILLED_STATUS).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
