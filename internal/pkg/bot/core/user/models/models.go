package models

import "encoding/json"

type User struct {
	Login     string  `json:"login"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Weight    float32 `json:"weight"`
	Height    uint    `json:"height"`
	Age       uint    `json:"age"`
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
