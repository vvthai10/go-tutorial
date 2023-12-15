package entities

import "time"

type Account struct {
	ID       string    `bson:"_id,omitempty"`
	Owner    string    `bson:"owner"`
	Balance  int64     `bson:"balance"`
	Currency string    `bson:"currency"`
	CreateAt time.Time `bson:"createdAt"`
}