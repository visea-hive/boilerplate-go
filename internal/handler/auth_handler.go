package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visea-hive/auth-core/internal/request"
	"github.com/visea-hive/auth-core/internal/service"
	"github.com/visea-hive/auth-core/pkg/helpers"
	"github.com/visea-hive/auth-core/pkg/messages"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))
	ip := c.ClientIP()

	// 1. Rate Limiting Check (Fail fast before binding/validation)
	if err := h.authService.CheckRegistrationRateLimit(c.Request.Context(), ip); err != nil {
		if errors.Is(err, messages.ErrTooManyRequests) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": messages.Translate(lang, messages.ErrTooManyRequests)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer)})
		return
	}

	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := helpers.GenerateErrorValidationResponse(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": messages.Translate(lang, messages.ErrRegisterValidation), "errors": validationError.Errors})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req, ip)
	if err != nil {
		if errors.Is(err, messages.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"message": messages.Translate(lang, messages.ErrEmailAlreadyExists)})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": messages.Translate(lang, messages.ErrRegisterValidation), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	var req request.VerifyEmailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationError := helpers.GenerateErrorValidationResponse(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrGeneralInvalidInput), "errors": validationError.Errors})
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	resp, err := h.authService.VerifyEmail(c.Request.Context(), req.Token, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrBadRequest), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
