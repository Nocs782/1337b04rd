package handler

import (
	"1337b04rd/internal/service"
	"html/template"
	"net/http"
	"time"
)

type CatalogPost struct {
	ID            int
	Title         string
	Content       string // ✅ Added Content field
	IMGURL        string
	CommentCount  int
	TimeRemaining int
}

type CatalogPageData struct {
	Posts []CatalogPost
}

func ShowCatalog(postService *service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get posts from service
		posts, err := postService.GetActivePosts()
		if err != nil {
			http.Error(w, "Failed to load posts", http.StatusInternalServerError)
			return
		}

		// Prepare posts for the catalog template
		var catalogPosts []CatalogPost
		for _, post := range posts {
			imgURL := ""
			if len(post.IMGsURLs) > 0 {
				imgURL = post.IMGsURLs[0]
			}

			// Calculate expiration time
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
				CommentCount:  0, // Placeholder for now
				TimeRemaining: timeRemaining,
			})
		}

		data := CatalogPageData{
			Posts: catalogPosts,
		}

		// Parse and execute template
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
