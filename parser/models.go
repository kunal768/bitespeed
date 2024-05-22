package parser

import (
	"time"
)

type ContactData struct {
	ID             int        `json:"id"`
	PhoneNumber    string     `json:"phoneNumber"`
	Email          string     `json:"email"`
	LinkedID       *int       `json:"linkedId"`
	LinkPrecedence string     `json:"linkPrecedence"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
}
