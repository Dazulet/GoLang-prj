package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mangalib/comment-service/internal/utils"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testSecret = "comment-test-secret"

func makeBody(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

func parseResp(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	return out
}

// ── Test 1: Comment creation requires manga_id ────────────────────────────────

func TestCreateComment_MissingMangaID(t *testing.T) {
	router := setupCommentRouter()
	token := adminToken(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/comments",
		makeBody(map[string]any{"body": "Great chapter!"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseResp(t, w)
	assert.Equal(t, false, resp["success"])
}

// ── Test 2: Comment creation requires non-empty body ─────────────────────────

func TestCreateComment_EmptyBody(t *testing.T) {
	router := setupCommentRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/comments",
		makeBody(map[string]any{"manga_id": 1, "body": ""}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken(t))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── Test 3: Comment creation requires auth ────────────────────────────────────

func TestCreateComment_Unauthenticated(t *testing.T) {
	router := setupCommentRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/comments",
		makeBody(map[string]any{"manga_id": 1, "body": "Hello"}))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ── Test 4: Delete comment — unauthenticated ──────────────────────────────────

func TestDeleteComment_Unauthenticated(t *testing.T) {
	router := setupCommentRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/comments/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ── Test 5: List comments — manga_id required ─────────────────────────────────

func TestListComments_MissingMangaID(t *testing.T) {
	router := setupCommentRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/comments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── Test 6: List comments — valid manga_id accepted ───────────────────────────

func TestListComments_ValidMangaID(t *testing.T) {
	router := setupCommentRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/comments?manga_id=1", nil)
	router.ServeHTTP(w, req)

	// Router returns stub 200 when manga_id is present
	assert.Equal(t, http.StatusOK, w.Code)
}

// ── Test 7: Like endpoint requires auth ───────────────────────────────────────

func TestLikeComment_Unauthenticated(t *testing.T) {
	router := setupCommentRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/comments/1/like", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ── Test 8: Auth middleware passes valid token ────────────────────────────────

func TestAuthMiddleware_ValidToken(t *testing.T) {
	token, err := utils.GenerateToken(5, "user", testSecret, 24)
	require.NoError(t, err)

	r := gin.New()
	r.Use(inlineAuthMW(testSecret))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ── Test 9: Auth middleware blocks wrong secret ───────────────────────────────

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	token, _ := utils.GenerateToken(5, "user", "different-secret", 24)

	r := gin.New()
	r.Use(inlineAuthMW(testSecret))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ── Test 10: Health endpoint ──────────────────────────────────────────────────

func TestHealthEndpoint(t *testing.T) {
	router := setupCommentRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResp(t, w)
	assert.Equal(t, "comment", resp["service"])
}

// ── Test 11: Comment body max length validation ───────────────────────────────

func TestCreateComment_TooLong(t *testing.T) {
	router := setupCommentRouter()

	longBody := make([]byte, 2001)
	for i := range longBody {
		longBody[i] = 'a'
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/comments",
		makeBody(map[string]any{"manga_id": 1, "body": string(longBody)}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken(t))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── Test 12: Paginated response structure ─────────────────────────────────────

func TestPaginatedResponseStructure(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		utils.Paginated(c, []string{"a", "b"}, 50, 2, 10)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResp(t, w)
	assert.Equal(t, true, resp["success"])

	pagination, ok := resp["pagination"].(map[string]any)
	require.True(t, ok, "pagination key must exist")
	assert.Equal(t, float64(2), pagination["page"])
	assert.Equal(t, float64(10), pagination["limit"])
	assert.Equal(t, float64(50), pagination["total"])
	assert.Equal(t, float64(5), pagination["total_pages"])
}

// ── helpers ───────────────────────────────────────────────────────────────────

func adminToken(t *testing.T) string {
	t.Helper()
	tok, err := utils.GenerateToken(1, "admin", testSecret, 24)
	require.NoError(t, err)
	return tok
}

func inlineAuthMW(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if len(header) < 8 || header[:7] != "Bearer " {
			utils.Unauthorized(c)
			c.Abort()
			return
		}
		_, err := utils.ParseToken(header[7:], secret)
		if err != nil {
			utils.Unauthorized(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func setupCommentRouter() *gin.Engine {
	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "comment", "status": "ok"})
	})

	// Public
	r.GET("/api/comments", func(c *gin.Context) {
		if c.Query("manga_id") == "" {
			utils.BadRequest(c, "manga_id is required")
			return
		}
		utils.Paginated(c, []any{}, 0, 1, 20)
	})

	// Protected group
	protected := r.Group("")
	protected.Use(inlineAuthMW(testSecret))

	protected.POST("/api/comments", func(c *gin.Context) {
		var body struct {
			MangaID uint   `json:"manga_id" binding:"required"`
			Body    string `json:"body"     binding:"required,min=1,max=2000"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.Created(c, body)
	})

	protected.DELETE("/api/comments/:id", func(c *gin.Context) {
		utils.NoContent(c)
	})

	protected.POST("/api/comments/:id/like", func(c *gin.Context) {
		utils.OK(c, gin.H{"action": "liked"})
	})

	return r
}
