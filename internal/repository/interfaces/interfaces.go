package interfaces

import (
	"context"
	"github.com/rshby/go-redis-lock/database/model"
	"gorm.io/gorm"
)

type StudentRepository interface {
	GetByID(ctx context.Context, id int) (*model.Student, error)
	GetByEmail(ctx context.Context, email string) (*model.Student, error)
	Insert(ctx context.Context, tx *gorm.DB, input *model.Student) error
}
