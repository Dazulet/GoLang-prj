package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mangalib/manga-service/internal/client"
	"github.com/mangalib/manga-service/internal/utils"
)

type UserProxyHandler struct {
	authClient *client.AuthClient
}

func NewUserProxyHandler(authClient *client.AuthClient) *UserProxyHandler {
	return &UserProxyHandler{authClient: authClient}
}

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
