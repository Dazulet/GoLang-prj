package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mangalib/manga-service/internal/client"
	"github.com/mangalib/manga-service/internal/utils"
)

// UserProxyHandler exposes Auth Service user data through the Manga Service.
// This demonstrates real Resty v2 inter-service communication:
//
//	Manga Service  --[Resty v2 HTTP GET]--> Auth Service
type UserProxyHandler struct {
	authClient *client.AuthClient
}

func NewUserProxyHandler(authClient *client.AuthClient) *UserProxyHandler {
	return &UserProxyHandler{authClient: authClient}
}

// GET /api/users/:id/profile
// Fetches user profile from Auth Service via Resty v2 and returns it.
// Useful for the frontend to get author info without calling Auth Service directly.
func (h *UserProxyHandler) GetUserProfile(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.authClient.GetUser(id)
	if err != nil {
		utils.NotFound(c, "user")
		return
	}

	utils.OK(c, user)
}
