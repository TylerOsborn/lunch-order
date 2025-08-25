package middleware

import (
	"lunchorder/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type AuthMiddleware struct {
	store    *sessions.CookieStore
	userRepo *repository.UserRepository
}

func NewAuthMiddleware(sessionSecret string, userRepo *repository.UserRepository) *AuthMiddleware {
	store := sessions.NewCookieStore([]byte(sessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	return &AuthMiddleware{
		store:    store,
		userRepo: userRepo,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := m.store.Get(c.Request, "auth-session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		userID, ok := session.Values["user_id"].(uint)
		if !ok || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		// Get user from database
		var user repository.User
		if err := m.userRepo.GetDB().First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		userData := user.(*repository.User)
		if userData.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) SetUserSession(c *gin.Context, user *repository.User) error {
	session, err := m.store.Get(c.Request, "auth-session")
	if err != nil {
		return err
	}

	session.Values["user_id"] = user.ID
	session.Values["user_email"] = user.Email
	session.Values["user_role"] = user.Role

	return session.Save(c.Request, c.Writer)
}

func (m *AuthMiddleware) ClearUserSession(c *gin.Context) error {
	session, err := m.store.Get(c.Request, "auth-session")
	if err != nil {
		return err
	}

	session.Values["user_id"] = nil
	session.Values["user_email"] = nil
	session.Values["user_role"] = nil
	session.Options.MaxAge = -1

	return session.Save(c.Request, c.Writer)
}