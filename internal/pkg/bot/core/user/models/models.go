package models

type User struct {
	Login     string  `json:"login"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Weight    float32 `json:"weight"`
	Height    uint    `json:"height"`
	Age       uint    `json:"age"`
}
