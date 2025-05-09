package handler

import (
	"1337b04rd/internal/adapter/postgres"
	rickmorty "1337b04rd/internal/adapter/rickandmorty"
	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

type CatalogPost struct {
	ID            int
	Title         string
	Content       string
	IMGURL        string
	CommentCount  int
	TimeRemaining int
}

type CatalogPageData struct {
	Posts   []CatalogPost
	Session *domain.Session
}

func ShowCatalog(postService *service.PostService, session *domain.Session) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		posts, err := postService.GetActivePosts()
		if err != nil {
			http.Error(w, "Failed to load posts", http.StatusInternalServerError)
			return
		}

		var catalogPosts []CatalogPost
		for _, post := range posts {
			imgURL := ""
			if len(post.IMGsURLs) > 0 {
				imgURL = post.IMGsURLs[0]
			}

			expireTime := post.LastCommented.Add(15 * time.Minute)
			timeRemaining := int(time.Until(expireTime).Minutes())
			if timeRemaining < 0 {
				timeRemaining = 0
			}

			catalogPosts = append(catalogPosts, CatalogPost{
				ID:            post.ID,
				Title:         post.Title,
				Content:       post.Content,
				IMGURL:        imgURL,
				CommentCount:  0,
				TimeRemaining: timeRemaining,
			})
		}

		data := CatalogPageData{
			Posts:   catalogPosts,
			Session: session,
		}

		tmpl, err := template.ParseFiles("templates/catalog.html")
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render catalog: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func EnsureSession(w http.ResponseWriter, r *http.Request, sessionRepo *postgres.SessionRepo, rmClient *rickmorty.Client) (*domain.Session, error) {
	sessionID := GetSessionID(r)

	if sessionID != "" {
		existingSession, err := sessionRepo.GetSessionByID(sessionID)
		if err == nil {
			if existingSession.ExpiresAt.After(time.Now()) {

				return existingSession, nil
			}

		}

	}

	newSessionID := GenerateSessionID()
	characterID := rand.Intn(826) + 1

	character, err := rmClient.FetchCharacterByID(characterID)
	if err != nil {
		return nil, err
	}

	newSession := domain.Session{
		ID:        newSessionID,
		Name:      character.Name,
		AvatarURL: character.Image,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err = sessionRepo.CreateSession(newSession)
	if err != nil {
		return nil, err
	}

	SetSessionCookie(w, newSessionID)

	return &newSession, nil
}

func ShowArchive(postService *service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := postService.GetArchivePosts()
		if err != nil {
			http.Error(w, "Failed to load archive", http.StatusInternalServerError)
			return
		}

		var catalogPosts []CatalogPost
		for _, post := range posts {
			img := ""
			if len(post.IMGsURLs) > 0 {
				img = post.IMGsURLs[0]
			}
			catalogPosts = append(catalogPosts, CatalogPost{
				ID:            post.ID,
				Title:         post.Title,
				Content:       post.Content,
				IMGURL:        img,
				CommentCount:  0,
				TimeRemaining: 0,
			})
		}

		data := CatalogPageData{
			Posts: catalogPosts,
		}

		tmpl, err := template.ParseFiles("templates/archive.html")
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render archive", http.StatusInternalServerError)
		}
	}
}
