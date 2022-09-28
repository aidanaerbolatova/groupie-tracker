package api

import "net/http"

func (h *Handler) SetEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/", h.HandleMainPage)
	mux.HandleFunc("/artists/", h.HandleArtistPage)
	mux.HandleFunc("/search/", h.HandleSearch)
	mux.HandleFunc("/filters/", h.HandleFilterPage)
	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
}
