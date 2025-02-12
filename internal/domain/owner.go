package domain

import "time"

type Owner struct {
	ID          uint      `db:"owner_id" json:"owner_id"`
	Name        string    `db:"name" json:"name"`
	Email       string    `db:"email" json:"email"`
	Phone       string    `db:"phone" json:"phone"`
	Password  	string    `db:"password" json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	// Relations
	Restaurants []Restaurant `json:"restaurants"`
}