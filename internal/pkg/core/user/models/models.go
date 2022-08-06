package models

type User struct {
	Login     string  `db:"login"`
	FirstName string  `db:"first_name"`
	LastName  string  `db:"last_name"`
	Weight    float32 `db:"weight"`
	Height    uint    `db:"height"`
	Age       uint    `db:"age"`
}
