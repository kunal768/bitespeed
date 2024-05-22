package parser

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repo struct {
	client *pgxpool.Pool
}

type Repository interface {
	InsertContact(ctx *gin.Context, data ContactData) (ContactData, error)
	QueryContactsByPhoneNumber(ctx *gin.Context, phoneNumber string) ([]ContactData, error)
	QueryContactsByEmail(ctx *gin.Context, email string) ([]ContactData, error)
	UpdateContact(ctx *gin.Context, data ContactData) (ContactData, error)
	QueryContactsByLinkedID(ctx *gin.Context, linkedID int) ([]ContactData, error)
}

func NewRepository(client *pgxpool.Pool) Repository {
	return &repo{
		client: client,
	}
}

func (r *repo) InsertContact(ctx *gin.Context, data ContactData) (ContactData, error) {
	query := `
        INSERT INTO bitespeed.Contact (phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt, deletedAt
    `
	row := r.client.QueryRow(ctx, query, data.PhoneNumber, data.Email, data.LinkedID, data.LinkPrecedence, data.CreatedAt, data.UpdatedAt)

	var insertedData ContactData
	err := row.Scan(&insertedData.ID, &insertedData.PhoneNumber, &insertedData.Email, &insertedData.LinkedID, &insertedData.LinkPrecedence, &insertedData.CreatedAt, &insertedData.UpdatedAt, &insertedData.DeletedAt)
	if err != nil {
		return ContactData{}, fmt.Errorf("failed to insert contact: %v", err)
	}

	return insertedData, nil
}

func (r *repo) QueryContactsByPhoneNumber(ctx *gin.Context, phoneNumber string) ([]ContactData, error) {
	rows, err := r.client.Query(ctx, `
        SELECT id, phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt, deletedAt
        FROM bitespeed.Contact
        WHERE phoneNumber = $1
        ORDER BY createdAt ASC
    `, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query contacts by phone number: %v", err)
	}
	defer rows.Close()

	var contacts []ContactData
	for rows.Next() {
		var contact ContactData
		err := rows.Scan(&contact.ID, &contact.PhoneNumber, &contact.Email, &contact.LinkedID, &contact.LinkPrecedence, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contact row: %v", err)
		}
		contacts = append(contacts, contact)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %v", err)
	}

	return contacts, nil
}

func (r *repo) QueryContactsByEmail(ctx *gin.Context, email string) ([]ContactData, error) {
	rows, err := r.client.Query(ctx, `
        SELECT id, phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt, deletedAt
        FROM bitespeed.Contact
        WHERE email = $1
        ORDER BY createdAt ASC
    `, email)
	if err != nil {
		return nil, fmt.Errorf("failed to query contacts by email: %v", err)
	}
	defer rows.Close()

	var contacts []ContactData
	for rows.Next() {
		var contact ContactData
		err := rows.Scan(&contact.ID, &contact.PhoneNumber, &contact.Email, &contact.LinkedID, &contact.LinkPrecedence, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contact row: %v", err)
		}
		contacts = append(contacts, contact)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %v", err)
	}

	return contacts, nil
}

func (r *repo) UpdateContact(ctx *gin.Context, data ContactData) (ContactData, error) {
	query := `
        UPDATE bitespeed.Contact
        SET phoneNumber = $1, email = $2, linkedId = $3, linkPrecedence = $4, updatedAt = $5
        WHERE id = $6
        RETURNING id, phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt, deletedAt
    `
	row := r.client.QueryRow(ctx, query, data.PhoneNumber, data.Email, data.LinkedID, data.LinkPrecedence, time.Now(), data.ID)

	var updatedData ContactData
	err := row.Scan(&updatedData.ID, &updatedData.PhoneNumber, &updatedData.Email, &updatedData.LinkedID, &updatedData.LinkPrecedence, &updatedData.CreatedAt, &updatedData.UpdatedAt, &updatedData.DeletedAt)
	if err != nil {
		return ContactData{}, fmt.Errorf("failed to update contact: %v", err)
	}

	return updatedData, nil
}

// Repository implementation
func (r *repo) QueryContactsByLinkedID(ctx *gin.Context, linkedID int) ([]ContactData, error) {
	rows, err := r.client.Query(ctx, `
        SELECT id, phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt, deletedAt
        FROM bitespeed.Contact
        WHERE linkedId = $1
    `, linkedID)
	if err != nil {
		return nil, fmt.Errorf("failed to query contacts by linked ID: %v", err)
	}
	defer rows.Close()

	var contacts []ContactData
	for rows.Next() {
		var contact ContactData
		err := rows.Scan(&contact.ID, &contact.PhoneNumber, &contact.Email, &contact.LinkedID, &contact.LinkPrecedence, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contact row: %v", err)
		}
		contacts = append(contacts, contact)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %v", err)
	}

	return contacts, nil
}
