package token

import (
	"errors"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID 			int64 		`json:"id"`
	FullName 	string 		`json:"fullName"`
	Email 		string 		`json:"email"`
	IssuedAt 	time.Time 	`json:"issuedAt"`
	ExpiredAt 	time.Time 	`json:"expiredAt"`
}

type PayloadParams struct {
	ID 			int64 		`json:"id"`
	FullName 	string 		`json:"fullName"`
	Email 		string 		`json:"email"`
}

func NewPayload(arg PayloadParams, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		ID: arg.ID,
		FullName: arg.FullName,
		Email: arg.Email,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}