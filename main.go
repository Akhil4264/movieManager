package main

import (
	"fmt"
	"log"
	"net/http"

	handler "github.com/Akhil4264/movieManager/handlers"
	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "server is up and running")
}

func main() {

	mux := mux.NewRouter()

	authRoute := mux.PathPrefix("/auth").Subrouter()
	usersRoute := mux.PathPrefix("/users").Subrouter()
	moviesRoute := mux.PathPrefix("/movies").Subrouter()
	playlistRoute := mux.PathPrefix("/playlists").Subrouter()

	mux.HandleFunc("/", homeHandler)
	authRoute.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "base auth route")
	}).Methods("GET")
	authRoute.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	authRoute.HandleFunc("/signup", handler.SignupHandler).Methods("POST")
	authRoute.HandleFunc("/github/callback", handler.GithubCallback).Methods("POST")
	authRoute.HandleFunc("/github/callback/", handler.GithubCallback).Methods("POST")

	usersRoute.HandleFunc("", handler.GetAllUsers).Methods("GET")
	usersRoute.HandleFunc("/", handler.GetAllUsers).Methods("GET")
	usersRoute.HandleFunc("/", handler.CreateUser).Methods("POST")
	usersRoute.HandleFunc("/{userId}", handler.GetUserById).Methods("GET")
	usersRoute.HandleFunc("/{userId}", handler.UpdateUserById).Methods("PUT")
	usersRoute.HandleFunc("/{userId}", handler.DeleteUserById).Methods("DELETE")
	usersRoute.HandleFunc("/query/{Query}", handler.GetUsersByQuery).Methods("GET")

	moviesRoute.HandleFunc("/search", handler.GetMovieByQuery).Methods("POST")

	playlistRoute.HandleFunc("", handler.CreatePlaylist).Methods("POST")
	playlistRoute.HandleFunc("/", handler.CreatePlaylist).Methods("POST")
	playlistRoute.HandleFunc("/{playlistId}", handler.GetPlayListById).Methods("GET")
	playlistRoute.HandleFunc("/{playlistId}", handler.UpdatePlayListById).Methods("PUT")
	playlistRoute.HandleFunc("/{playlistId}", handler.DeletePlayListById).Methods("DELETE")
	playlistRoute.HandleFunc("/{playlistId}/share/{userId}", handler.SharePlaylistToUser).Methods("POST")
	playlistRoute.HandleFunc("/{playlistId}/movie/{movieId}", handler.AddMovieToPlaylist).Methods("POST")
	playlistRoute.HandleFunc("/{playlistId}/movie/{movieId}", handler.RemoveMovieFromPlaylist).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":80", mux))

}
