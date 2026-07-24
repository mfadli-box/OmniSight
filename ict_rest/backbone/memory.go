package backbone

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type useMemory struct {
	db *sql.DB
}

type SessionMemory struct {
	UserId    string    `json:"user_id"`
	CompanyId string    `json:"company_id"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
	IsAdmin   bool      `json:"is_admin"`
	IsHris    bool      `json:"is_hris"`
}

func USLook(c *gin.Context) string {
	return c.GetString("userId")
}

func USLoad() gin.HandlerFunc {
	var db = PgSQL
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":       "UNAUTHORIZED",
				"error":      "Token required in Authorization header",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		token := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":       "UNAUTHORIZED",
				"error":      "Token is required",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		query := `
			SELECT s.user_id, u.company_id, s.expires_at, u.is_active, u.is_admin, u.is_hris
			FROM   "dat_user_session" s
			JOIN   "dat_user" u ON s.user_id = u.id
			WHERE  s.token = $1
		`
		var session SessionMemory
		ctx := context.Background()
		err := db.QueryRowContext(ctx, query, token).Scan(
			&session.UserId,
			&session.CompanyId,
			&session.ExpiresAt,
			&session.IsActive,
			&session.IsAdmin,
			&session.IsHris,
		)
		if err != nil || session.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":       "UNAUTHORIZED",
				"error":      "Session has expired or is invalid",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		if !session.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":       "FORBIDDEN",
				"error":      "User account is deactivated",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		activeCompanyId := session.CompanyId
		if headerCID := c.GetHeader("X-Company-ID"); headerCID != "" && headerCID != session.CompanyId {
			var exists bool
			err = db.QueryRowContext(ctx,
				`SELECT EXISTS(SELECT 1 FROM "dat_user_company" WHERE user_id = $1 AND company_id = $2)`,
				session.UserId, headerCID,
			).Scan(&exists)
			if err != nil || !exists {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":       "FORBIDDEN",
					"error":      "You do not have access to the selected company",
					"request_id": c.GetString("request_id"),
				})
				return
			}
			activeCompanyId = headerCID
		}

		if activeCompanyId == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":       "FORBIDDEN",
				"error":      "No company selected. Please select a company first.",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		c.Set("userId", session.UserId)
		c.Set("companyId", activeCompanyId)
		c.Set("isAdmin", session.IsAdmin)
		c.Set("isHris", session.IsHris)
		c.Next()
	}
}

func USAuth() gin.HandlerFunc {
	var db = PgSQL
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":       "UNAUTHORIZED",
				"error":      "Token required in Authorization header",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		token := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":       "UNAUTHORIZED",
				"error":      "Token is required",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		query := `
			SELECT s.user_id, s.expires_at, u.is_active, u.is_admin, u.is_hris
			FROM   "dat_user_session" s
			JOIN   "dat_user" u ON s.user_id = u.id
			WHERE  s.token = $1
		`
		var (
			userID    string
			expiresAt time.Time
			isActive  bool
			isAdmin   bool
			isHris    bool
		)
		ctx := context.Background()
		err := db.QueryRowContext(ctx, query, token).Scan(
			&userID, &expiresAt, &isActive, &isAdmin, &isHris,
		)
		if err != nil || expiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":       "UNAUTHORIZED",
				"error":      "Session has expired or is invalid",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		if !isActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":       "FORBIDDEN",
				"error":      "User account is deactivated",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		c.Set("userId", userID)
		c.Set("isAdmin", isAdmin)
		c.Set("isHris", isHris)
		c.Next()
	}
}

func USLock() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":       "FORBIDDEN",
				"error":      "Admin access only",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		c.Next()
	}
}

func USLogs(moduleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() >= http.StatusBadRequest {
			return
		}
		userID := c.GetString("userId")
		if userID == "" {
			return
		}
		companyID := c.GetString("companyId")
		action := strings.TrimSpace(c.Request.Method)
		if action == "" {
			action = "ACTION"
		}
		module := strings.TrimSpace(moduleCode)
		path := strings.TrimSpace(c.Request.URL.Path)
		ipAddress := strings.TrimSpace(c.ClientIP())
		userAgent := strings.TrimSpace(c.GetHeader("User-Agent"))
		auditID := uuid.New().String()
		queryP := `
			INSERT INTO "dat_user_action" (
				id, user_id, company_id, module_code, action, path, ip_address, user_agent, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		`
		_, err := PgSQL.ExecContext(context.Background(), queryP,
			auditID,
			userID,
			companyID,
			module,
			action,
			path,
			ipAddress,
			userAgent)
		if err == nil {
			return
		}
		queryF := `
			INSERT INTO "dat_user_action" (
				id, user_id, company_id, module_code, action, path, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, NOW())
		`
		_, err = PgSQL.ExecContext(context.Background(), queryF,
			auditID,
			userID,
			companyID,
			module,
			action,
			path)
		if err == nil {
			return
		}
		queryM := `
			INSERT INTO "dat_user_action" (
				user_id, action, path
			) VALUES ($1, $2, $3)
		`
		if _, err = PgSQL.ExecContext(context.Background(), queryM,
			userID,
			action,
			path); err != nil {
			Log.Warn().Err(err).Str("module", moduleCode).Msg("failed to write audit log")
		}
	}
}

func USRole(moduleCode string, requiredLevel string) gin.HandlerFunc {
	var db = PgSQL
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if exists && isAdmin.(bool) {
			c.Next()
			return
		}

		userID := c.GetString("userId")
		companyID := c.GetString("companyId")
		if userID == "" || companyID == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":       "FORBIDDEN",
				"error":      "User context not loaded",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		query := `
			SELECT up.level::text
			FROM   "dat_user_privilege" up
			JOIN   "dat_module" m ON m.id = up.module_id
			JOIN   "dat_user_company" uc ON uc.id = up.user_company_id
			WHERE  uc.user_id = $1
			AND    uc.company_id = $2
			AND    m.code = $3
			LIMIT  1
		`
		var level string
		err := db.QueryRowContext(context.Background(), query, userID, companyID, moduleCode).Scan(&level)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":       "FORBIDDEN",
				"error":      "Access denied — no privilege found",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		level = strings.TrimSpace(level)
		if level == "hide" || level == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":       "FORBIDDEN",
				"error":      "Access denied",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		if requiredLevel == "view" {
			if level != "view" && level != "book" && level != "post" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":       "FORBIDDEN",
					"error":      "Insufficient privilege — view required",
					"request_id": c.GetString("request_id"),
				})
				return
			}
		}

		if requiredLevel == "book" {
			if level != "book" && level != "post" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":       "FORBIDDEN",
					"error":      "Insufficient privilege — book required",
					"request_id": c.GetString("request_id"),
				})
				return
			}
		}

		if requiredLevel == "post" {
			if level != "post" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":       "FORBIDDEN",
					"error":      "Insufficient privilege — post required",
					"request_id": c.GetString("request_id"),
				})
				return
			}
		}

		c.Next()
	}
}
