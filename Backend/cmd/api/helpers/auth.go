package helpers

import (
	"bytes"
	"embed"
	"html/template"
	"strconv"
	"strings"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:embed "templates"
var templateFS embed.FS

// GenerateNonceMessage generates a random nonce with a welcome message if an account is to be created,
// and with a login message otherwise.
func GenerateNonceMessage(action string) (string, error) {
	uuidWithHyphen := uuid.New()
	nonce := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	// determine if the message is a welcome or a login prompt based on the given action
	templateFile := ""
	if action == models.NonceActionCreate {
		templateFile = "register_message.gotmpl"
	} else if action == models.NonceActionUpdate {
		templateFile = "login_message.gotmpl"
	}

	tmpl, err := template.New("message").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return "", err
	}

	// embed generated nonce into message
	message := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(message, "greeting", nonce)
	if err != nil {
		return "", err
	}

	return message.String(), nil
}

// RecoverAddressFromSignature retrieves the signer's address from a signature,
// given the original message.
func RecoverAddressFromSignature(message string, signature string) (string, error) {
	// hash the unsigned message using EIP-191 specifications.
	// source: https://eips.ethereum.org/EIPS/eip-191#specification
	prefix := "\x19Ethereum Signed Message:\n"
	dataLength := strconv.Itoa(len(message))
	formattedMessage := []byte(prefix + dataLength + message)
	messageHash := crypto.Keccak256Hash(formattedMessage)

	// get the bytes of the signature
	decodedMessage, err := hexutil.Decode(signature)
	if err != nil {
		return "", err
	}

	// handle cases where EIP-115 is not implemented (most wallets don't implement it)
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}

	// recover the public key from the signature
	publicKey, err := crypto.SigToPub(messageHash.Bytes(), decodedMessage)
	if err != nil {
		return "", err
	}

	// return the address gotten from the public jey
	address := crypto.PubkeyToAddress(*publicKey).String()
	return address, nil
}

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func MatchPassword(hash []byte, password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return false, err
	}

	return true, err
}
