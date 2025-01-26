package domain

type Menu struct {
	ID           uint        `db:"menu_id" json:"menu_id"`
	RestaurantID uint        `db:"restaurant_id" json:"restaurant_id"`
	Categories   []Category  `json:"categories"`
}

type Category struct {
	ID          uint    `db:"category_id" json:"category_id"`
	MenuID      uint    `db:"menu_id" json:"menu_id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Items       []MenuItem `json:"items"`
}

type MenuItem struct {
	ID          uint    `db:"item_id" json:"item_id"`
	CategoryID  uint    `db:"category_id" json:"category_id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Price       float64 `db:"price" json:"price"`
	IsAvailable bool    `db:"is_available" json:"is_available"`
} 