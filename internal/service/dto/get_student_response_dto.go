package dto

import (
	"github.com/rshby/go-redis-lock/database/model"
	"github.com/rshby/go-redis-lock/internal/utils/helper"
)

type GetStudentResponseDTO struct {
	ID             int    `json:"id,omitempty" `
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	IdentityNumber string `json:"identity_number,omitempty"`
	Email          string `json:"email,omitempty"`
	Address        string `json:"address,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdateAt       string `json:"update_at,omitempty"`
}

// ConvertStudentEntityToStudentResponse is function to convert entity student to response student
func ConvertStudentEntityToStudentResponse(input *model.Student) *GetStudentResponseDTO {
	student := GetStudentResponseDTO{
		ID:             input.ID,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		IdentityNumber: input.IdentityNumber,
		Email:          input.Email,
		Address:        input.Address,
		CreatedAt:      helper.TimeToStringFormat(input.CreatedAt),
		UpdateAt:       helper.TimeToStringFormat(input.UpdatedAt),
	}

	return &student
}
