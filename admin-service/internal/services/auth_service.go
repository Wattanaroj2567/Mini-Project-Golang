package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gamegear/admin-service/internal/clients"
	"github.com/gamegear/admin-service/internal/models"
)

// AuthService coordinates admin authentication workflows.
type AuthService interface {
	Register(ctx context.Context, req models.AdminRegisterRequest) (*models.AdminAuthResponse, error)
	Login(ctx context.Context, req models.AdminLoginRequest) (*models.AdminAuthResponse, error)
	ForgotPassword(ctx context.Context, req models.AdminForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req models.AdminResetPasswordRequest) error
	Logout(ctx context.Context, token string) error
}

type authService struct {
	authClient clients.AuthClient
}

// NewAuthService constructs AuthService.
func NewAuthService(authClient clients.AuthClient) AuthService {
	return &authService{authClient: authClient}
}

func (s *authService) Register(ctx context.Context, req models.AdminRegisterRequest) (*models.AdminAuthResponse, error) {
	resp, err := s.authClient.RegisterAdmin(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("auth service: register admin: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("auth service: register admin: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to register admin")
		return nil, newServiceError(resp.StatusCode, message)
	}

	authResp, err := decodeAdminAuthResponse(body)
	if err != nil {
		return nil, fmt.Errorf("auth service: decode register response: %w", err)
	}

	return authResp, nil
}

func (s *authService) Login(ctx context.Context, req models.AdminLoginRequest) (*models.AdminAuthResponse, error) {
	identifier := strings.TrimSpace(req.Email)
	if identifier == "" {
		identifier = strings.TrimSpace(req.Identifier)
	}
	if identifier == "" {
		return nil, newServiceError(http.StatusBadRequest, "email or identifier is required")
	}

	payload := map[string]string{
		"identifier": identifier,
		"password":   req.Password,
	}

	resp, err := s.authClient.LoginAdmin(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("auth service: login admin: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("auth service: login admin: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to login admin")
		return nil, newServiceError(resp.StatusCode, message)
	}

	authResp, err := decodeAdminAuthResponse(body)
	if err != nil {
		return nil, fmt.Errorf("auth service: decode login response: %w", err)
	}

	return authResp, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req models.AdminForgotPasswordRequest) error {
	resp, err := s.authClient.ForgotPassword(ctx, req)
	if err != nil {
		return fmt.Errorf("auth service: forgot-password: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return fmt.Errorf("auth service: forgot-password: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to submit forgot-password request")
		return newServiceError(resp.StatusCode, message)
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req models.AdminResetPasswordRequest) error {
	resp, err := s.authClient.ResetPassword(ctx, req)
	if err != nil {
		return fmt.Errorf("auth service: reset-password: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return fmt.Errorf("auth service: reset-password: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to reset password")
		return newServiceError(resp.StatusCode, message)
	}

	return nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return newServiceError(http.StatusBadRequest, "authorization token is required")
	}

	resp, err := s.authClient.LogoutAdmin(ctx, token)
	if err != nil {
		return fmt.Errorf("auth service: logout admin: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return fmt.Errorf("auth service: logout admin: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to logout admin")
		return newServiceError(resp.StatusCode, message)
	}

	return nil
}

func decodeAdminAuthResponse(body []byte) (*models.AdminAuthResponse, error) {
	var direct models.AdminAuthResponse
	if err := json.Unmarshal(body, &direct); err == nil && !isZeroAdminProfile(direct.Admin) {
		return &direct, nil
	}

	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Data) > 0 {
		return mapAuthPayload(envelope.Data)
	}

	return mapAuthPayload(body)
}

func mapAuthPayload(payload []byte) (*models.AdminAuthResponse, error) {
	type upstreamProfile struct {
		ID          uint   `json:"id"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}
	type upstreamAuth struct {
		Token string           `json:"token"`
		User  *upstreamProfile `json:"user"`
		Admin *upstreamProfile `json:"admin"`
	}

	var upstream upstreamAuth
	if err := json.Unmarshal(payload, &upstream); err != nil {
		return nil, fmt.Errorf("decode auth payload: %w", err)
	}

	profile := models.AdminProfile{}
	switch {
	case upstream.Admin != nil:
		profile = models.AdminProfile{
			ID:          upstream.Admin.ID,
			Email:       upstream.Admin.Email,
			DisplayName: upstream.Admin.DisplayName,
		}
	case upstream.User != nil:
		profile = models.AdminProfile{
			ID:          upstream.User.ID,
			Email:       upstream.User.Email,
			DisplayName: upstream.User.DisplayName,
		}
	}

	resp := &models.AdminAuthResponse{
		Token: upstream.Token,
		Admin: profile,
	}
	if isZeroAdminProfile(resp.Admin) {
		return nil, fmt.Errorf("decode auth payload: missing user/admin data")
	}

	return resp, nil
}

func isZeroAdminProfile(profile models.AdminProfile) bool {
	return profile.ID == 0 && profile.Email == "" && profile.DisplayName == ""
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return body, nil
}

func parseErrorMessage(body []byte, fallback string) string {
	if len(bytes.TrimSpace(body)) == 0 {
		return fallback
	}

	var payload struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	if err := json.Unmarshal(body, &payload); err == nil {
		switch {
		case payload.Message != "":
			return payload.Message
		case payload.Error != "":
			return payload.Error
		}
	}

	return string(body)
}
