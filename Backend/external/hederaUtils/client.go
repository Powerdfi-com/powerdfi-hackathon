package hederaUtils

import (
	"github.com/Powerdfi-com/Backend/config"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type Client struct {
	Client     *hedera.Client
	TokenID    *hedera.TokenID
	PrivateKey *hedera.PrivateKey
	AccountID  *hedera.AccountID
	Token      *Token
}

func CreateClient(cfg config.Config) (*Client, error) {

	var client *hedera.Client
	if cfg.Env == "production" {
		client = hedera.ClientForMainnet()
	} else {
		client = hedera.ClientForTestnet()
	}

	tokenID := cfg.Hedera.TokenId
	treasuryAccountID := cfg.Hedera.TreasuryAccountId
	treasuryPrivateKey := cfg.Hedera.TreasuryPrivateKey

	if treasuryAccountID == "" || treasuryPrivateKey == "" || tokenID == "" {
		return nil, ErrIncompleteTreasuryDetails
	}

	hederaTokenID, err := hedera.TokenIDFromString(tokenID)
	if err != nil {
		return nil, err
	}

	accountID, err := hedera.AccountIDFromString(treasuryAccountID)
	if err != nil {
		return nil, err
	}

	privateKey, err := hedera.PrivateKeyFromString(treasuryPrivateKey)
	if err != nil {
		return nil, err
	}

	client.SetOperator(accountID, privateKey)

	token := InitializeToken(client, &hederaTokenID, &privateKey, &accountID)

	return &Client{
		Client:     client,
		TokenID:    &hederaTokenID,
		PrivateKey: &privateKey,
		AccountID:  &accountID,
		Token:      token,
	}, nil
}
