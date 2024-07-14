package service

import (
	"context"
	"github.com/rshby/go-redis-lock/http/httpresponse"
	repositoryInterfaces "github.com/rshby/go-redis-lock/internal/repository/interfaces"
	"github.com/rshby/go-redis-lock/internal/service/dto"
	"github.com/rshby/go-redis-lock/internal/service/interfaces"
	"github.com/rshby/go-redis-lock/internal/utils/helper"
	"github.com/rshby/go-redis-lock/internal/utils/validatorutils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type studentService struct {
	db          *gorm.DB
	studentRepo repositoryInterfaces.StudentRepository
}

// NewStudentService is function to create student service
func NewStudentService(db *gorm.DB, studentRepo repositoryInterfaces.StudentRepository) interfaces.StudentService {
	return &studentService{
		db:          db,
		studentRepo: studentRepo,
	}
}

// GetByID is method to get data student by id
func (s *studentService) GetByID(ctx context.Context, id int) (*dto.GetStudentResponseDTO, *httpresponse.HttpError) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"id":      id,
	})

	// call method in repository
	student, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, httpresponse.ErrorInternalServerError
	}

	if student == nil {
		return nil, httpresponse.ErrorStudentNotFound
	}

	// mapping to response
	response := dto.ConvertStudentEntityToStudentResponse(student)
	return response, nil
}

// GetByEmail is method to get data student by email
func (s *studentService) GetByEmail(ctx context.Context, email string) (*dto.GetStudentResponseDTO, *httpresponse.HttpError) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"email":   email,
	})

	// call method in repository
	student, err := s.studentRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.Error(err)
		return nil, httpresponse.ErrorInternalServerError
	}

	if student == nil {
		return nil, httpresponse.ErrorStudentNotFound
	}

	// mapping to response
	response := dto.ConvertStudentEntityToStudentResponse(student)
	return response, nil
}

// CreateNewStudent is method to create new student
func (s *studentService) CreateNewStudent(ctx context.Context, request *dto.CreateStudentRequestDTO) *httpresponse.HttpError {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"request": helper.Dump(request),
	})

	// Lock with redis by email
	mutexUnlock, err := s.studentRepo.LockCreateNewStudentByEmail(ctx, request.Email)
	defer mutexUnlock()
	if err != nil {
		logger.Error(err)
		return httpresponse.ErrorInternalServerError
	}

	// validate
	if err = validatorutils.Validate.Struct(*request); err != nil {
		logger.Error(err)
		return validatorutils.GetHttpErrorByTag(err)
	}

	// start transaction
	tx := s.db.Begin()
	defer tx.Rollback()

	// call method in repository
	if err := s.studentRepo.Insert(ctx, tx, request.ToStudentEntity()); err != nil {
		logger.Error(err)
		return httpresponse.ErrorInternalServerError
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		logger.Error(err)
		return httpresponse.ErrorInternalServerError
	}

	return nil
}
