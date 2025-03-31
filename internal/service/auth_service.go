package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lutestringamend/perwebbe/internal/config"
	"github.com/lutestringamend/perwebbe/internal/middleware"
	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, jwtConfig config.JWTConfig) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

func (s *authService) Register(username, email, password string) (*model.User, error) {
	existingUser, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingEmail, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		Active:       true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("email is not registered")
	}

	if !user.Active {
		return nil, errors.New("account is inactive")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := middleware.GenerateJWT(user.ID, user.Username, user.Role, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	authReponse := &model.AuthResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authReponse, nil
}

func (s *authService) RefreshToken(refreshToken string) (*model.AuthResponse, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	var userID uint
	_, err = fmt.Sscanf(sub, "%d", &userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if !user.Active {
		return nil, errors.New("user is inactive")
	}

	accessToken, err := middleware.GenerateJWT(user.ID, user.Username, user.Role, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := middleware.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	authReponse := &model.AuthResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	return authReponse, nil
}

func (s *authService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if ok && int64(exp) < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
