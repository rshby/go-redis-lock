package dto

import (
	"github.com/rshby/go-redis-lock/database/model"
)

type CreateStudentRequestDTO struct {
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"omitempty"`
	IdentityNumber string `json:"identity_number"`
	Email          string `json:"email" validate:"required,email"`
	Address        string `json:"address"`
}

// ToStudentEntity is method to convert object from dto to entity student
func (c *CreateStudentRequestDTO) ToStudentEntity() *model.Student {
	student := model.Student{
		FirstName:      c.FirstName,
		LastName:       c.LastName,
		IdentityNumber: c.IdentityNumber,
		Email:          c.Email,
		Address:        c.Address,
	}

	return &student
}
