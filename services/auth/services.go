package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/auth/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo IAuthRepository
}

func NewAuthService(repo IAuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register hashes user password and stores the new user
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	user := &models.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Phone:     req.Phone,
		Role:      models.UserRole(req.Role),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login validates user credentials and issues Access & Refresh Tokens
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate Access Token
	accessToken, err := s.generateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshTokenString, err := s.generateSecureToken()
	if err != nil {
		return nil, err
	}

	expiryDaysStr := os.Getenv("REFRESH_TOKEN_EXPIRY_DAYS")
	expiryDays, err := strconv.Atoi(expiryDaysStr)
	if err != nil {
		expiryDays = 7 // Fallback
	}

	refreshTokenModel := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiredAt: time.Now().AddDate(0, 0, expiryDays),
		CreatedAt: time.Now(),
	}

	err = s.repo.CreateRefreshToken(ctx, refreshTokenModel)
	if err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		User: UserDTO{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

// RefreshToken validates a refresh token and generates a new access token
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AuthResponse, error) {
	rt, err := s.repo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(rt.ExpiredAt) {
		s.repo.DeleteRefreshToken(ctx, req.RefreshToken)
		return nil, errors.New("refresh token expired")
	}

	user, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := s.generateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rt.Token,
		User: UserDTO{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

// generateAccessToken creates a signed JWT string
func (s *AuthService) generateAccessToken(userID uuid.UUID, role string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		secret = []byte("super-secret-jwt-key")
	}

	expiryMinStr := os.Getenv("JWT_EXPIRY_MINUTES")
	expiryMin, err := strconv.Atoi(expiryMinStr)
	if err != nil {
		expiryMin = 15 // Fallback
	}

	claims := jwt.MapClaims{
		"userId": userID.String(),
		"role":   role,
		"exp":    time.Now().Add(time.Minute * time.Duration(expiryMin)).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token")
	}

	return tokenString, nil
}

// generateSecureToken creates a random hex string for Refresh Tokens
func (s *AuthService) generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
