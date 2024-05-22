package parser

type ContactRequest struct {
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type ContactResponse struct {
	PrimaryContactID    int      `json:"primaryContatctId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIDs []int    `json:"secondaryContactIds"`
}
