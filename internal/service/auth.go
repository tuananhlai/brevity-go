package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/tuananhlai/brevity-go/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type AuthService interface {
	Register(ctx context.Context, email, username, password string) error
	Login(ctx context.Context, emailOrUsername, password string) (*LoginReturn, error)
}

type TokenGenerationConfig struct {
	Secret string
	Expiry time.Duration
}

type authServiceImpl struct {
	authRepo repository.AuthRepository
	// accessTokenConfig is the configuration used for creating access tokens.
	accessTokenConfig TokenGenerationConfig
	// bcryptCost is the cost of the bcrypt hash. The larger this value, the more secure
	// the hash is, but the slower it is to generate.
	bcryptCost int
}

type AuthServiceOption func(*authServiceImpl)

func WithAccessTokenSecret(secret string) AuthServiceOption {
	return func(s *authServiceImpl) {
		s.accessTokenConfig.Secret = secret
	}
}

func WithAccessTokenExpiry(expiry time.Duration) AuthServiceOption {
	return func(s *authServiceImpl) {
		s.accessTokenConfig.Expiry = expiry
	}
}

func WithBcryptCost(cost int) AuthServiceOption {
	return func(s *authServiceImpl) {
		s.bcryptCost = cost
	}
}

func NewAuthService(authRepo repository.AuthRepository, opts ...AuthServiceOption) AuthService {
	defaultAccessTokenConfig := TokenGenerationConfig{
		// TODO: Replace with static secret.
		Secret: rand.Text(),
		Expiry: time.Hour * 24 * 30,
	}

	service := &authServiceImpl{
		authRepo:          authRepo,
		accessTokenConfig: defaultAccessTokenConfig,
		bcryptCost:        bcrypt.DefaultCost,
	}
	for _, opt := range opts {
		opt(service)
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

func (s *authServiceImpl) generateToken(userID string, secret string, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiry).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func (s *authServiceImpl) generateAccessToken(userID string) (string, error) {
	return s.generateToken(userID, s.accessTokenConfig.Secret, s.accessTokenConfig.Expiry)
}

type LoginReturn struct {
	ID          string
	Username    string
	Email       string
	AccessToken string
}
