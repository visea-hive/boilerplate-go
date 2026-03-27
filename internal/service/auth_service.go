package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/visea-hive/auth-core/internal/model"
	"github.com/visea-hive/auth-core/internal/repository"
	"github.com/visea-hive/auth-core/internal/request"
	"github.com/visea-hive/auth-core/internal/response"
	"github.com/visea-hive/auth-core/pkg/crypto"
	"github.com/visea-hive/auth-core/pkg/jwt"
	"github.com/visea-hive/auth-core/pkg/mail"
	"github.com/visea-hive/auth-core/pkg/messages"
	"github.com/visea-hive/auth-core/pkg/validate"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication and registration logic.
type AuthService interface {
	CheckRegistrationRateLimit(ctx context.Context, ip string) error
	Register(ctx context.Context, req request.RegisterRequest, ip string) (response.RegisterResponse, error)
	VerifyEmail(ctx context.Context, token string, ip string, userAgent string) (response.TokenResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
	db       *gorm.DB
	rdb      *redis.Client
	hasher   *crypto.Hasher
	jwt      *jwt.Manager
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	db *gorm.DB,
	rdb *redis.Client,
	hasher *crypto.Hasher,
	jwt *jwt.Manager,
) AuthService {
	return &authService{
		userRepo: userRepo,
		authRepo: authRepo,
		db:       db,
		rdb:      rdb,
		hasher:   hasher,
		jwt:      jwt,
	}
}

func (s *authService) CheckRegistrationRateLimit(ctx context.Context, ip string) error {
	key := fmt.Sprintf("signups:%s", ip)
	count, err := s.rdb.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	if count == 1 {
		s.rdb.Expire(ctx, key, 1*time.Hour)
	}
	if count > 5 {
		return messages.ErrTooManyRequests
	}
	return nil
}

func (s *authService) Register(ctx context.Context, req request.RegisterRequest, ip string) (response.RegisterResponse, error) {
	// 1. Password Strength (Score >= 3)
	if err := validate.Strong(req.Password); err != nil {
		return response.RegisterResponse{}, err
	}

	// 3. Validate email availability
	existing, err := s.userRepo.GetByEmail(s.db, req.Email)
	if err != nil {
		return response.RegisterResponse{}, err
	}
	if existing != nil {
		return response.RegisterResponse{}, messages.ErrEmailAlreadyExists
	}

	// 4. Hash Password
	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return response.RegisterResponse{}, err
	}

	var userUUID string
	var rawToken string
	var tokenHashStr string

	err = s.db.Transaction(func(db *gorm.DB) error {
		// 5. Create User (pending)
		user := &model.User{
			Email:       req.Email,
			Password:    &passwordHash,
			DisplayName: &req.DisplayName,
			Status:      model.UserStatusPending,
		}
		user.UUID = model.GenerateUUIDv7()
		user.CreatedBy = user.UUID // creator is the user themselves during signup

		if err := s.userRepo.Create(db, user); err != nil {
			return err
		}
		userUUID = user.UUID

		// 6. Generate Verification Token
		var err error
		rawToken, tokenHashStr, err = crypto.NewToken()
		if err != nil {
			return err
		}

		verification := &model.EmailVerification{
			UserUUID:  user.UUID,
			TokenHash: tokenHashStr,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		if err := s.authRepo.CreateEmailVerification(db, verification); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return response.RegisterResponse{}, err
	}

	// 7. Send Email (Async)
	go func() {
		body := mail.VerifyEmailBody(req.DisplayName, fmt.Sprintf("http://localhost:8080/api/v1/auth/verify-email?token=%s", rawToken))
		if err := mail.Send(req.Email, "Verifikasi Email Anda", body); err != nil {
			slog.Error("Failed to send verification email", "email", req.Email, "error", err)
		}
	}()

	return response.RegisterResponse{
		UserUUID: userUUID,
		Message:  "Verify your email",
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, rawToken string, ip string, userAgent string) (response.TokenResponse, error) {
	tokenHashStr := crypto.HashToken(rawToken)

	verification, err := s.authRepo.GetEmailVerificationByTokenHash(s.db, tokenHashStr)
	if err != nil {
		return response.TokenResponse{}, err
	}

	if verification == nil {
		return response.TokenResponse{}, messages.ErrBadRequest
	}

	if time.Now().After(verification.ExpiresAt) {
		return response.TokenResponse{}, messages.ErrBadRequest
	}

	if verification.VerifiedAt != nil {
		return response.TokenResponse{}, nil
	}

	var tokenResp response.TokenResponse

	err = s.db.Transaction(func(db *gorm.DB) error {
		// 1. Update User
		user, err := s.userRepo.GetByUUID(db, verification.UserUUID)
		if err != nil {
			return err
		}
		now := time.Now()
		user.EmailVerifiedAt = &now
		user.Status = model.UserStatusActive
		if err := s.userRepo.Update(db, user); err != nil {
			return err
		}

		// 2. Update Verification record
		verification.VerifiedAt = &now
		if err := s.authRepo.UpdateEmailVerification(db, verification); err != nil {
			return err
		}

		// 3. Issue Token Pair
		tokenResp, err = s.issueTokenResponse(db, user, ip, userAgent)
		return err
	})

	if err != nil {
		return response.TokenResponse{}, err
	}

	return tokenResp, nil
}

// issueTokenResponse orchestrates the creation of a session, refresh token, and access token.
func (s *authService) issueTokenResponse(db *gorm.DB, user *model.User, ip, userAgent string) (response.TokenResponse, error) {
	// 1. Create Session
	session, err := s.createSession(db, user.UUID, ip, userAgent)
	if err != nil {
		return response.TokenResponse{}, err
	}

	// 2. Create Refresh Token
	rtValue, _, err := s.createRefreshToken(db, session.ID)
	if err != nil {
		return response.TokenResponse{}, err
	}

	// 3. Create Access Token (JWT)
	accessToken, err := s.createAccessToken(user, session.ID)
	if err != nil {
		return response.TokenResponse{}, err
	}

	return response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rtValue,
		ExpiresIn:    15 * 60, // 15 mins
	}, nil
}

// createAccessToken generates a JWT access token for the given user and session.
func (s *authService) createAccessToken(user *model.User, sessionID uint) (string, error) {
	// 1. Build Access Token (JWT) Claims
	claims := jwt.Claims{
		UserUUID:  user.UUID,
		Email:     user.Email,
		SessionID: fmt.Sprintf("%d", sessionID),
	}
	if user.DisplayName != nil {
		claims.Name = *user.DisplayName
	}
	if user.AvatarURL != nil {
		claims.AvatarURL = user.AvatarURL
	}

	// 2. Sign Access Token
	return s.jwt.Sign(claims)
}

// createSession creates a new session record in the database.
func (s *authService) createSession(db *gorm.DB, userUUID string, ip, userAgent string) (*model.Session, error) {
	session := &model.Session{
		UserUUID:  userUUID,
		IPAddress: &ip,
		UserAgent: &userAgent,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.authRepo.CreateSession(db, session); err != nil {
		return nil, err
	}
	return session, nil
}

// createRefreshToken generates a secure family and token, stores the hash, and returns the raw values.
func (s *authService) createRefreshToken(db *gorm.DB, sessionID uint) (string, string, error) {
	familyBytes, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return "", "", err
	}
	family := crypto.SHA256Hex(familyBytes)

	tokenBytes, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return "", "", err
	}
	tokenValue := crypto.SHA256Hex(tokenBytes)
	tokenHash := crypto.HashToken(tokenValue)

	refreshToken := &model.RefreshToken{
		SessionID: sessionID,
		TokenHash: tokenHash,
		Family:    family,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.authRepo.CreateRefreshToken(db, refreshToken); err != nil {
		return "", "", err
	}

	return tokenValue, family, nil
}
