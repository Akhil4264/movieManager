package handlers

import (
	"encoding/json"
	"fmt"

	// "io"
	"net/http"

	movieRepository "github.com/Akhil4264/movieManager/Repositories/movieRepository"
	"github.com/Akhil4264/movieManager/Repositories/userRepository"
	"github.com/gorilla/mux"
)

type searchMovieQuery struct {
	SearchQuery string `json:"search_query"`
	YearQuery   int    `json:"year_query"`
	Page        int    `json:"page"`
}

func GetMovieByQuery(w http.ResponseWriter , r *http.Request){
	_, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var searchq searchMovieQuery
	err = json.NewDecoder(r.Body).Decode(&searchq)
	if(err != nil){
		fmt.Println("error acccessing data from req : ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if (len(searchq.SearchQuery) < 3) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mveres,err := movieRepository.GetMoviesByQuery(searchq.SearchQuery,searchq.YearQuery,searchq.Page)
	if(err != nil){
		fmt.Println("error getting response from api : ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if(mveres == nil){
		err = json.NewEncoder(w).Encode(map[string]any{
			"search" : []string{},
			"totalResults" : 0,
			"response" : "True",
			"pages" : 0,
		})
	}else{
		err = json.NewEncoder(w).Encode(mveres)
	}
	if(err != nil){
		fmt.Println("error secind response : ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetMovieById(w http.ResponseWriter,r *http.Request){
	_, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	movieId := vars["movieId"]
	movie,err := movieRepository.GetMovieById(movieId)
	// movie,err := movieRepository.GetMovieById("tt0499549")
	if(err != nil){
		fmt.Println("error getting movie by id : ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if(movie == nil){
		err = json.NewEncoder(w).Encode(map[string]string{})
	}else{
		err = json.NewEncoder(w).Encode(*movie)
	}
	if(err != nil){
		fmt.Println("error sending response : ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
