package handlers

import (
	"encoding/json"
	"fmt"
	// "io"
	"net/http"

	movieRepository "github.com/Akhil4264/movieManager/Repositories/movieRepository"
)

type searchMovieQuery struct{
	SearchQuery string
	YearQuery int
	Page int
}

func GetMovieByQuery(w http.ResponseWriter , r *http.Request){
	var searchq searchMovieQuery
	err := json.NewDecoder(r.Body).Decode(&searchq)
	if(err != nil){
		fmt.Println("error acccessing data from req : ",err)
	}
	mveres,err := movieRepository.GetMoviesByQuery(searchq.SearchQuery,searchq.YearQuery,searchq.Page)
	if(err != nil){
		fmt.Println("error getting response from api : ",err)
	}
	err = json.NewEncoder(w).Encode(mveres)
	if(err != nil){
		fmt.Println("error secind response : ",err)
	}
}

func GetMovieById(w http.ResponseWriter,r *http.Request){
	movie,err := movieRepository.GetMovieById("tt0499549")
	if(err != nil){
		fmt.Println("error getting movie by id : ",err)
	}
	if(movie == nil){
		err = json.NewEncoder(w).Encode(map[string]string{})
	}else{
		err = json.NewEncoder(w).Encode(*movie)
	}
	if(err != nil){
		fmt.Println("error sending response : ",err)
	}
	
}
func AddMovie(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"Adding movie to table")
}
