package webui

import (
	"embed"
	"net/http"
)

//go:embed web/*
var uiFS embed.FS

// Index возвращает /web/index.html
func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, uiFS, "web/index.html")
	}
}

func Static() http.Handler {
	return http.StripPrefix("/static/", http.FileServerFS(uiFS))
}
