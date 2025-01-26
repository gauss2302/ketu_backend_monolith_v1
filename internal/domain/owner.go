package domain

type Owner struct {
	ID uint `db:"owner_id" json:"owner_id"`
	Company string `db:"company" json: "company"`
}