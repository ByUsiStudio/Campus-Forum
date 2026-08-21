package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Code    int    `json:"code"`    // HTTP状态码
	Error   string `json:"error"`   // 错误类型
	Message string `json:"message"` // 错误消息
	Detail  string `json:"detail"`  // 详细信息（可选）
}

// 常见错误类型
const (
	ErrTypeNotFound     = "NOT_FOUND"
	ErrTypeForbidden    = "FORBIDDEN"
	ErrTypeUnauthorized = "UNAUTHORIZED"
	ErrTypeBadRequest   = "BAD_REQUEST"
	ErrTypeInternal     = "INTERNAL_ERROR"
	ErrTypeRateLimit    = "RATE_LIMITED"
	ErrTypeValidation   = "VALIDATION_ERROR"
	ErrTypeConflict     = "CONFLICT"
)

// NotFound 404错误
func NotFound(c *gin.Context, message string, detail ...string) {
	errResponse := ErrorResponse{
		Code:    http.StatusNotFound,
		Error:   ErrTypeNotFound,
		Message: message,
	}
	if len(detail) > 0 {
		errResponse.Detail = detail[0]
	}
	c.JSON(http.StatusNotFound, errResponse)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string, detail ...string) {
	errResponse := ErrorResponse{
		Code:    http.StatusForbidden,
		Error:   ErrTypeForbidden,
		Message: message,
	}
	if len(detail) > 0 {
		errResponse.Detail = detail[0]
	}
	c.JSON(http.StatusForbidden, errResponse)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string, detail ...string) {
	errResponse := ErrorResponse{
		Code:    http.StatusUnauthorized,
		Error:   ErrTypeUnauthorized,
		Message: message,
	}
	if len(detail) > 0 {
		errResponse.Detail = detail[0]
	}
	c.JSON(http.StatusUnauthorized, errResponse)
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string, detail ...string) {
	errResponse := ErrorResponse{
		Code:    http.StatusBadRequest,
		Error:   ErrTypeBadRequest,
		Message: message,
	}
	if len(detail) > 0 {
		errResponse.Detail = detail[0]
	}
	c.JSON(http.StatusBadRequest, errResponse)
}

// InternalError 500错误
func InternalError(c *gin.Context, message string, detail ...string) {
	errResponse := ErrorResponse{
		Code:    http.StatusInternalServerError,
		Error:   ErrTypeInternal,
		Message: message,
	}
	if len(detail) > 0 {
		errResponse.Detail = detail[0]
	}
	c.JSON(http.StatusInternalServerError, errResponse)
}

// RateLimit 429错误
func RateLimit(c *gin.Context, message string, detail ...string) {
	errResponse := ErrorResponse{
		Code:    429,
		Error:   ErrTypeRateLimit,
		Message: message,
	}
	if len(detail) > 0 {
		errResponse.Detail = detail[0]
	}
	c.JSON(429, errResponse)
}

// Conflict 409错误
func Conflict(c *gin.Context, message string, detail ...string) {
	errResponse := ErrorResponse{
		Code:    http.StatusConflict,
		Error:   ErrTypeConflict,
		Message: message,
	}
	if len(detail) > 0 {
		errResponse.Detail = detail[0]
	}
	c.JSON(http.StatusConflict, errResponse)
}

// SendErrorResponse 通用错误响应函数
func SendErrorResponse(c *gin.Context, code int, message string, detail ...string) {
	var errType string

	switch code {
	case http.StatusNotFound:
		errType = ErrTypeNotFound
	case http.StatusForbidden:
		errType = ErrTypeForbidden
	case http.StatusUnauthorized:
		errType = ErrTypeUnauthorized
	case http.StatusBadRequest:
		errType = ErrTypeBadRequest
	case http.StatusInternalServerError:
		errType = ErrTypeInternal
	case 429:
		errType = ErrTypeRateLimit
	case http.StatusConflict:
		errType = ErrTypeConflict
	default:
		errType = ErrTypeInternal
	}

	errResponse := ErrorResponse{
		Code:    code,
		Error:   errType,
		Message: message,
	}
	if len(detail) > 0 {
		errResponse.Detail = detail[0]
	}
	c.JSON(code, errResponse)
}
