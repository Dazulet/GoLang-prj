package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mangalib/auth-service/internal/middleware"
	"github.com/mangalib/auth-service/internal/models"
	"github.com/mangalib/auth-service/internal/services"
	"github.com/mangalib/auth-service/internal/utils"
	"github.com/mangalib/auth-service/internal/validators"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.users.GetAllUsers()
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.OK(c, users)
}
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		Bio      string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.users.GetProfile(uint(id))
	if err != nil {
		utils.NotFound(c, "user")
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Role != "" {
		user.Role = models.Role(req.Role)
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := h.users.Update(user); err != nil {
		utils.InternalError(c)
		return
	}

	utils.OK(c, user)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req validators.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	result, err := h.auth.Register(req)
	if err != nil {
		c.JSON(http.StatusConflict, utils.Response{Success: false, Message: err.Error()})
		return
	}
	utils.Created(c, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req validators.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	result, err := h.auth.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Response{Success: false, Message: err.Error()})
		return
	}
	utils.OK(c, result)
}

func (h *AuthHandler) Validate(c *gin.Context) {
	var body struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	claims, err := h.auth.ValidateToken(body.Token)
	if err != nil {
		utils.Unauthorized(c)
		return
	}
	utils.OK(c, gin.H{
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}

type UserHandler struct {
	users *services.UserService
}

func NewUserHandler(users *services.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	user, err := h.users.GetProfile(userID)
	if err != nil {
		utils.NotFound(c, "user")
		return
	}
	utils.OK(c, user)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}
	user, err := h.users.GetProfile(uint(id))
	if err != nil {
		utils.NotFound(c, "user")
		return
	}
	utils.OK(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req validators.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	user, err := h.users.UpdateBio(userID, req)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.OK(c, user)
}
