// Package client provides HTTP clients for communicating with other
// MangaLib microservices using the Resty v2 library.
package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// ── Shared response envelope ──────────────────────────────────────────────────

type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data"`
}

// ── UserInfo (matches auth-service public profile) ────────────────────────────

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

// ── AuthClient ────────────────────────────────────────────────────────────────

// AuthClient is a Resty v2-based HTTP client for the Auth Service.
// It is used by the Comment Service to fetch user profiles and validate tokens.
type AuthClient struct {
	client  *resty.Client
	baseURL string
}

// NewAuthClient constructs an AuthClient pointed at the given base URL.
// All requests time out after 5 seconds and retry once on transient failures.
func NewAuthClient(baseURL string) *AuthClient {
	r := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(5*time.Second).
		SetRetryCount(1).
		SetRetryWaitTime(500*time.Millisecond).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	return &AuthClient{client: r, baseURL: baseURL}
}

// GetUser fetches the public profile of a user from the Auth Service.
// Primary inter-service call: Comment Service -> Auth Service.
func (c *AuthClient) GetUser(userID uint) (*UserInfo, error) {
	var result apiResponse[UserInfo]

	resp, err := c.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/users/%d", userID))

	if err != nil {
		return nil, fmt.Errorf("auth service unreachable: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("user %d not found", userID)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("auth service error %d: %s", resp.StatusCode(), resp.String())
	}
	if !result.Success {
		return nil, fmt.Errorf("auth service: %s", result.Message)
	}

	return &result.Data, nil
}

// TokenClaims holds decoded JWT data returned by the Auth Service validate endpoint.
type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

// ValidateToken asks the Auth Service to validate a raw JWT string.
// Returns decoded user_id and role on success.
func (c *AuthClient) ValidateToken(token string) (*TokenClaims, error) {
	var result apiResponse[TokenClaims]

	resp, err := c.client.R().
		SetBody(map[string]string{"token": token}).
		SetResult(&result).
		Post("/api/auth/validate")

	if err != nil {
		return nil, fmt.Errorf("auth service unreachable: %w", err)
	}
	if resp.StatusCode() == 401 {
		return nil, fmt.Errorf("token invalid or expired")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("auth service error %d", resp.StatusCode())
	}

	return &result.Data, nil
}
