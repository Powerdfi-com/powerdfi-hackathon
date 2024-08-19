package models

var KYC_STATUS_FAILED = "failed"
var KYC_STATUS_SUCCESS = "success"

type UserKyc struct {
	UserId      string
	URL         string
	Platform    string
	ReferenceId string
	Status      string
	Comment     string
}
