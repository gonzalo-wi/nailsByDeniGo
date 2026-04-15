package handlers

import (
	"net/http"

	authapp "apiGoShei/internal/application/auth"
	"apiGoShei/internal/domain/client"
	"apiGoShei/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	registerUC   *authapp.RegisterUseCase
	loginUC      *authapp.LoginUseCase
	adminLoginUC *authapp.AdminLoginUseCase
}

func NewAuthHandler(registerUC *authapp.RegisterUseCase, loginUC *authapp.LoginUseCase, adminLoginUC *authapp.AdminLoginUseCase) *AuthHandler {
	return &AuthHandler{registerUC: registerUC, loginUC: loginUC, adminLoginUC: adminLoginUC}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.registerUC.Execute(authapp.RegisterInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
	})
	if err != nil {
		if err == client.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar la cuenta"})
		return
	}

	c.JSON(http.StatusCreated, dto.TokenResponse{ClientID: out.ClientID, Token: out.Token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.loginUC.Execute(authapp.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if err == client.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al iniciar sesión"})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		ClientID:  out.ClientID,
		Token:     out.Token,
		FirstName: out.FirstName,
		LastName:  out.LastName,
		Email:     out.Email,
		Phone:     out.Phone,
	})
}

func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.adminLoginUC.Execute(authapp.AdminLoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if err == authapp.ErrAdminInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al iniciar sesión"})
		return
	}

	c.JSON(http.StatusOK, dto.AdminTokenResponse{AdminID: out.AdminID, Token: out.Token})
}
