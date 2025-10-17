package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/gamegear/users-service/internal/models"
	"github.com/gamegear/users-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrOldPasswordIncorrect = errors.New("old password is not correct")
	ErrPasswordRequired     = errors.New("password confirmation is required for this action")
)

// ProfileService defines the interface for profile management business logic.
type ProfileService interface {
	GetProfile(ctx context.Context, userID uint) (*models.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req models.UpdateProfileRequest) (*models.ProfileResponse, error)
	DeleteAccount(ctx context.Context, userID uint, password string) error
}

// profileService is the implementation of the ProfileService interface.
type profileService struct {
	userRepo repositories.UserRepository
}

// NewProfileService creates a new ProfileService.
func NewProfileService(userRepo repositories.UserRepository) ProfileService {
	return &profileService{userRepo: userRepo}
}

// GetProfile retrieves a user's profile information.
func (s *profileService) GetProfile(ctx context.Context, userID uint) (*models.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &models.ProfileResponse{
		ID:           user.ID,
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
	}

	return response, nil
}

// UpdateProfile orchestrates the logic for updating profile info, changing the password, or deleting the account.
func (s *profileService) UpdateProfile(ctx context.Context, userID uint, req models.UpdateProfileRequest) (*models.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Case 2: Change password
	if req.NewPassword != "" {
		if req.OldPassword == "" || req.ConfirmPassword == "" {
			return nil, errors.New("old password and confirmation are required to set a new password")
		}
		if req.NewPassword != req.ConfirmPassword {
			return nil, ErrPasswordMismatch
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
			return nil, ErrOldPasswordIncorrect
		}
		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash new password: %w", err)
		}
		user.PasswordHash = string(newHashedPassword)
	}

	// Case 3: Update profile info
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return s.GetProfile(ctx, userID)
}

// DeleteAccount removes a user after verifying the password.
func (s *profileService) DeleteAccount(ctx context.Context, userID uint, password string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}
