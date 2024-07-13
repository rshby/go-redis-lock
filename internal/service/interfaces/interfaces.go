package interfaces

import (
	"context"
	"github.com/rshby/go-redis-lock/http/httpresponse"
	"github.com/rshby/go-redis-lock/internal/service/dto"
)

type StudentService interface {
	GetByID(ctx context.Context, id int) (*dto.GetStudentResponseDTO, *httpresponse.HttpError)
	GetByEmail(ctx context.Context, email string) (*dto.GetStudentResponseDTO, *httpresponse.HttpError)
	CreateNewStudent(ctx context.Context, request *dto.CreateStudentRequestDTO) *httpresponse.HttpError
}
