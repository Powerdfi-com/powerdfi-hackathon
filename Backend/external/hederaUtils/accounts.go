package hederaUtils

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type Account struct {
	privateKey *hedera.PrivateKey
	publicKey  *hedera.PublicKey
	accountID  *hedera.AccountID

	client *hedera.Client
}

func InitializeAccount(client *hedera.Client, accountID *hedera.AccountID) *Account {
	return &Account{
		client:    client,
		accountID: accountID,
	}
}

func (t *Account) SetPrivateKey(privateKey *hedera.PrivateKey) {
	publicKey := privateKey.PublicKey()

	t.privateKey = privateKey
	t.publicKey = &publicKey
}

func CreateNewAccount(client *hedera.Client) (*Account, error) {
	hederaPrivKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		return nil, err
	}

	publicKey := hederaPrivKey.PublicKey()
	transaction := hedera.NewAccountCreateTransaction().
		SetKey(publicKey).
		SetAlias(publicKey.ToEvmAddress()).
		SetInitialBalance(hedera.NewHbar(0))

	//Sign the transaction with the client operator private key and submit to a Hedera network
	txResponse, err := transaction.Execute(client)
	if err != nil {
		panic(err)
	}

	//Request the receipt of the transaction
	receipt, err := txResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	//Get the account ID
	newAccountId := *receipt.AccountID

	return &Account{
		privateKey: &hederaPrivKey,
		publicKey:  &publicKey,
		accountID:  &newAccountId,
		client:     client,
	}, nil
}

func (a *Account) UnFreezeAccount(freezeKey hedera.PrivateKey, tokenID hedera.TokenID) (*hedera.TransactionReceipt, error) {
	freeze, err := hedera.NewTokenGrantKycTransaction().SetAccountID(*a.GetAccountID()).SetTokenID(tokenID).FreezeWith(a.client)
	if err != nil {
		return nil, err
	}

	resp, err := freeze.Sign(freezeKey).Execute(a.client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(a.client)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func (a *Account) FreezeAccount(freezeKey hedera.PrivateKey, tokenID hedera.TokenID) (*hedera.TransactionResponse, error) {
	freeze, err := hedera.NewTokenRevokeKycTransaction().SetAccountID(*a.GetAccountID()).SetTokenID(tokenID).FreezeWith(a.client)
	if err != nil {
		return nil, err
	}

	resp, err := freeze.Sign(freezeKey).Execute(a.client)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (a *Account) AssociateAccountToToken(tokenID hedera.TokenID) (*hedera.TransactionReceipt, error) {
	token, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(*a.GetAccountID()).
		AddTokenID(tokenID).
		FreezeWith(a.client)
	if err != nil {
		return nil, err
	}

	resp, err := token.Sign(*a.GetAccountPrivateKey()).Execute(a.client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(a.client)
	if err != nil {
		if receipt.Status != hedera.StatusTokenAlreadyAssociatedToAccount {
			return nil, err
		}

	}

	return &receipt, nil
}

func (a *Account) AllowTokenAssociation() (*hedera.TransactionReceipt, error) {
	account, err := hedera.NewAccountUpdateTransaction().
		SetAccountID(*a.accountID).
		FreezeWith(a.client)
	if err != nil {
		return nil, err
	}

	resp, err := account.Sign(*a.privateKey).Execute(a.client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(a.client)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func (a *Account) AddPrivateKey(privateKey *hedera.PrivateKey) {
	a.privateKey = privateKey
	pubKey := privateKey.PublicKey()
	a.publicKey = &pubKey
}

func (a *Account) GetAccountPrivateKey() *hedera.PrivateKey {
	return a.privateKey
}

func (a *Account) GetAccountPublicKey() *hedera.PublicKey {
	return a.publicKey
}

func (a *Account) GetAccountID() *hedera.AccountID {
	return a.accountID
}

func (a *Account) GetTokenBalanceWithID(tokenID hedera.TokenID) (uint64, error) {
	balanceMap, err := a.getTokenBalances()
	if err != nil {
		return 0, err
	}

	return balanceMap.Get(tokenID), nil
}

func (a *Account) GetTokenBalances() (*hedera.TokenBalanceMap, error) {
	return a.getTokenBalances()
}

func (a *Account) TransferFungibleToken(amount uint64, receiverAccountID *hedera.AccountID, tokenID *hedera.TokenID) (*hedera.TransactionReceipt, error) {
	return InitializeToken(a.client, tokenID, a.privateKey, a.accountID).
		TransferFungibleToken(amount, receiverAccountID)
}

func (a *Account) TransferNonFungibleToken(receiverAccountID *hedera.AccountID, nftID *hedera.NftID) (*hedera.TransactionReceipt, []byte, error) {
	return InitializeToken(a.client, &nftID.TokenID, a.privateKey, a.accountID).
		TransferNFToken(nftID, receiverAccountID)
}

func (a *Account) getTokenBalances() (*hedera.TokenBalanceMap, error) {
	r, err := hedera.NewAccountBalanceQuery().
		SetAccountID(*a.GetAccountID()).
		Execute(a.client)
	if err != nil {
		return nil, err
	}

	return &r.Tokens, nil
}

func (a *Account) GetAccountInfo() (hedera.AccountInfo, error) {
	return hedera.NewAccountInfoQuery().SetAccountID(*a.GetAccountID()).Execute(a.client)
}
