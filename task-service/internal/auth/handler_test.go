package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockService implements ServiceInterface for handler testing
type mockService struct {
	registerResp *AuthResponse
	loginResp    *AuthResponse
	registerErr  error
	loginErr     error
}

func (m *mockService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	if m.registerErr != nil {
		return nil, m.registerErr
	}
	return m.registerResp, nil
}

func (m *mockService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	if m.loginErr != nil {
		return nil, m.loginErr
	}
	return m.loginResp, nil
}

func TestHandlerRegister(t *testing.T) {
	svc := &mockService{
		registerResp: &AuthResponse{
			Token: "test-token",
			User:  UserInfo{ID: "u1", Name: "Test", Email: "test@test.com", Role: "member"},
		},
	}
	handler := NewHandler(svc)

	body, _ := json.Marshal(RegisterRequest{Name: "Test", Email: "test@test.com", Password: "Test@1234"})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}

	var resp AuthResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Token != "test-token" {
		t.Errorf("expected token 'test-token', got '%s'", resp.Token)
	}
}

func TestHandlerRegisterInvalidBody(t *testing.T) {
	handler := NewHandler(&mockService{})

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestHandlerLogin(t *testing.T) {
	svc := &mockService{
		loginResp: &AuthResponse{
			Token: "login-token",
			User:  UserInfo{ID: "u1", Email: "test@test.com", Role: "member"},
		},
	}
	handler := NewHandler(svc)

	body, _ := json.Marshal(LoginRequest{Email: "test@test.com", Password: "Test@1234"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
