// Package client provides Resty v2-based HTTP clients for
// inter-service communication within the MangaLib platform.
package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// apiResponse is the standard envelope all MangaLib services return.
type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data"`
}

// ── DTOs ──────────────────────────────────────────────────────────────────────

// UserInfo is the public user profile returned by the Auth Service.
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

// TokenClaims holds the decoded JWT payload returned by /api/auth/validate.
type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

// ── AuthClient ────────────────────────────────────────────────────────────────

// AuthClient lets the Manga Service make HTTP calls to the Auth Service
// using the Resty v2 HTTP client library.
type AuthClient struct {
	resty   *resty.Client
	baseURL string
}

// NewAuthClient builds a configured Resty v2 client targeting the Auth Service.
//
// Configuration:
//   - 5s request timeout
//   - 1 automatic retry with 500ms wait (handles momentary restarts in Docker)
//   - JSON headers set by default
func NewAuthClient(baseURL string) *AuthClient {
	r := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(5*time.Second).
		SetRetryCount(1).
		SetRetryWaitTime(500*time.Millisecond).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	return &AuthClient{resty: r, baseURL: baseURL}
}

// GetUser fetches a public user profile from the Auth Service.
//
// Inter-service call: Manga Service → Auth Service GET /api/users/:id
// Used to attach author info to manga and bookmark responses.
func (c *AuthClient) GetUser(userID uint) (*UserInfo, error) {
	var result apiResponse[UserInfo]

	resp, err := c.resty.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/users/%d", userID))

	if err != nil {
		return nil, fmt.Errorf("auth service unreachable: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("user %d not found", userID)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("auth service returned %d", resp.StatusCode())
	}
	return &result.Data, nil
}

// ValidateToken calls the Auth Service token validation endpoint.
//
// Inter-service call: Manga Service → Auth Service POST /api/auth/validate
// Used in scenarios where the Manga Service needs server-side token verification
// (e.g., admin action auditing).
func (c *AuthClient) ValidateToken(token string) (*TokenClaims, error) {
	var result apiResponse[TokenClaims]

	resp, err := c.resty.R().
		SetBody(map[string]string{"token": token}).
		SetResult(&result).
		Post("/api/auth/validate")

	if err != nil {
		return nil, fmt.Errorf("auth service unreachable: %w", err)
	}
	if resp.StatusCode() == 401 {
		return nil, fmt.Errorf("invalid or expired token")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("auth service returned %d", resp.StatusCode())
	}
	return &result.Data, nil
}
