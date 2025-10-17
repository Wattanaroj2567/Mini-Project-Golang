package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gamegear/users-service/internal/models"
	"github.com/gamegear/users-service/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserAlreadyExists     = errors.New("user with this email or username already exists")
	ErrPasswordMismatch      = errors.New("passwords do not match")
	ErrTokenInvalidOrExpired = errors.New("token is invalid or has expired")
	ErrUnauthorizedRole      = errors.New("user does not have required role")
)

// AuthService defines the interface for authentication-related business logic.
type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
	RegisterAdmin(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	LoginAdmin(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	Logout(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, req models.ForgotPasswordRequest) error
	ForgotPasswordAdmin(ctx context.Context, req models.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error
	ResetPasswordAdmin(ctx context.Context, req models.ResetPasswordRequest) error
}

// authService is the implementation of the AuthService interface.
type authService struct {
	userRepo          repositories.UserRepository
	passwordResetRepo repositories.PasswordResetRepository
	tokenService      TokenService
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo repositories.UserRepository,
	passwordResetRepo repositories.PasswordResetRepository,
	tokenService TokenService,
) AuthService {
	return &authService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		tokenService:      tokenService,
	}
}

// Register orchestrates the user registration process for members.
func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	return s.registerWithRole(ctx, req, "member")
}

// RegisterAdmin creates a new administrator account.
func (s *authService) RegisterAdmin(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	return s.registerWithRole(ctx, req, "admin")
}

func (s *authService) registerWithRole(ctx context.Context, req models.RegisterRequest, role string) (*models.AuthResponse, error) {
	if req.Password != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	existingUser, err := s.userRepo.FindByEmailOrUsername(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	username, err := s.resolveUsername(ctx, req.Username, req.Email)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	user := &models.User{
		Username:     username,
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	token, err := s.tokenService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	user.PasswordHash = ""
	return &models.AuthResponse{User: user, Token: token}, nil
}

func (s *authService) resolveUsername(ctx context.Context, requested, email string) (string, error) {
	base := strings.TrimSpace(requested)
	if base == "" {
		if at := strings.Index(email, "@"); at > 0 {
			base = email[:at]
		} else {
			base = email
		}
	}
	if base == "" {
		base = "admin"
	}

	candidate := base
	suffix := 0
	for {
		existing, err := s.userRepo.FindByEmailOrUsername(ctx, candidate)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return candidate, nil
		}
		suffix++
		candidate = fmt.Sprintf("%s%d", base, suffix)
	}
}

// Login orchestrates the user login process.
func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	return s.loginWithRole(ctx, req, "")
}

// LoginAdmin authenticates an administrator.
func (s *authService) LoginAdmin(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	return s.loginWithRole(ctx, req, "admin")
}

func (s *authService) loginWithRole(ctx context.Context, req models.LoginRequest, requiredRole string) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmailOrUsername(ctx, req.Identifier)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if requiredRole != "" && user.Role != requiredRole {
		return nil, ErrUnauthorizedRole
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	user.PasswordHash = ""
	return &models.AuthResponse{User: user, Token: token}, nil
}

// Logout contains logic for user logout (e.g., token blacklisting if needed).
func (s *authService) Logout(ctx context.Context, token string) error {
	log.Printf("User logged out. Token revocation (if implemented) would apply to: %s", token)
	return nil
}

// ForgotPassword orchestrates the password reset initiation process.
func (s *authService) ForgotPassword(ctx context.Context, req models.ForgotPasswordRequest) error {
	return s.forgotPasswordWithRole(ctx, req, "")
}

// ForgotPasswordAdmin initiates password reset for administrators.
func (s *authService) ForgotPasswordAdmin(ctx context.Context, req models.ForgotPasswordRequest) error {
	return s.forgotPasswordWithRole(ctx, req, "admin")
}

func (s *authService) forgotPasswordWithRole(ctx context.Context, req models.ForgotPasswordRequest, requiredRole string) error {
	user, err := s.userRepo.FindByEmailOrUsername(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		log.Printf("Password reset requested for non-existent email: %s", req.Email)
		return nil
	}

	if requiredRole != "" && user.Role != requiredRole {
		return ErrUnauthorizedRole
	}

	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.passwordResetRepo.Create(ctx, resetToken); err != nil {
		return fmt.Errorf("could not create password reset token: %w", err)
	}

	// --- START: โค้ดที่ปรับปรุงใหม่ ---
	log.Println("!!! COPY THE TOKEN BELOW FOR THE NEXT STEP !!!")
	log.Printf("RESET_TOKEN_IS:%s", resetToken.Token)
	log.Println("!!! END OF TOKEN !!!")
	// --- END: โค้ดที่ปรับปรุงใหม่ ---

	return nil
}

// ResetPassword orchestrates the final step of the password reset process.
func (s *authService) ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error {
	return s.resetPasswordWithRole(ctx, req, "")
}

// ResetPasswordAdmin completes password reset for administrators.
func (s *authService) ResetPasswordAdmin(ctx context.Context, req models.ResetPasswordRequest) error {
	return s.resetPasswordWithRole(ctx, req, "admin")
}

func (s *authService) resetPasswordWithRole(ctx context.Context, req models.ResetPasswordRequest, requiredRole string) error {
	if req.NewPassword != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	resetToken, err := s.passwordResetRepo.FindByToken(ctx, req.Token)
	if err != nil {
		return ErrTokenInvalidOrExpired
	}

	if time.Now().After(resetToken.ExpiresAt) {
		_ = s.passwordResetRepo.DeleteByToken(ctx, req.Token)
		return ErrTokenInvalidOrExpired
	}

	user, err := s.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	if requiredRole != "" && user.Role != requiredRole {
		return ErrUnauthorizedRole
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash new password: %w", err)
	}

	user.PasswordHash = string(newHashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("could not update password: %w", err)
	}

	return s.passwordResetRepo.DeleteByToken(ctx, req.Token)
}
