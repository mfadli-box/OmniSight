package SP00

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	"ict_rest/mechanic"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NHand(u UseCase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) ListHrisCompany(c *gin.Context) {
	list, err := h.usecase.ListHrisCompany()
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "HRIS companies retrieved",
		"data":    list,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Username and password are required"))
		return
	}

	if req.CompanyID == "" || req.CompanyID == "Non-HRIS" {
		result, err := h.usecase.Login(
			req.Username,
			req.Password,
			c.ClientIP(),
			c.GetHeader("User-Agent"),
		)
		if err != nil {
			mechanic.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"data":    result,
		})
	} else {
		hrisLink, err := h.usecase.GetCompanyHrisLink(req.CompanyID)
		if err != nil || hrisLink == "" {
			mechanic.Error(c, mechanic.ValidationError("Company does not have HRIS integration configured"))
			return
		}

		externalUSR := strings.ReplaceAll(hrisLink, "{$U}", req.Username)
		externalURL := strings.ReplaceAll(externalUSR, "{$P}", req.Password)

		client := &http.Client{Timeout: 10 * time.Second}
		httpResp, err := client.Get(externalURL)
		if err != nil {
			mechanic.Error(c, mechanic.ExternalServiceError("Failed to connect to HRIS server", err))
			return
		}
		defer httpResp.Body.Close()

		respBody, err := io.ReadAll(httpResp.Body)
		if err != nil {
			mechanic.Error(c, mechanic.ExternalServiceError("Failed to read HRIS response", err))
			return
		}

		var soapResp HrisSoapResponse
		if err := xml.Unmarshal(respBody, &soapResp); err != nil {
			mechanic.Error(c, mechanic.ExternalServiceError("Invalid HRIS response format", err))
			return
		}

		respText := strings.TrimSpace(soapResp.Text)
		if !strings.HasPrefix(respText, "1|") {
			mechanic.Error(c, mechanic.Unauthorized("HRIS authentication failed"))
			return
		}

		result, err := h.usecase.LoginHris(
			req.Username,
			req.Password,
			req.CompanyID,
			c.ClientIP(),
			c.GetHeader("User-Agent"),
		)
		if err != nil {
			mechanic.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "HRIS login successful",
			"data":    result,
		})
	}
}

func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := ""
	if len(authHeader) > 7 {
		token = authHeader[7:]
	}
	if token == "" {
		mechanic.Error(c, mechanic.ValidationError("Token is required"))
		return
	}
	if err := h.usecase.Logout(token); err != nil {
		mechanic.Error(c, mechanic.InternalError("Failed to logout", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
