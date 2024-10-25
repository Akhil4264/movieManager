package handlers

import (
	"fmt"
	"net/http"
)

func GetMovieByQuery(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"Accessing movie by reference...")
}
