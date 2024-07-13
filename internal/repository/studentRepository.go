package repository

import (
	"context"
	"github.com/rshby/go-redis-lock/database/model"
	"github.com/rshby/go-redis-lock/internal/cache"
	cacheInterfaces "github.com/rshby/go-redis-lock/internal/cache/interfaces"
	"github.com/rshby/go-redis-lock/internal/config"
	interfaces "github.com/rshby/go-redis-lock/internal/repository/interfaces"
	"github.com/rshby/go-redis-lock/internal/utils/helper"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type studentRepository struct {
	db    *gorm.DB
	cache cacheInterfaces.CacheManager
}

// NewStudentRepository is function instance studentRepository
func NewStudentRepository(db *gorm.DB, cache cacheInterfaces.CacheManager) interfaces.StudentRepository {
	return &studentRepository{
		db:    db,
		cache: cache,
	}
}

// GetByID is method to get data student by id
func (s *studentRepository) GetByID(ctx context.Context, id int) (*model.Student, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"id":      id,
	})

	cacheKey := cache.GetStudentCacheKeyByID(id)
	if config.EnableCache() {
		cachedItem, err := cache.GetByKey[*model.Student](s.cache, cacheKey)
		if err == nil {
			if cachedItem != nil {
				logger.Info("get data student from redis cache")
				return cachedItem, nil
			}
		}
	}

	// get from mysql
	var student model.Student
	err := s.db.WithContext(ctx).Model(&model.Student{}).Take(&student, "id = ?", id).Error

	switch err {
	case nil:
		if config.EnableCache() {
			if err := s.cache.Set(cacheKey, helper.Dump(&student)); err != nil {
				logger.Error(err)
			}
		}

		return &student, nil
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}

// GetByEmail is method to get data student by email
func (s *studentRepository) GetByEmail(ctx context.Context, email string) (*model.Student, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"email":   email,
	})

	cacheKey := cache.GetStudentCacheKeyByEmail(email)
	if config.EnableCache() {
		cachedItem, err := cache.GetByKey[int](s.cache, cacheKey)
		if err == nil {
			if cachedItem > 0 {
				return s.GetByID(ctx, cachedItem)
			}
		}
	}

	// get from mysql
	var id int
	if err := s.db.WithContext(ctx).Model(&model.Student{}).Where("email = ?", email).Pluck("id", &id).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	// if not found
	if id == 0 {
		return nil, nil
	}

	if config.EnableCache() {
		if err := s.cache.Set(cacheKey, id); err != nil {
			logger.Error(err)
		}
	}

	return s.GetByID(ctx, id)
}

// Insert is function to insert new student
func (s *studentRepository) Insert(ctx context.Context, tx *gorm.DB, input *model.Student) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"input":   helper.Dump(input),
	})

	// using transaction
	if tx == nil {
		tx = s.db
	}

	// insert to mysql
	if err := tx.WithContext(ctx).Model(&model.Student{}).Create(input).Error; err != nil {
		logger.Error(err)
		return err
	}

	// delete related chache from redis
	if config.EnableCache() {
		if err := s.cache.DeleteByKeys([]string{
			cache.GetStudentCacheKeyByID(input.ID),
			cache.GetStudentCacheKeyByEmail(input.Email),
		}); err != nil {
			logger.Error(err)
		}
	}

	return nil
}
