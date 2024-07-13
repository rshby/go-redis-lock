package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rshby/go-redis-lock/http/httpresponse"
	"github.com/rshby/go-redis-lock/internal/service/interfaces"
	"github.com/rshby/go-redis-lock/internal/utils/helper"
	"github.com/sirupsen/logrus"
)

type StudentHandler struct {
	studentService interfaces.StudentService
}

// NewStudentHandler is method to create instance studentService
func NewStudentHandler(studentService interfaces.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// UploadStudentsData is handler to upload student data
func (s *StudentHandler) UploadStudentsData(ctx *gin.Context) {
	// TODO : buat controller untuk upload students data csv
}

// GetByID is method handler to get data student by id
func (s *StudentHandler) GetByID(ctx *gin.Context) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
	})

	// get params
	id := helper.ExpectNumber[int](ctx.Param("id"))

	// call method in service
	student, httpError := s.studentService.GetByID(ctx, id)
	if httpError != nil {
		logger.Error(httpError)
		httpresponse.ResponseError(ctx, httpError)
		return
	}

	httpresponse.ResponseOK(ctx, "success get data student", student)
	return
}
