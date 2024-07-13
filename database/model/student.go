package model

import "time"

type Student struct {
	ID             int       `json:"id" gorm:"column:id;type:int;not null;primaryKey;autoIncrement"`
	FirstName      string    `json:"first_name" gorm:"column:first_name;type:varchar(256);not null"`
	LastName       string    `json:"last_name" gorm:"column:last_name;type:varchar(256);default:null"`
	IdentityNumber string    `json:"identity_number" gorm:"column:identity_number;type:varchar(256);unique;default:null"`
	Email          string    `json:"email" gorm:"column:email;type:varchar(256);unique;not null"`
	Address        string    `json:"address" gorm:"column:address;type:text;default:null"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;not null;autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null;autoCreateTime;autoUpdateTime"`
}

func (s *Student) TableName() string {
	return "students"
}
