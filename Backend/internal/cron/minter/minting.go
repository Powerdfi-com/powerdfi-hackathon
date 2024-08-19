package minter

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type Minter struct {
	app internal.Application
}

func NewMinter(app internal.Application) Minter {

	return Minter{
		app: app,
	}
}

// WatchDbForApprovedAssets continuously runs, checking and handling approved unminted assets.
func (i Minter) WatchDbForApprovedAssets() {

	for {
		assets, err := i.app.Repositories.Asset.GetApprovedUnmintedAssets(models.Filter{
			Limit: 10,
		})
		log.Printf("getting approved unminted assets: %v\n", len(assets))
		if err != nil {
			log.Printf("error getting expired assets: %s\n", err.Error())
			continue
		}

		for _, asset := range assets {

			urlsAsBytes := make([][]byte, 1)

			urlsAsBytes[0] = []byte(asset.MetadataUrl)
			// for i, metadata := range asset.MetadataUrl {
			// 	urlsAsBytes[i] = []byte(metadata)
			// }

			hederaTokenID, err := hedera.TokenIDFromString(asset.TokenId)
			if err != nil {

				log.Printf("error during mint %s: %s\n", asset.Id, err.Error())
				continue
			}
			assetOwner, err := i.app.Repositories.AssetOwner.GetOwnerAsset(
				asset.Id,
				asset.CreatorUserID,
			)

			if err != nil {
				if !errors.Is(err, repository.ErrRecordNotFound) {
					log.Printf("error getting asset ownership %s: %s\n", asset.Id, err.Error())
					continue
				}
			}

			serialNumbers := assetOwner.SerialNumbers

			for range asset.TotalSupply {

				nft, err := hedera.NewTokenMintTransaction().
					SetTokenID(hederaTokenID).
					SetMetadatas(urlsAsBytes).
					FreezeWith(i.app.HederaClient.Client)
				if err != nil {

					log.Printf("error during freeze mint txn %s: %s\n", asset.Id, err.Error())
					continue
				}

				// Sign token with treasury key
				resp, err := nft.Sign(*i.app.HederaClient.PrivateKey).
					Execute(i.app.HederaClient.Client)
				if err != nil {

					log.Printf("error during execute mint %s: %s\n", asset.Id, err.Error())
					continue
				}

				reciept, err := resp.GetReceipt(i.app.HederaClient.Client)
				if err != nil {

					if reciept.Status == hedera.StatusTokenMaxSupplyReached {
						log.Printf("token fully minted %s: %s\n", asset.Id, err.Error())
						break
					}
					log.Printf("error during get receipt minit %s: %s\n", asset.Id, err.Error())
					continue

				}

				creator, err := i.app.Repositories.User.GetById(asset.CreatorUserID)
				if err != nil {

					log.Printf("error during creator retrieval %s: %s\n", asset.Id, err.Error())
					continue
				}

				serialNumber := reciept.SerialNumbers[0]

				accountID, err := hedera.AccountIDFromString(creator.AccountID)
				if err != nil {

					log.Printf("error parse accountId %s: %s\n", asset.Id, err.Error())
					continue
				}

				fmt.Println("serialNumbers", serialNumber)
				fmt.Println("reciept", reciept.TokenID)

				tokenTransferTx, err := hedera.NewTransferTransaction().
					AddNftTransfer(hedera.NftID{TokenID: hederaTokenID, SerialNumber: serialNumber}, *i.app.HederaClient.AccountID, accountID).
					FreezeWith(i.app.HederaClient.Client)
				if err != nil {
					log.Printf("error during transfer txn %s: %s\n", asset.Id, err.Error())
					continue
				}
				// Sign with the treasury key to authorize the transfer
				signTransferTx := tokenTransferTx.Sign(*i.app.HederaClient.PrivateKey)

				tokenTransferSubmit, err := signTransferTx.Execute(i.app.HederaClient.Client)
				if err != nil {
					log.Printf("error during transfer execute %s: %s\n", asset.Id, err.Error())
					continue
				}

				_, err = tokenTransferSubmit.GetReceipt(i.app.HederaClient.Client)
				if err != nil {
					log.Printf("error during get receipt %s: %s\n", asset.Id, err.Error())
					continue
				}

				serialNumbers = append(serialNumbers, serialNumber)
				err = i.app.Repositories.AssetOwner.Update(models.AssetOwner{
					UserId:        creator.Id,
					AssetId:       asset.Id,
					SerialNumbers: serialNumbers,
				})
				if err != nil {
					log.Printf("error during get receipt %s: %s\n", asset.Id, err.Error())
					continue
				}
			}
			log.Printf("mint finished for %s\n", asset.Id)
			err = i.app.Repositories.Asset.UpdateMintStatus(asset.Id)
			if err != nil {
				log.Printf("error during update mint status %s: %s\n", asset.Id, err.Error())
				continue
			}
			_, err = i.app.Repositories.Activity.Add(models.Activity{
				AssetId: asset.Id,

				Action:   models.ACTIVITY_ACTION_MINT,
				ToUserId: asset.CreatorUserID,

				Quantity: asset.TotalSupply,
			})
			if err != nil {
				log.Printf("error during record transaction %s: %s\n", asset.Id, err.Error())
				continue
			}
		}

		// cooldown period of 1 minute
		time.Sleep(5 * time.Minute)
	}

}

func (c Minter) Mint() {

}
