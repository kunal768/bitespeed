package parser

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

type svc struct {
	repo Repository
}

type Service interface {
	ParseIncomingRequest(ctx *gin.Context, req ContactRequest) (ContactResponse, error)
}

func NewParserService(repo Repository) Service {
	return &svc{
		repo: repo,
	}
}

// collectContactDetails aggregates emails, phone numbers, and secondary IDs from a list of contacts.
func collectContactDetails(contacts []ContactData) ([]string, []string, []int) {
	emails := []string{}
	phones := []string{}
	secondaryIDs := []int{}
	for _, contact := range contacts {
		if contact.Email != "" {
			emails = append(emails, contact.Email)
		}
		if contact.PhoneNumber != "" {
			phones = append(phones, contact.PhoneNumber)
		}
		if contact.LinkPrecedence == "secondary" {
			secondaryIDs = append(secondaryIDs, contact.ID)
		}
	}
	return emails, phones, secondaryIDs
}

// insertNewPrimaryContact creates a new primary contact and returns the corresponding response.
func insertNewPrimaryContact(ctx *gin.Context, repo Repository, email, phone string) (ContactResponse, error) {
	newContact := ContactData{
		PhoneNumber:    phone,
		Email:          email,
		LinkPrecedence: "primary",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	insertedContact, err := repo.InsertContact(ctx, newContact)
	if err != nil {
		return ContactResponse{}, err
	}

	return ContactResponse{
		PrimaryContactID:    insertedContact.ID,
		Emails:              []string{email},
		PhoneNumbers:        []string{phone},
		SecondaryContactIDs: []int{},
	}, nil
}

// insertNewSecondaryContact creates a new secondary contact linked to a primary contact and returns the corresponding response.
func insertNewSecondaryContact(ctx *gin.Context, repo Repository, email, phone string, primaryContact ContactData) (ContactResponse, error) {
	newContact := ContactData{
		PhoneNumber:    phone,
		Email:          email,
		LinkPrecedence: "secondary",
		LinkedID:       &primaryContact.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	_, err := repo.InsertContact(ctx, newContact)
	if err != nil {
		return ContactResponse{}, err
	}

	// Query all contacts linked to the primary contact
	allContacts, err := repo.QueryContactsByLinkedID(ctx, primaryContact.ID)
	if err != nil {
		return ContactResponse{}, err
	}

	allContacts = append([]ContactData{primaryContact}, allContacts...)

	emails, phones, secondaryIDs := collectContactDetails(allContacts)

	return ContactResponse{
		PrimaryContactID:    primaryContact.ID,
		Emails:              emails,
		PhoneNumbers:        phones,
		SecondaryContactIDs: secondaryIDs,
	}, nil
}

// gatherContactResponse collects and returns the response based on the primary contact and its related contacts.
func gatherContactResponse(primaryContact ContactData, contacts []ContactData) (ContactResponse, error) {
	emails, phones, secondaryIDs := collectContactDetails(contacts)
	return ContactResponse{
		PrimaryContactID:    primaryContact.ID,
		Emails:              emails,
		PhoneNumbers:        phones,
		SecondaryContactIDs: secondaryIDs,
	}, nil
}

// updateSecondaryContact updates a secondary contact and returns the corresponding response.
func (s *svc) updateSecondaryContact(ctx *gin.Context, repo Repository, primaryContact, secondaryContact ContactData) (ContactData, error) {
	secondaryContact.LinkPrecedence = "secondary"
	secondaryContact.LinkedID = &primaryContact.ID
	return s.repo.UpdateContact(ctx, secondaryContact)
}

func (s *svc) ParseIncomingRequest(ctx *gin.Context, req ContactRequest) (ContactResponse, error) {
	if req.Email == "" && req.PhoneNumber == "" {
		return ContactResponse{}, errors.New("either email or phone number is required")
	}

	email, phone := req.Email, req.PhoneNumber
	givenPhone := len(phone) > 0
	givenEmail := len(email) > 0
	ctxStd := ctx
	if givenPhone && !givenEmail {
		contactsByPhone, err := s.repo.QueryContactsByPhoneNumber(ctxStd, phone)
		if err != nil {
			return ContactResponse{}, err
		}
		if len(contactsByPhone) == 0 {
			return ContactResponse{}, errors.New("no contacts found for provided phone number")
		}

		primaryContact := contactsByPhone[0]
		secondaryContacts, _ := s.repo.QueryContactsByLinkedID(ctx, primaryContact.ID)
		secondaryContacts = append([]ContactData{primaryContact}, secondaryContacts...)
		return gatherContactResponse(primaryContact, secondaryContacts)
	}

	if givenEmail && !givenPhone {
		contactsByEmail, err := s.repo.QueryContactsByEmail(ctxStd, email)
		if err != nil {
			return ContactResponse{}, err
		}
		if len(contactsByEmail) == 0 {
			return ContactResponse{}, errors.New("no contacts found for provided email")
		}
		var primaryContact ContactData
		primaryContact = contactsByEmail[0]
		if *primaryContact.LinkedID > 0 {
			contact, err := s.repo.QueryContactByID(ctx, *primaryContact.LinkedID)
			if err != nil {
				return ContactResponse{}, err
			}
			primaryContact = contact
		}

		secondaryContacts, _ := s.repo.QueryContactsByLinkedID(ctx, primaryContact.ID)
		secondaryContacts = append([]ContactData{primaryContact}, secondaryContacts...)
		return gatherContactResponse(primaryContact, secondaryContacts)
	}

	if givenEmail && givenPhone {
		contactsByEmail, err := s.repo.QueryContactsByEmail(ctxStd, email)
		if err != nil {
			return ContactResponse{}, err
		}
		contactsByPhone, err := s.repo.QueryContactsByPhoneNumber(ctxStd, phone)
		if err != nil {
			return ContactResponse{}, err
		}

		noEmails := len(contactsByEmail) == 0
		noPhones := len(contactsByPhone) == 0

		if noEmails && noPhones {
			return insertNewPrimaryContact(ctxStd, s.repo, email, phone)
		}

		if !noEmails && noPhones {
			primaryContact := contactsByEmail[0]
			return insertNewSecondaryContact(ctxStd, s.repo, email, phone, primaryContact)
		}

		if noEmails && !noPhones {
			primaryContact := contactsByPhone[0]
			return insertNewSecondaryContact(ctxStd, s.repo, email, phone, primaryContact)
		}

		if !noEmails && !noPhones {
			primaryContactEmail := contactsByEmail[0]
			primaryContactPhone := contactsByPhone[0]

			var primaryContact ContactData
			var secondaryContact ContactData
			if primaryContactEmail.CreatedAt.Before(primaryContactPhone.CreatedAt) {
				primaryContact = primaryContactEmail
				secondaryContact = primaryContactPhone
			} else {
				primaryContact = primaryContactPhone
				secondaryContact = primaryContactEmail
			}

			_, err := s.updateSecondaryContact(ctxStd, s.repo, primaryContact, secondaryContact)
			if err != nil {
				return ContactResponse{}, err
			}

			secondaryContacts, _ := s.repo.QueryContactsByLinkedID(ctx, primaryContact.ID)
			secondaryContacts = append([]ContactData{primaryContact}, secondaryContacts...)
			return gatherContactResponse(primaryContact, secondaryContacts)
		}
	}

	return ContactResponse{}, errors.New("case not handled")
}
