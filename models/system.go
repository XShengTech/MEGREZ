package models

type System struct {
	Key   string `json:"key" gorm:"primary_key;unique;index"`
	Value string `json:"value"`
}
