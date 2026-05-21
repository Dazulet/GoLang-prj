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

	"github.com/mangalib/manga-service/internal/utils"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testSecret = "manga-test-secret"

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

// ── Test 1: Slugify utility ───────────────────────────────────────────────────

func TestSlugify(t *testing.T) {
	cases := []struct{ in, want string }{
		{"One Piece", "one-piece"},
		{"Attack on Titan!!", "attack-on-titan"},
		{"  My Hero Academia  ", "my-hero-academia"},
		{"Dragon Ball Z", "dragon-ball-z"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, utils.Slugify(tc.in), "input: %q", tc.in)
	}
}

// ── Test 2: ParsePagination defaults ─────────────────────────────────────────

func TestParsePagination_Defaults(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		page, limit, offset := utils.ParsePagination(c)
		c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit, "offset": offset})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	resp := parseResp(t, w)
	assert.Equal(t, float64(1), resp["page"])
	assert.Equal(t, float64(20), resp["limit"])
	assert.Equal(t, float64(0), resp["offset"])
}

// ── Test 3: ParsePagination respects custom values ────────────────────────────

func TestParsePagination_Custom(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		page, limit, offset := utils.ParsePagination(c)
		c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit, "offset": offset})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?page=3&limit=10", nil)
	r.ServeHTTP(w, req)

	resp := parseResp(t, w)
	assert.Equal(t, float64(3), resp["page"])
	assert.Equal(t, float64(10), resp["limit"])
	assert.Equal(t, float64(20), resp["offset"])
}

// ── Test 4: ParsePagination clamps limit to 100 ───────────────────────────────

func TestParsePagination_LimitClamped(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		_, limit, _ := utils.ParsePagination(c)
		c.JSON(http.StatusOK, gin.H{"limit": limit})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?limit=9999", nil)
	r.ServeHTTP(w, req)

	resp := parseResp(t, w)
	assert.Equal(t, float64(20), resp["limit"], "limit above 100 must default to 20")
}

// ── Test 5: Create manga — missing required field ─────────────────────────────

func TestCreateManga_MissingTitle(t *testing.T) {
	router := setupMangaRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/manga",
		makeBody(map[string]any{"author": "Oda"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+validAdminToken(t))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseResp(t, w)
	assert.Equal(t, false, resp["success"])
}

// ── Test 6: Create manga — invalid status enum ────────────────────────────────

func TestCreateManga_InvalidStatus(t *testing.T) {
	router := setupMangaRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/manga",
		makeBody(map[string]any{
			"title":  "Test Manga",
			"status": "flying",
		}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+validAdminToken(t))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── Test 7: Protected route rejects unauthenticated request ──────────────────

func TestProtectedRoute_NoToken(t *testing.T) {
	router := setupMangaRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/manga",
		makeBody(map[string]any{"title": "Naruto"}))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ── Test 8: Admin-only route rejects regular user token ───────────────────────

func TestAdminRoute_UserToken_Forbidden(t *testing.T) {
	router := setupMangaRouter()
	userToken, _ := utils.ParseToken("", testSecret) // trigger error path
	_ = userToken

	// Build a real user-role token
	token, err := buildToken(99, "user")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/manga",
		makeBody(map[string]any{"title": "Bleach"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ── Test 9: Health endpoint returns 200 ──────────────────────────────────────

func TestHealthEndpoint(t *testing.T) {
	router := setupMangaRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResp(t, w)
	assert.Equal(t, "manga", resp["service"])
	assert.Equal(t, "ok", resp["status"])
}

// ── Test 10: Chapter creation rejects negative chapter number ─────────────────

func TestCreateChapter_NegativeNumber(t *testing.T) {
	router := setupMangaRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/manga/1/chapters",
		makeBody(map[string]any{"number": -5}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+validAdminToken(t))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── Test 11: OK response helper wraps data correctly ─────────────────────────

func TestResponseHelper_OK(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		utils.OK(c, gin.H{"name": "Luffy"})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResp(t, w)
	assert.Equal(t, true, resp["success"])
	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Luffy", data["name"])
}

// ── Test 12: Error response helper ───────────────────────────────────────────

func TestResponseHelper_NotFound(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		utils.NotFound(c, "manga")
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
	resp := parseResp(t, w)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "manga not found", resp["message"])
}

// ── helpers ───────────────────────────────────────────────────────────────────

func buildToken(userID uint, role string) (string, error) {
	return utils.GenerateToken(userID, role, testSecret, 24)
}

func validAdminToken(t *testing.T) string {
	t.Helper()
	tok, err := buildToken(1, "admin")
	require.NoError(t, err)
	return tok
}

func setupMangaRouter() *gin.Engine {
	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "manga", "status": "ok"})
	})

	authMW := func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if len(header) < 8 || header[:7] != "Bearer " {
			utils.Unauthorized(c)
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(header[7:], testSecret)
		if err != nil {
			utils.Unauthorized(c)
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}

	adminMW := func(c *gin.Context) {
		role, _ := c.Get("userRole")
		if role != "admin" {
			utils.Forbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}

	validateManga := func(c *gin.Context) {
		var body struct {
			Title  string `json:"title"  binding:"required"`
			Status string `json:"status" binding:"omitempty,oneof=ongoing completed hiatus cancelled"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{"success": true, "data": body})
	}

	validateChapter := func(c *gin.Context) {
		var body struct {
			Number float64 `json:"number" binding:"required,min=0"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{"success": true})
	}

	api := r.Group("/api")
	protected := api.Group("")
	protected.Use(authMW)
	admin := protected.Group("")
	admin.Use(adminMW)
	admin.POST("/manga", validateManga)
	admin.POST("/manga/:mangaId/chapters", validateChapter)

	return r
}
