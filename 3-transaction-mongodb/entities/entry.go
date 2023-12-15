package entities

import "time"

type Entry struct {
	ID        string    `bson:"_id,omitempty"`
	AccountID string    `bson:"accountID"`
	Amount    int64     `bson:"amount"`
	CreatedAt time.Time `bson:"createdAt"`
}