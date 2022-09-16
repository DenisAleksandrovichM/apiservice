package models

import "encoding/json"

type User struct {
	Login     string  `db:"login" json:"login"`
	FirstName string  `db:"first_name" json:"first_name"`
	LastName  string  `db:"last_name" json:"last_name"`
	Weight    float32 `db:"weight" json:"weight"`
	Height    uint    `db:"height" json:"height"`
	Age       uint    `db:"age" json:"age"`
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

type KafkaMessage struct {
	Id   string `json:"id"`
	Data []byte `json:"data"`
}

func (k KafkaMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(k)
}

func (k *KafkaMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, k)
}
