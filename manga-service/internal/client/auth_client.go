package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

type AuthClient struct {
	resty   *resty.Client
	baseURL string
}

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
