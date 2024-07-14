package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rshby/go-redis-lock/database"
	"github.com/rshby/go-redis-lock/internal/cache/interfaces"
	"github.com/rshby/go-redis-lock/internal/handler"
	"github.com/rshby/go-redis-lock/internal/repository"
	"github.com/rshby/go-redis-lock/internal/service"
	"gorm.io/gorm"
)

type AppRouter struct {
	db    *gorm.DB
	app   *gin.RouterGroup
	cache interfaces.CacheManager
}

// NewAppRouter is method to create new instance AppRouter
func NewAppRouter(app *gin.RouterGroup, cache interfaces.CacheManager) *AppRouter {
	return &AppRouter{
		db:    database.DatabaseMySQL,
		app:   app,
		cache: cache,
	}
}

func (r *AppRouter) InitEndpoint() {
	// register repository
	studentRepository := repository.NewStudentRepository(r.db, r.cache)

	// register service
	studentService := service.NewStudentService(r.db, studentRepository)

	// register handler
	studentHandler := handler.NewStudentHandler(studentService)

	// api v1
	apiV1 := r.app.Group("/v1")
	{
		studentGroup := apiV1.Group("/student")
		{
			studentGroup.GET("/:id", studentHandler.GetByID)
			studentGroup.POST("/", studentHandler.CreateNewStudent)
		}
	}
}
