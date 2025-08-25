package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lunchorder/repository"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	oauthConfig  *oauth2.Config
	adminEmail   string
}

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewAuthService(userRepo *repository.UserRepository, clientID, clientSecret, redirectURL string) *AuthService {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		userRepo:    userRepo,
		oauthConfig: config,
		adminEmail:  "tyler.osborn@impact.com",
	}
}

func (s *AuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *AuthService) HandleCallback(code string) (*repository.User, error) {
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	userInfo, err := s.getUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	// Check if user exists
	user, err := s.userRepo.GetUserByGoogleID(userInfo.ID)
	if err != nil {
		// User doesn't exist, create new user
		role := "standard"
		if userInfo.Email == s.adminEmail {
			role = "admin"
		}

		user = &repository.User{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			GoogleID: userInfo.ID,
			Role:     role,
		}

		err = s.userRepo.CreateUser(user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %v", err)
		}
	}

	return user, nil
}

func (s *AuthService) getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info from Google")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *AuthService) IsAdmin(user *repository.User) bool {
	return user.Role == "admin"
}