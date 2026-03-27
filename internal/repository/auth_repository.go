package repository

import (
	"github.com/visea-hive/auth-core/internal/model"
	"gorm.io/gorm"
)

// AuthRepository defines the interface for auth-related database operations.
type AuthRepository interface {
	// Email Verification
	CreateEmailVerification(db *gorm.DB, ev *model.EmailVerification) error
	GetEmailVerificationByTokenHash(db *gorm.DB, tokenHash string) (*model.EmailVerification, error)
	UpdateEmailVerification(db *gorm.DB, ev *model.EmailVerification) error

	// Session & Refresh Token
	CreateSession(db *gorm.DB, session *model.Session) error
	GetSessionByID(db *gorm.DB, id uint) (*model.Session, error)
	CreateRefreshToken(db *gorm.DB, rt *model.RefreshToken) error
	GetRefreshTokenByTokenHash(db *gorm.DB, tokenHash string) (*model.RefreshToken, error)
}

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateEmailVerification(db *gorm.DB, ev *model.EmailVerification) error {
	return db.Create(ev).Error
}

func (r *authRepository) GetEmailVerificationByTokenHash(db *gorm.DB, tokenHash string) (*model.EmailVerification, error) {
	var ev model.EmailVerification
	if err := db.Where("token_hash = ?", tokenHash).First(&ev).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ev, nil
}

func (r *authRepository) UpdateEmailVerification(db *gorm.DB, ev *model.EmailVerification) error {
	return db.Save(ev).Error
}

func (r *authRepository) CreateSession(db *gorm.DB, session *model.Session) error {
	return db.Create(session).Error
}

func (r *authRepository) GetSessionByID(db *gorm.DB, id uint) (*model.Session, error) {
	var session model.Session
	if err := db.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *authRepository) CreateRefreshToken(db *gorm.DB, rt *model.RefreshToken) error {
	return db.Create(rt).Error
}

func (r *authRepository) GetRefreshTokenByTokenHash(db *gorm.DB, tokenHash string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	if err := db.Where("token_hash = ?", tokenHash).First(&rt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}
