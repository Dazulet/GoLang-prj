package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mangalib/auth-service/internal/utils"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func makeBody(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	return out
}

const testSecret = "test-jwt-secret"

func TestGenerateToken_Valid(t *testing.T) {
	token, err := utils.GenerateToken(42, "user", testSecret, 24)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := utils.ParseToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "user", claims.Role)
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, _ := utils.GenerateToken(1, "user", testSecret, 24)
	_, err := utils.ParseToken(token, "wrong-secret")
	assert.Error(t, err)
}

func TestParseToken_Expired(t *testing.T) {
	claims := utils.Claims{
		UserID: 99,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))

	_, err := utils.ParseToken(signed, testSecret)
	assert.Error(t, err, "expired token must be rejected")
}

func TestPasswordHashing(t *testing.T) {
	hash, err := utils.HashPassword("super-secret-123")
	require.NoError(t, err)
	assert.NotEqual(t, "super-secret-123", hash, "hash must not equal plaintext")
	assert.True(t, utils.CheckPassword(hash, "super-secret-123"))
	assert.False(t, utils.CheckPassword(hash, "wrong-password"))
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	router := gin.New()
	router.Use(authMiddlewareForTest(testSecret))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := gin.New()
	router.Use(authMiddlewareForTest(testSecret))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	token, _ := utils.GenerateToken(7, "user", testSecret, 24)

	router := gin.New()
	router.Use(authMiddlewareForTest(testSecret))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterHandler_ShortPassword(t *testing.T) {
	router := setupRegisterRouter()

	body := makeBody(map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"password": "short",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseResponse(t, w)
	assert.Equal(t, false, resp["success"])
}

func TestRegisterHandler_InvalidEmail(t *testing.T) {
	router := setupRegisterRouter()

	body := makeBody(map[string]string{
		"username": "bob",
		"email":    "not-an-email",
		"password": "ValidPass123",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler_EmptyBody(t *testing.T) {
	router := setupLoginRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateHandler_MissingToken(t *testing.T) {
	router := setupValidateRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/validate", makeBody(map[string]string{}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateHandler_ValidToken(t *testing.T) {
	token, _ := utils.GenerateToken(55, "admin", testSecret, 24)
	router := setupValidateRouterWithSecret(testSecret)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/validate",
		makeBody(map[string]string{"token": token}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.Equal(t, true, resp["success"])
}

func authMiddlewareForTest(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if len(header) < 8 || header[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
			c.Abort()
			return
		}
		_, err := utils.ParseToken(header[7:], secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func setupRegisterRouter() *gin.Engine {
	r := gin.New()
	r.POST("/api/auth/register", func(c *gin.Context) {
		var body struct {
			Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
			Email    string `json:"email"    binding:"required,email"`
			Password string `json:"password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"success": true})
	})
	return r
}

func setupLoginRouter() *gin.Engine {
	r := gin.New()
	r.POST("/api/auth/login", func(c *gin.Context) {
		var body struct {
			Email    string `json:"email"    binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
	return r
}

func setupValidateRouter() *gin.Engine {
	return setupValidateRouterWithSecret(testSecret)
}

func setupValidateRouterWithSecret(secret string) *gin.Engine {
	r := gin.New()
	r.POST("/api/auth/validate", func(c *gin.Context) {
		var body struct {
			Token string `json:"token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		claims, err := utils.ParseToken(body.Token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    gin.H{"user_id": claims.UserID, "role": claims.Role},
		})
	})
	return r
}
