package hederaUtils

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

type Token struct {
	client     *hedera.Client
	privateKey *hedera.PrivateKey
	accountID  *hedera.AccountID
	tokenID    *hedera.TokenID
}

func InitializeToken(client *hedera.Client, tokenID *hedera.TokenID, privateKey *hedera.PrivateKey, accountID *hedera.AccountID) *Token {
	return &Token{
		client,
		privateKey,
		accountID,
		tokenID,
	}
}

func (t *Token) SetTokenID(tokenID *hedera.TokenID) {
	t.tokenID = tokenID
}

func (t *Token) SetPrivateKey(privateKey *hedera.PrivateKey) {
	t.privateKey = privateKey
}

func (t *Token) CreateToken(supplyKey hedera.Key, tokenType hedera.TokenType, decimalization uint, tokenName, tokenSymbol string, initialSupply uint64, signatures ...hedera.PrivateKey) (*hedera.TransactionReceipt, error) {
	token, err := hedera.NewTokenCreateTransaction().
		SetTokenName(tokenName).
		SetTokenType(tokenType).
		SetTokenSymbol(tokenSymbol).
		// TODO: Talk to Chiagozie.
		// We Are are limited to the amount of supply
		// when we use decimalization.
		SetDecimals(decimalization).
		SetInitialSupply(initialSupply).
		SetTreasuryAccountID(*t.accountID).
		// TODO: I am not sure how expiration works
		// What happens when a token expires?
		// How can we ensure tokens are not lost
		// Who renews?
		SetExpirationTime(time.Now().Add(time.Hour * 24 * 90)).
		// KYC key should be an admin generated key
		// that validates a user right after filling KYC
		// successfully.
		SetKycKey(t.privateKey.PublicKey()).
		SetSupplyKey(supplyKey).
		FreezeWith(t.client)

	if err != nil {
		return nil, err
	}

	// Add admin key
	token.Sign(*t.privateKey)
	// Add other signatures COO, CEO, etc
	for _, signature := range signatures {
		token.Sign(signature)
	}

	resp, err := token.Execute(t.client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(t.client)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func (t *Token) TransferFungibleToken(amount uint64, receiverAccountID *hedera.AccountID) (*hedera.TransactionReceipt, error) {
	token, err := hedera.NewTransferTransaction().
		AddTokenTransfer(*t.tokenID, *t.accountID, -int64(amount)).
		AddTokenTransfer(*t.tokenID, *receiverAccountID, int64(amount)).
		FreezeWith(t.client)

	if err != nil {
		return nil, err
	}

	resp, err := token.Sign(*t.privateKey).Execute(t.client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(t.client)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func (t *Token) ScheduleTransferNFT(nftID *hedera.NftID, receiver *hedera.AccountID, expirationTime time.Time) (*hedera.TransactionReceipt, error) {
	token, err := hedera.NewTransferTransaction().
		AddNftTransfer(*nftID, *t.accountID, *receiver).
		FreezeWith(t.client)
	if err != nil {
		return nil, err
	}

	scheduler, err := hedera.NewScheduleCreateTransaction().
		SetExpirationTime(expirationTime).
		SetScheduledTransaction(token)
	if err != nil {
		return nil, err
	}

	resp, err := scheduler.Execute(t.client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(t.client)
	return &receipt, err
}

func (t *Token) TransferNFToken(nftID *hedera.NftID, receiver *hedera.AccountID) (*hedera.TransactionReceipt, []byte, error) {
	token, err := hedera.NewTransferTransaction().
		AddNftTransfer(*nftID, *t.accountID, *receiver).
		FreezeWith(t.client)

	if err != nil {
		return nil, nil, err
	}

	resp, err := token.Sign(*t.privateKey).Execute(t.client)
	if err != nil {
		return nil, nil, err
	}

	receipt, err := resp.GetReceipt(t.client)
	if err != nil {
		return nil, nil, err
	}

	return &receipt, resp.Hash, nil
}

func (t *Token) ScheduleMintNFToken(urls []string, expirationTime time.Time) (*hedera.TransactionReceipt, error) {
	urlsAsBytes := make([][]byte, len(urls))

	for i, url := range urls {
		urlsAsBytes[i] = []byte(url)
	}

	nft, err := hedera.NewTokenMintTransaction().
		SetTokenID(*t.tokenID).
		SetMetadatas(urlsAsBytes).
		FreezeWith(t.client)
	if err != nil {
		return nil, err
	}

	scheduler, err := hedera.NewScheduleCreateTransaction().
		SetExpirationTime(expirationTime).
		SetScheduledTransaction(nft)
	if err != nil {
		return nil, err
	}

	resp, err := scheduler.Execute(t.client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(t.client)
	return &receipt, err
}

func (t *Token) SignScheduledTransaction(scheduleID string, signature *hedera.PrivateKey) (*hedera.TransactionReceipt, error) {
	scheduled, err := hedera.ScheduleIDFromString(scheduleID)
	if err != nil {
		return nil, err
	}

	toSign, err := hedera.NewScheduleSignTransaction().
		SetScheduleID(scheduled).
		FreezeWith(t.client)
	if err != nil {
		return nil, err
	}

	resp, err := toSign.Sign(*signature).Execute(t.client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(t.client)
	if err != nil {
		return nil, err
	}

	return &receipt, err
}

func (t *Token) QueryScheduledTransaction(scheduleID string) (*hedera.ScheduleInfo, error) {
	scheduled, err := hedera.ScheduleIDFromString(scheduleID)
	if err != nil {
		return nil, err
	}

	scheduleInfo, err := hedera.NewScheduleInfoQuery().
		SetScheduleID(scheduled).
		Execute(t.client)

	return &scheduleInfo, err
}

func (t *Token) MintNFToken(metadatas []string, signatories ...hedera.PrivateKey) (*hedera.TransactionReceipt, error) {
	urlsAsBytes := make([][]byte, len(metadatas))

	for i, metadata := range metadatas {
		urlsAsBytes[i] = []byte(metadata)
	}

	nft, err := hedera.NewTokenMintTransaction().
		SetTokenID(*t.tokenID).
		SetMetadatas(urlsAsBytes).
		FreezeWith(t.client)
	if err != nil {
		return nil, err
	}

	for _, signatory := range signatories {
		nft.Sign(signatory)
	}

	// Sign token with treasury key
	resp, err := nft.Sign(*t.privateKey).
		Execute(t.client)

	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(t.client)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func (t *Token) BurnNFToken(serialNo int64, cooKey *hedera.PrivateKey, ceoKey *hedera.PrivateKey) (*hedera.TransactionReceipt, error) {
	nft, err := hedera.NewTokenBurnTransaction().
		SetTokenID(*t.tokenID).
		SetSerialNumber(serialNo).
		FreezeWith(t.client)
	if err != nil {
		return nil, err
	}

	// Sign token with treasury key
	resp, err := nft.Sign(*t.privateKey).
		// Sign token with COO key
		Sign(*cooKey).
		// Sign token with CEO key
		Sign(*ceoKey).
		Execute(t.client)

	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(t.client)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}
