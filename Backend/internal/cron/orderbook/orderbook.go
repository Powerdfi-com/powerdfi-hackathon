package orderbook

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Powerdfi-com/Backend/external/hederaUtils"
	utils "github.com/Powerdfi-com/Backend/helpers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type OrderBook struct {
	app internal.Application
}

func NewOrderBook(app internal.Application) OrderBook {
	return OrderBook{
		app: app,
	}
}

func (o OrderBook) WatchDbOrderMatches() {

	log.Printf("Starting order matching service")
	for {
		// Check for matching buy and sell orders
		buyOrder, sellOrder := o.findMatchingOrders()
		if buyOrder != nil && sellOrder != nil {
			// Perform the transfer of assets and USDC tokens
			// log.Printf("fulfilling sell order %s and buy order id %s", sellOrder.Id, buyOrder.Id)
			// TODO: prevent fufilling own user orders
			err := o.fulfillOrders(buyOrder, sellOrder)

			if err != nil {
				log.Printf("err fuffilling sell order %s and buy order id %s %s", sellOrder.Id, buyOrder.Id, err.Error())

				continue
				// Handle error
			}

			// log.Fatalf("test completed")

		}

		// Sleep for a certain duration before checking again
		time.Sleep(1 * time.Minute)
	}
}

func (o OrderBook) markOrderAsFailed(order *models.Order) {
	order.Status = models.ORDER_FAILED_STATUS
	err := o.app.Repositories.Order.Update(*order)
	if err != nil {
		panic(err)
	}
}

func (o OrderBook) fulfillOrders(buyOrder, sellOrder *models.Order) error {
	// Check if partial fulfillment is required
	quantityToFill := buyOrder.Quantity
	if quantityToFill > sellOrder.Quantity {
		quantityToFill = sellOrder.Quantity
	}

	// Transfer assets and tokens for the determined quantity
	return o.transferAssetsAndTokens(buyOrder, sellOrder, quantityToFill)
}

func (o OrderBook) findMatchingOrders() (*models.Order, *models.Order) {
	limit := 100 // Adjust this value based on your system capabilities
	var page = 0

	for {
		// Get a batch of open orders from the database

		openOrders, err := o.app.Repositories.Order.GetUnFilledBuyOrders(models.Filter{
			Limit: limit,
			Page:  page,
		})

		if err != nil {
			log.Fatalf("err getting open orders %s", err.Error())
		}

		log.Printf("found %d open orders", len(openOrders))
		// Check for matching orders within the current batch
		for _, order := range openOrders {
			if order.Type == models.ORDER_BUY_TYPE {
				sellOrder, err := o.app.Repositories.Order.FindMatchingSellOrder(order)

				if err != nil {
					if errors.Is(err, repository.ErrRecordNotFound) {
						continue
					}
					log.Printf("err finding matching sell order for buy order id %s %s", order.Id, err.Error())
				}

				return &order, &sellOrder
			}
		}

		// Check if more orders need to be fetched
		if len(openOrders) < limit {
			break // Reached the end of open orders
		}

		// Increment offset for the next batch
		page += 1
	}

	return nil, nil
}

func (o OrderBook) transferAssetsAndTokens(buyOrder, sellOrder *models.Order, quantityToFill int64) error {
	// Implement logic to transfer assets and USDC tokens between users
	// Get the authenticated user

	// Get the listing from the database

	fmt.Println("quantity", quantityToFill)
	// Get the asset associated with the listing
	asset, err := o.app.Repositories.Asset.GetById(sellOrder.AssetId)
	if err != nil {

		return err
	}

	// Get the buyer and creator accounts
	buyer, err := o.app.Repositories.User.GetById(buyOrder.UserId)
	if err != nil {

		return err
	}
	seller, err := o.app.Repositories.User.GetById(sellOrder.UserId)
	if err != nil {

		return err
	}

	// Decrypt the user's private key
	data, err := utils.DecryptPlainText(buyer.EncryptedPrivateKey, o.app.Config.MasterKey)
	if err != nil {

		return err
	}

	// Convert the decrypted data to a private key
	key, err := hedera.PrivateKeyFromBytesECDSA(data)
	if err != nil {

		return err
	}

	// Decrypt the creator's private key
	cdata, err := utils.DecryptPlainText(seller.EncryptedPrivateKey, o.app.Config.MasterKey)
	if err != nil {

		return err
	}

	// Convert the decrypted data to a private key
	ckey, err := hedera.PrivateKeyFromBytesECDSA(cdata)
	if err != nil {
		return err
	}

	// Convert the account IDs to Hedera account IDs
	accountID, err := hedera.AccountIDFromString(buyer.AccountID)
	if err != nil {
		return err
	}
	creatorID, err := hedera.AccountIDFromString(seller.AccountID)
	if err != nil {
		return err
	}

	// Convert the token IDs to Hedera token IDs
	assetTokenId, err := hedera.TokenIDFromString(asset.TokenId)
	if err != nil {
		return err
	}
	usdcTokenId, err := hedera.TokenIDFromString(o.app.Config.Hedera.TokenIdUSDC)
	if err != nil {
		return err
	}

	// Initialize the user buyerAccount and set the private key
	buyerAccount := hederaUtils.InitializeAccount(o.app.HederaClient.Client, &accountID)
	buyerAccount.SetPrivateKey(&key)

	// Get the balance of the user's account
	balance, err := buyerAccount.GetTokenBalanceWithID(usdcTokenId)
	if err != nil {
		return err
	}

	// Convert the balance to Hbar units
	parsedBalance := hedera.HbarFrom(float64(balance), hedera.HbarUnits.Microbar)

	fmt.Println("balance", parsedBalance.As(hedera.HbarUnits.Hbar), sellOrder.Price*float64(quantityToFill))
	// Check if the balance is sufficient for the listing price
	if parsedBalance.As(hedera.HbarUnits.Hbar) < sellOrder.Price*float64(quantityToFill) {

		o.markOrderAsFailed(buyOrder)
		return fmt.Errorf("insufficient balance")
	}

	// Initialize the creator account and set the private key
	sellerAccount := hederaUtils.InitializeAccount(o.app.HederaClient.Client, &creatorID)
	sellerAccount.SetPrivateKey(&ckey)

	// Get the asset ownership for the creator and user
	sellerAssetOwnership, err := o.app.Repositories.AssetOwner.GetOwnerAsset(asset.Id, asset.CreatorUserID)
	if err != nil {
		o.markOrderAsFailed(sellOrder)

		return err
	}
	buyerAssetOwnership, err := o.app.Repositories.AssetOwner.GetOwnerAsset(asset.Id, buyer.Id)
	if err != nil {
		if !errors.Is(err, repository.ErrRecordNotFound) {
			o.markOrderAsFailed(sellOrder)
			return err
		}

	}

	// Check if the creator has enough serial numbers for the requested quantity
	if len(sellerAssetOwnership.SerialNumbers) < int(quantityToFill) {
		o.markOrderAsFailed(sellOrder)
		return err
	}

	// Transfer USDC tokens from the user to the creator
	_, err = buyerAccount.TransferFungibleToken(uint64(hedera.HbarFrom(sellOrder.Price*float64(quantityToFill), hedera.HbarUnits.Hbar).As(hedera.HbarUnits.Microbar)), &creatorID, &usdcTokenId)
	if err != nil {
		return err
	}

	// Associate the user's account with the asset token
	_, err = buyerAccount.AssociateAccountToToken(assetTokenId)
	if err != nil {
		return err
	}

	// Transfer the requested quantity of non-fungible tokens from the creator to the user
	for range quantityToFill {
		serialNumber := sellerAssetOwnership.SerialNumbers[0]
		_, _, err = sellerAccount.TransferNonFungibleToken(&accountID, &hedera.NftID{
			TokenID:      assetTokenId,
			SerialNumber: serialNumber,
		})
		if err != nil {
			return err
		}

		buyerAssetOwnership.SerialNumbers = append(buyerAssetOwnership.SerialNumbers, serialNumber)
		sellerAssetOwnership.SerialNumbers = sellerAssetOwnership.SerialNumbers[1:]

		err = o.app.Repositories.AssetOwner.Update(sellerAssetOwnership)
		if err != nil {
			return err
		}

		err = o.app.Repositories.AssetOwner.Update(buyerAssetOwnership)
		if err != nil {
			return err
		}
	}
	// Update the status of the buy and sell orders
	if quantityToFill == sellOrder.Quantity {
		// The sell order was completely filled
		sellOrder.Status = models.ORDER_FILLED_STATUS
	} else {
		// The sell order was partially filled
		sellOrder.Status = models.ORDER_PARTIALLY_FILLED_STATUS
	}

	if quantityToFill == buyOrder.Quantity {
		// The buy order was completely filled
		buyOrder.Status = models.ORDER_FILLED_STATUS
	} else {
		// The buy order was partially filled
		buyOrder.Status = models.ORDER_PARTIALLY_FILLED_STATUS

	}

	buyOrder.Quantity -= quantityToFill
	sellOrder.Quantity -= quantityToFill

	// Update the buy and sell orders in the database
	err = o.app.Repositories.Order.Update(*buyOrder)
	if err != nil {
		return err
	}

	err = o.app.Repositories.Order.Update(*sellOrder)
	if err != nil {
		return err
	}

	_, err = o.app.Repositories.Activity.Add(models.Activity{
		AssetId:    asset.Id,
		FromUserId: &sellOrder.UserId,
		ToUserId:   buyOrder.UserId,
		Action:     models.ACTIVITY_ACTION_SALE,
		Price:      sellOrder.Price,
		Currency:   string(internal.USDC),
		Quantity:   quantityToFill,
	})
	if err != nil {
		return err
	}

	log.Printf("Order fulfilled successfully")
	return nil
}
