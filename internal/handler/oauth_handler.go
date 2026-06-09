package handler

import (
	"BAZ/internal/usecase"
	"BAZ/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	oAuthUsecase *usecase.OAuthUsecase
}

func NewOAuthHandler(uc *usecase.OAuthUsecase) *OAuthHandler {
	return &OAuthHandler{oAuthUsecase: uc}
}

func (h *OAuthHandler) Login(c *gin.Context) {
	provider := c.Param("provider")
	state := utils.GenerateState()
	url, err := h.oAuthUsecase.GetAuthURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, url)
}

func (h *OAuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "missing code"})
		return
	}
	log.Printf("OAuth callback received: provider=%s, code=%s", provider, code)

	resp, err := h.oAuthUsecase.HandleCallback(c.Request.Context(), provider, code)
	if err != nil {
		log.Printf("OAuth callback error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
