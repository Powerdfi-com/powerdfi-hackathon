package models

var KYB_STATUS_FAILED = "failed"
var KYB_STATUS_SUCCESS = "success"
var KYB_STATUS_PENDING = "pending"

type UserKyB struct {
	UserId           string
	CertificateOfInc string
	Platform         string
	ReferenceId      string
	Status           string
	Comment          string
	CompanyName      string
	CompanyLocation  string
	CompanyAddress   string
	CompanyRegNo     string
}
