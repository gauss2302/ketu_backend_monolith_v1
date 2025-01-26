package domain

import "time"

type Restaurant struct {
	ID          uint      `db:"restaurant_id" json:"restaurant_id"`
	OwnerID     uint      `db:"owner_id" json:"owner_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	MainImage   string   `json:"main_image"`
	Images      []string `json:"images"`
	IsVerified  bool      `db:"is_verified" json:"is_verified"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	// Relations
	Details         RestaurantDetails  `json:"details"`
	Location RestaurantLocation `json:"location"`
	Menu            Menu              `json:"menu"`
}

type RestaurantLocation struct {
	ID           uint    `db:"location_id" json:"location_id"`
	RestaurantID uint    `db:"restaurant_id" json:"restaurant_id"`
	Address      Address  `db:"address" json:"address"`
	Latitude     float64 `db:"latitude" json:"latitude"`
	Longitude    float64 `db:"longitude" json:"longitude"`
} 


type Address struct {
	City     string `db:"city" json:"city"`
	District string `db:"district" json:"district"`
}

type RestaurantDetails struct {
	ID           uint    `db:"details_id" json:"details_id"`
	RestaurantID uint    `db:"restaurant_id" json:"restaurant_id"`
	Rating       float64 `db:"rating" json:"rating"`
	Capacity     int     `db:"capacity" json:"capacity"`
	OpeningHours string  `db:"opening_hours" json:"opening_hours"`
} 