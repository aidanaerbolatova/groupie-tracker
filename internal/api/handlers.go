package api

import (
	"fmt"
	"groupie-tracker/internal/models"
	"groupie-tracker/internal/store"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	templateArtist = "templates/artist.html" // name of the html file for artists' page
	templateIndex  = "templates/index.html"  // name of the html file for the main page
	templateError  = "templates/error.html"  // name of the html file for the error page
)

type Handler struct {
	store *store.Store
}

func NewHandler() *Handler {
	return &Handler{
		store: &store.Store{},
	}
}

// HandleMainPage is a handler for the main page
func (h *Handler) HandleMainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleErrorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if r.URL.Path != "/" {
		h.handleErrorPage(w, http.StatusNotFound, nil)
		return
	}

	if err := h.store.GetAllArtists(); err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
		return
	}

	tmp, err := template.ParseFiles(templateIndex)
	if err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
		return
	}

	h.store.Result = h.store.AllArtists
	if err := tmp.Execute(w, h.store); err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
	}
}

// HandleSearch is a handler for "/search/" path
func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleErrorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}

	input := r.FormValue("search-artist")
	if input == "" {
		h.handleErrorPage(w, http.StatusBadRequest, nil)
		return
	}
	h.store.Result = h.store.GetSearchResult(input)

	tmp, err := template.ParseFiles(templateIndex)
	if err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
		return
	}

	if err := tmp.Execute(w, h.store); err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
	}
}

// HandleFilterPage is handler for "/filters/" page
func (h *Handler) HandleFilterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleErrorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}
	var err error
	fromCreation := r.FormValue("fromCreation")
	toCreation := r.FormValue("toCreation")
	fromAlbum := r.FormValue("first-album_from")
	toAlbum := r.FormValue("first-album_to")
	location := r.FormValue("searchLocation")
	err = r.ParseForm()
	if err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
	}
	members := r.Form["memberNumber"]

	filter := models.FilterParams{
		CreationDate:    [2]string{fromCreation, toCreation},
		FirstAlbumDate:  [2]string{fromAlbum, toAlbum},
		MembersCheckbox: members,
		Location:        location,
	}
	h.store.Result, err = h.store.GetFilterResult(filter)
	if err != nil {
		h.handleErrorPage(w, http.StatusBadRequest, err)
		return
	}

	tmp, err := template.ParseFiles(templateIndex)
	if err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
		return
	}

	if err := tmp.Execute(w, h.store); err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
	}
}

// HandleArtistPage is a handler for "/artists/" path
func (h *Handler) HandleArtistPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleErrorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if err := h.store.GetAllArtists(); err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Println(r.URL.Path)
	id, err1 := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/artists/"))
	artist, err2 := h.store.GetArtistByID(id)
	if err1 != nil || err2 != nil {
		h.handleErrorPage(w, http.StatusNotFound, err1)
		return
	}

	tmp, err := template.ParseFiles(templateArtist)
	if err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
		return
	}

	if err := tmp.Execute(w, artist); err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
	}
}

type errorHTTP struct {
	Status  int
	Message string
}

// handleErrorPage is a handler for the error page
func (h *Handler) handleErrorPage(w http.ResponseWriter, status int, serverError error) {
	if status >= http.StatusInternalServerError {
		log.Printf("something went wrong: %s", serverError)
	}

	errHTTP := errorHTTP{
		Status:  status,
		Message: http.StatusText(status),
	}

	w.WriteHeader(status)
	tmp, err := template.ParseFiles(templateError)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmp.Execute(w, errHTTP); err != nil {
		h.handleErrorPage(w, http.StatusInternalServerError, err)
	}
}
