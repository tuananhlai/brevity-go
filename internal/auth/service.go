package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/tuananhlai/brevity-go/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// Service defines authentication business logic.
type Service interface {
	Register(ctx context.Context, email, username, password string) error
	Login(ctx context.Context, emailOrUsername, password string) (*LoginResult, error)
	VerifyAccessToken(ctx context.Context, accessToken string) (string, error)
	GetCurrentUser(ctx context.Context, userID string) (*repository.User, error)
}

type serviceImpl struct {
	authRepo          repository.Repository
	accessTokenSecret string
	accessTokenExpiry time.Duration
	bcryptCost        int
}

func NewService(authRepo repository.Repository, accessTokenSecret string) Service {
	return &serviceImpl{
		bcryptCost:        bcrypt.DefaultCost,
		accessTokenExpiry: time.Hour * 24 * 30,
		authRepo:          authRepo,
		accessTokenSecret: accessTokenSecret,
	}
}

func (s *serviceImpl) Register(ctx context.Context, email, username, password string) error {
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
			return fmt.Errorf("%w: %s", repository.ErrUserAlreadyExists, err)
		}
		return err
	}

	return nil
}

// Login logs in a user and returns a JWT token.
func (s *serviceImpl) Login(ctx context.Context, emailOrUsername string, password string) (*LoginResult, error) {
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

	return &LoginResult{
		ID:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}

func (s *serviceImpl) GetCurrentUser(ctx context.Context, userID string) (*repository.User, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *serviceImpl) VerifyAccessToken(ctx context.Context, accessToken string) (string, error) {
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

	subject, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("missing token subject")
	}

	return subject, nil
}

func (s *serviceImpl) generateToken(userID string, secret string, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiry).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func (s *serviceImpl) generateAccessToken(userID string) (string, error) {
	return s.generateToken(userID, s.accessTokenSecret, s.accessTokenExpiry)
}

type LoginResult struct {
	ID          string
	Username    string
	Email       string
	AccessToken string
}
