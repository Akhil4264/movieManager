package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Akhil4264/movieManager/Repositories/playlistRepository"
	"github.com/Akhil4264/movieManager/Repositories/userRepository"
	dbconnection "github.com/Akhil4264/movieManager/connections"
	"github.com/gorilla/mux"
)

func CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var playlist playlistRepository.Playlist
	err = json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		fmt.Println("Error decoding req body : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	playlist.OwnerID = user.Id
	if !playlistRepository.AccessType[playlist.ShareAccess] {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playlist.CreatedOn = time.Now().Unix()
	playlist.ModifiedAt = playlist.CreatedOn
	plist, err := playlistRepository.CreatePlaylist(playlist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if plist == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(plist)
}

func UpdatePlayListById(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	playlistIdStr := vars["playlistId"]
	playlistId, err := strconv.Atoi(playlistIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = dbconnection.DB.QueryRow("SELECT id from playlist WHERE id=$1 AND ownerId=$2", playlistId, user.Id).Scan(&playlistId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	var reqBody map[string]string
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if reqBody["field"] == "" || reqBody["val"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = playlistRepository.UpdatePlayListById(playlistId, reqBody["field"], reqBody["val"])
	if err != nil {
		fmt.Println("error updating playlist : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AddMovieToPlaylist(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("Error finding user from request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Println("No user found in request.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	playlistIdStr := vars["playlistId"]
	movieId := vars["movieId"]
	playlistId, err := strconv.Atoi(playlistIdStr)
	if err != nil {
		fmt.Println("Error converting playlist ID to integer:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = dbconnection.DB.QueryRow("SELECT playlist_Id FROM userplaylist WHERE user_id=$1 AND playlist_id=$2 AND writeaccess=true", user.Id, playlistId).Scan(&playlistId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No playlist found for user with ID:", user.Id, "and playlist ID:", playlistId)
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			fmt.Println("Error querying the database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	exists, err := playlistRepository.MovieExistsInPlaylist(playlistId, movieId)
	if err != nil {
		fmt.Println("Error checking movie in playlist:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "Movie already exists in this playlist",
		})
		return
	}
	err = playlistRepository.AddMovieToPlaylist(playlistId, movieId)
	if err != nil {
		fmt.Println("Error adding movie to playlist:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RemoveMovieFromPlaylist(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("Error finding user from request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Println("No user found in request.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	playlistIdStr := vars["playlistId"]
	movieId := vars["movieId"]
	playlistId, err := strconv.Atoi(playlistIdStr)
	if err != nil {
		fmt.Println("Error converting playlist ID to integer:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	access, err := playlistRepository.CheckAccessForPlaylist(playlistId, user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if access < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = playlistRepository.RemoveMovieFromPlaylist(playlistId, movieId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SharePlaylistToUser(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("Error finding user from request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Println("No user found in request.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var userplaylist playlistRepository.UserPlaylist
	err = json.NewDecoder(r.Body).Decode(&userplaylist)
	if err != nil {
		fmt.Println("Error parsing req body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// checking access for the user with token
	access, err := playlistRepository.CheckAccessForPlaylist(userplaylist.PlaylistId, user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if access < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//checking if the relation exists for the user for whom the request is made
	playlist, err := playlistRepository.PlaylistSharedToUser(userplaylist.PlaylistId, userplaylist.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if playlist != nil {
		if userplaylist.WriteAccess == playlist.WriteAccess {
			err = json.NewEncoder(w).Encode(map[string]string{
				"msg": "Access is already present",
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		err = playlistRepository.UpdateAccessToUser(userplaylist.PlaylistId, userplaylist.UserId, userplaylist.WriteAccess)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	err = playlistRepository.SharePlaylistToUser(userplaylist.PlaylistId, userplaylist.UserId, userplaylist.WriteAccess)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeletePlayListById(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("Error finding user from request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Println("No user found in request.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	playlistIdStr := vars["playlistId"]
	playlistId, err := strconv.Atoi(playlistIdStr)
	if err != nil {
		fmt.Println("Error converting playlist ID to integer:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	access, err := playlistRepository.CheckAccessForPlaylist(playlistId, user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if access < 3 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = playlistRepository.DeletePlayListById(playlistId)
	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CopyPlayListById(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("Error finding user from request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Println("No user found in request.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	playlistIdStr := vars["playlistId"]
	playlistId, err := strconv.Atoi(playlistIdStr)
	if err != nil {
		fmt.Println("Error converting playlist ID to integer:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	access, err := playlistRepository.CheckAccessForPlaylist(playlistId, user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if access < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	playlist,err := playlistRepository.CopyPlayList(playlistId,user.Id)
	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&playlist)
}

func GetPlayListById(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("Error finding user from request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Println("No user found in request.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	playlistIdStr := vars["playlistId"]
	playlistId, err := strconv.Atoi(playlistIdStr)
	if err != nil {
		fmt.Println("Error converting playlist ID to integer:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	access, err := playlistRepository.CheckAccessForPlaylist(playlistId, user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if access < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	PlayListWithMovies,err := playlistRepository.GetPlayListById(playlistId)
	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&PlayListWithMovies)
}

func GetPlayLists(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("Error finding user from request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Println("No user found in request.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playlists,err := playlistRepository.GetPlayLists(user.Id)
	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&playlists)
}

func RemoveAccessToUser(w http.ResponseWriter,r *http.Request){
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("Error finding user from request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Println("No user found in request.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	playlistIdStr := vars["playlistId"]
	userIdStr := vars["userId"]
	playlistId, err := strconv.Atoi(playlistIdStr)
	if err != nil {
		fmt.Println("Error converting playlist ID to integer:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		fmt.Println("Error converting playlist ID to integer:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if(user.Id == userId){
		w.WriteHeader(http.StatusNotAcceptable)
		return 
	}
	access, err := playlistRepository.CheckAccessForPlaylist(playlistId, user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if access < 3 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
}


