package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func handleGoogleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	print(oauthStateString)
	println(url)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Handles Google OAuth2 callback, exchanges code for a token, and stores it in session
func handleGoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != oauthStateString {
		return c.String(http.StatusBadRequest, "Invalid OAuth state")
	}

	code := c.QueryParam("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
	}

	// Save the token in the session
	session, _ := store.Get(c.Request(), "session-name")
	session.Values["access_token"] = token.AccessToken
	session.Values["token_type"] = token.TokenType
	session.Values["expiry"] = token.Expiry.Format(time.RFC3339)
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, "/profile")
}

// Displays user profile based on session token
func handleProfile(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session-name")
	accessToken, ok := session.Values["access_token"].(string)
	if !ok || accessToken == "" {
		return c.String(http.StatusUnauthorized, "User not authenticated")
	}

	// Example profile response
	response := map[string]string{
		"access_token": accessToken,
		"expiry":       session.Values["expiry"].(string),
	}

	return c.JSON(http.StatusOK, response)
}

// Logs the user out by clearing the session
func handleLogout(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session-name")
	session.Options.MaxAge = -1 // Delete session
	session.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusSeeOther, "/")
}
