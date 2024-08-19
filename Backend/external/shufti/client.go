package shufti

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Powerdfi-com/Backend/helpers"
)

//go:embed shufti_countries.json
var folder embed.FS

type Client struct {
	baseUrl   string
	clientId  string
	secretKey string
	journeyId string
	countries []Country
}

func NewShuftiClient(baseUrl, clientId, secretKey, journeyId string) *Client {

	// print(string(fileByte))
	content, _ := folder.ReadFile("shufti_countries.json")
	countries := []Country{}
	err := json.Unmarshal(content, &countries)
	if err != nil {
		log.Fatalf("err loading shufti supported countries list %s", err.Error())
	}

	return &Client{
		baseUrl,
		clientId,
		secretKey,
		journeyId,
		countries,
	}
}

func (h *Client) setRequestHeaders(req *http.Request) *http.Request {
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(h.clientId, h.secretKey)
	return req
}

func (h *Client) calculateResponseHash(responseData string) [32]byte {
	passwordHash := sha256.Sum256([]byte(h.secretKey))
	return sha256.Sum256([]byte(responseData + hex.EncodeToString(passwordHash[:])))
}

func (h *Client) IsValidSig(sig string, data []byte) bool {

	// Calculate the response hash
	responseHash := h.calculateResponseHash(string(data))

	// Check if the response signature matches the calculated hash
	return sig == hex.EncodeToString(responseHash[:])
}

func (h *Client) GetCountries() ([]Country, error) {
	return h.countries, nil
}
func (h *Client) GetJourneyVerificationLink(referenceId string, email string) (KycJourneyResponse, error) {
	reqUrl := h.baseUrl

	var kycResponse KycJourneyResponse
	reqBody, _ := json.Marshal(map[string]string{
		"journey_id": h.journeyId,
		"reference":  referenceId,
		"email":      email,
	})
	b := bytes.NewBuffer(reqBody)

	mRequest, err := http.NewRequest("POST", reqUrl, b)
	if err != nil {
		return kycResponse, err
	}
	mRequest = h.setRequestHeaders(mRequest)
	mResponse, err := http.DefaultClient.Do(mRequest)

	if err != nil {
		return kycResponse, err
	}

	defer mResponse.Body.Close()

	mBody, err := io.ReadAll(mResponse.Body)

	if err != nil {
		return kycResponse, err
	}

	if !helpers.IsStatusOk(mResponse.StatusCode) {
		switch mResponse.StatusCode {
		case http.StatusNotFound:
			return kycResponse, fmt.Errorf("not found")
		default:
			log.Println(string(mBody))
			return kycResponse, fmt.Errorf("err: code %v", mResponse.StatusCode)
		}
	}

	err = json.Unmarshal(mBody, &kycResponse)
	if err != nil {
		return kycResponse, err
	}
	return kycResponse, nil

}
func (h *Client) GetKycStatus(referenceId string) (VerificationParams, error) {
	// reqUrl := h.baseUrl
	reqUrl := fmt.Sprintf("%sstatus", h.baseUrl)

	var kycResponse VerificationParams
	reqBody, _ := json.Marshal(map[string]string{
		"journey_id": h.journeyId,
		"reference":  referenceId,
	})
	b := bytes.NewBuffer(reqBody)

	mRequest, err := http.NewRequest("POST", reqUrl, b)
	if err != nil {
		return kycResponse, err
	}
	mRequest = h.setRequestHeaders(mRequest)
	mResponse, err := http.DefaultClient.Do(mRequest)

	if err != nil {
		return kycResponse, err
	}

	defer mResponse.Body.Close()

	mBody, err := io.ReadAll(mResponse.Body)

	if err != nil {
		return kycResponse, err
	}

	if !helpers.IsStatusOk(mResponse.StatusCode) {
		switch mResponse.StatusCode {
		case http.StatusNotFound:
			return kycResponse, fmt.Errorf("not found")
		default:
			log.Println(string(mBody))
			return kycResponse, fmt.Errorf("err: code %v", mResponse.StatusCode)
		}
	}

	err = json.Unmarshal(mBody, &kycResponse)
	if err != nil {
		return kycResponse, err
	}
	return kycResponse, nil

}

func (h *Client) InitiateKYBVerification(referenceId, regNumber, country string) error {
	reqUrl := h.baseUrl

	reqBody, _ := json.Marshal(map[string]interface{}{
		"reference": referenceId,
		"kyb": map[string]interface{}{
			"advanced_search":             "1",
			"company_registration_number": regNumber,
			"country_names":               []string{country},
			"search_type":                 "contains",
		},
	})
	b := bytes.NewBuffer(reqBody)

	mRequest, err := http.NewRequest("POST", reqUrl, b)
	if err != nil {
		return err
	}
	mRequest = h.setRequestHeaders(mRequest)
	mResponse, err := http.DefaultClient.Do(mRequest)

	if err != nil {
		return err
	}

	defer mResponse.Body.Close()

	mBody, err := io.ReadAll(mResponse.Body)

	if err != nil {
		return err
	}

	if !helpers.IsStatusOk(mResponse.StatusCode) {
		switch mResponse.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("not found")
		default:
			log.Println(string(mBody))
			return fmt.Errorf("err: code %v", mResponse.StatusCode)
		}
	}

	return nil

}
