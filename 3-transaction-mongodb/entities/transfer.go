package entities

import "time"

type Transfer struct {
	ID            string    `bson:"_id,omitempty"`
	FromAccountID string    `bson:"fromAccountId"`
	ToAccountID   string    `bson:"toAccountId"`
	Amount        int64     `bson:"amount"`
	CreatedAt     time.Time `bson:"createdAt"`
}