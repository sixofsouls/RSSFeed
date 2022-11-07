package api

import (
	"NewsFeed/pkg/db"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// API structure for application
type API struct {
	r  *mux.Router
	db *db.DB
}

// New API constructor
func New(db *db.DB) *API {
	api := API{
		r:  mux.NewRouter(),
		db: db,
	}
	api.endpoints()
	return &api
}

func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {
	// display n-amount of news
	api.r.HandleFunc("/news/{n}", api.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	// web application
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// newsHandler returns requested amount of news
func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	s := mux.Vars(r)["n"]
	amountToShow, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, _ := api.db.ReadPosts(amountToShow)
	json.NewEncoder(w).Encode(posts)
}
