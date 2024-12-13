package domain

type Place struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Location    Location `json:"location"`
	MainImage   string   `json:"main_image"`
	Images      []string `json:"images"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type Location struct {
	Address  string `json:"address"`
	City     string `json:"city"`
	Province string `json:"province"`
}
