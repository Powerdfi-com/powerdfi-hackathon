package shufti

var (
	EVENT_TYPE_ACCEPTED  string = "verification.accepted"
	EVENT_TYPE_DECLINED  string = "verification.declined"
	EVENT_TYPE_CANCELLED string = "verification.cancelled"
	EVENT_TYPE_INVALID   string = "request.invalid"
	EVENT_TYPE_TIMEOUT   string = "request.timeout"
)

/*
{
	"reference": "1234567",
	"event": "request.pending",
	"verification_url": "https://app.shuftipro.com/verification/process/tA8EP3JWgBjHeNpnE0vvS58m9IY6EaA1Xstveb9aMG51UfbWLAGzCJ0UTfvtp1ba",
	"email": "jhondeo@shufti.com",
	"country": null
}

*/

type Country struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type KycJourneyResponse struct {
	Reference       string `json:"reference"`
	Event           string `json:"event"`
	VerificationURL string `json:"verification_url"`
	Email           string `json:"email"`
	Country         string `json:"country"`
}

type EmailVerify struct {
	Email string `json:"email"`
}

type Name struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type Document struct {
	Name Name `json:"name"`
}
type VerificationData struct {
	EmailVerify EmailVerify `json:"email_verify"`
	Document    Document    `json:"document"`
}

type VerificationParams struct {
	ReferenceId      string           `json:"reference"`
	Event            string           `json:"event"`
	DeclinedReason   string           `json:"declined_reason"`
	VerificationData VerificationData `json:"verification_data"`
}
