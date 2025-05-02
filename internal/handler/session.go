package handler

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

func GenerateSessionID() string {
	bytes := make([]byte, 32) // 256-bit session ID
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 1-week expiration
	}
	http.SetCookie(w, cookie)
}
func GetSessionID(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	return cookie.Value
}
