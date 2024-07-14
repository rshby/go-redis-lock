package service

import (
	"context"
	"fmt"
	"github.com/rshby/go-redis-lock/http/httpresponse"
	"github.com/rshby/go-redis-lock/internal/cache"
	repositoryInterfaces "github.com/rshby/go-redis-lock/internal/repository/interfaces"
	"github.com/rshby/go-redis-lock/internal/service/dto"
	"github.com/rshby/go-redis-lock/internal/service/interfaces"
	"github.com/rshby/go-redis-lock/internal/utils/helper"
	"github.com/rshby/go-redis-lock/internal/utils/validatorutils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
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

// BurstStudentCount is method to burst student count times to demo lock
func (s *studentService) BurstStudentCount(ctx context.Context, count int) ([]dto.GetStudentResponseDTO, *httpresponse.HttpError) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"count":   count,
	})

	// validate count
	if err := validatorutils.Validate.Var(count, "gt=0"); err != nil {
		logger.Error(err)
		return nil, validatorutils.GetHttpErrorByTag(err)
	}

	var (
		wg          = &sync.WaitGroup{}
		chanStudent = make(chan dto.GetStudentResponseDTO, count)
		response    []dto.GetStudentResponseDTO
		countLocked int
	)

	// looping count times to create each student object
	wg.Add(1)
	go func(wg *sync.WaitGroup, count int) {
		defer wg.Done()

		var (
			wgStudent = &sync.WaitGroup{}
			epochTime = time.Now().Unix()
		)

		// looping count
		for i := 0; i < count; i++ {
			wgStudent.Add(1)
			go func(wgStudent *sync.WaitGroup) {
				defer wgStudent.Done()

				// lock with redis by epoch
				mutexUnlock, err := s.studentRepo.LockBurstStudentCountByEpoch(ctx, epochTime)
				defer mutexUnlock()
				if err != nil {
					// wait if another nodes are locked
					logrus.WithFields(logrus.Fields{
						"i":       i + 1,
						"lockKey": cache.BurstStudentCountLockKeyByEpoch(epochTime),
					}).Error(err)
				}

				// create student response
				student := dto.GetStudentResponseDTO{
					ID:        i + 1,
					FirstName: fmt.Sprintf("Student %d", i+1),
				}

				// send data to channel
				chanStudent <- student
			}(wgStudent)
		}

		// wait all goroutines student done
		wgStudent.Wait()

		// close channel
		close(chanStudent)
	}(wg, count)

	// append to response
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		// receive data from channel
		for student := range chanStudent {
			response = append(response, student)
		}
	}(wg)

	// wait
	wg.Wait()

	logger.Infof("total count response %d data. locked %d", len(response), countLocked)
	return response, nil
}
