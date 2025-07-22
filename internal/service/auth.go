package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type AuthService interface {
	Register(ctx context.Context, email, username, password string) error
	Login(ctx context.Context, emailOrUsername, password string) (*LoginReturn, error)
	VerifyAccessToken(ctx context.Context, accessToken string) (string, error)
	GetCurrentUser(ctx context.Context, userID string) (*model.AuthUser, error)
}

type authServiceImpl struct {
	authRepo          repository.AuthRepository
	accessTokenSecret string
	accessTokenExpiry time.Duration
	// bcryptCost is the cost of the bcrypt hash. The larger this value, the more secure
	// the hash is, but the slower it is to generate.
	bcryptCost int
}

func NewAuthService(authRepo repository.AuthRepository, accessTokenSecret string) AuthService {
	service := &authServiceImpl{
		bcryptCost:        bcrypt.DefaultCost,
		accessTokenExpiry: time.Hour * 24 * 30,
		authRepo:          authRepo,
		accessTokenSecret: accessTokenSecret,
	}

	return service
}

func (s *authServiceImpl) Register(ctx context.Context, email, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return err
	}

	_, err = s.authRepo.CreateUser(ctx, repository.CreateUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return fmt.Errorf("%w: %s", ErrUserAlreadyExists, err)
		}
		return err
	}

	return nil
}

// Login logs in a user and returns a JWT token.
func (s *authServiceImpl) Login(ctx context.Context, emailOrUsername string, password string) (*LoginReturn, error) {
	user, err := s.authRepo.GetUser(ctx, emailOrUsername)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidCredentials, err)
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCredentials, err)
	}

	accessToken, err := s.generateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &LoginReturn{
		ID:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}

func (s *authServiceImpl) GetCurrentUser(ctx context.Context, userID string) (*model.AuthUser, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authServiceImpl) VerifyAccessToken(ctx context.Context, accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(s.accessTokenSecret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims["sub"].(string), nil
}

func (s *authServiceImpl) generateToken(userID string, secret string, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiry).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func (s *authServiceImpl) generateAccessToken(userID string) (string, error) {
	return s.generateToken(userID, s.accessTokenSecret, s.accessTokenExpiry)
}

type LoginReturn struct {
	ID          string
	Username    string
	Email       string
	AccessToken string
}
