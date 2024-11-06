package movieRepository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "github.com/Akhil4264/movieManager/connections"
	"github.com/joho/godotenv"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Movie struct {
	Title      string `json:"title,omitempty"`
	Year       string `json:"year,omitempty"`
	Rated      string `json:"rated,omitempty"`
	Released   string `json:"released,omitempty"`
	Runtime    string `json:"runtime,omitempty"`
	Genre      string `json:"genre,omitempty"`
	Director   string `json:"director,omitempty"`
	Writer     string `json:"writer,omitempty"`
	Actors     string `json:"actors,omitempty"`
	Plot       string `json:"plot,omitempty"`
	Language   string `json:"language,omitempty"`
	Country    string `json:"country,omitempty"`
	Awards     string `json:"awards,omitempty"`
	Poster     string `json:"poster,omitempty"`
	Metascore  string `json:"metascore,omitempty"`
	ImdbRating string `json:"imdbRating,omitempty"`
	ImdbVotes  string `json:"imdbVotes,omitempty"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"type,omitempty"`
	DVD        string `json:"dvd,omitempty"`
	BoxOffice  string `json:"boxOffice,omitempty"`
	Production string `json:"production,omitempty"`
	Website    string `json:"website,omitempty"`
	Response   string `json:"response,omitempty"`
}

type Rating struct {
	ID      int    `json:"id"`
	MovieID string `json:"movieId"`
	Source  string `json:"source"`
	Value   string `json:"value"`
}

type MovieRatings struct {
	Movie
	Ratings []Rating
}

type MovieRes struct {
	Search       []Movie `json:"search"`
	TotalResults string  `json:"totalResults"`
	Response     string  `json:"response"`
	Pages        int     `json:"pages"`
}

// Intially we'll have only
// imdbId,title,year,type, poster

// give option to only search movies by search string ..that too.. min 3 characters needed, if characters are entered only then we can show them filters

func AddMovie(movieRatings MovieRatings) (*MovieRatings, error) {
	v := reflect.ValueOf(&(movieRatings.Movie)).Elem()
	t := v.Type()
	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		columns = append(columns, field.Name)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		values = append(values, v.Field(i).Interface())
	}

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		fmt.Println("Error creating a transaction: ", err)
		return nil, err
	}

	query := fmt.Sprintf("INSERT INTO movie (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err = tx.Exec(query, values...)
	if err != nil {
		log.Printf("Error inserting movie: %v", err)
		_ = tx.Rollback()
		return nil, err
	}

	// Insert ratings
	ratingQuery := `INSERT INTO ratings (movieId, source, value) VALUES ($1, $2, $3) RETURNING ID`
	for ind, rating := range movieRatings.Ratings {
		var ratingId int
		err = tx.QueryRow(ratingQuery, movieRatings.ImdbID, rating.Source, rating.Value).Scan(&ratingId)
		if err != nil {
			log.Printf("Error inserting rating: %v", err)
			_ = tx.Rollback()
			return nil, err
		}
		movieRatings.Ratings[ind].ID = ratingId
		movieRatings.Ratings[ind].MovieID = movieRatings.ImdbID
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, err
	}

	return &movieRatings, nil
}

func CheckAndAddMovie(movie MovieRatings, conf bool) (*MovieRatings, error) {
	if conf {
		return AddMovie(movie)
	}
	query := "SELECT row_to_json(movieratings) FROM (SELECT m.*, COALESCE((SELECT json_agg(r) FROM ratings r WHERE r.movieId = $1), '[]') AS ratings FROM movie m WHERE m.imdbID = $1) AS movieratings"
	row := db.DB.QueryRow(query, movie.ImdbID)
	var data string
	err := row.Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Movie with ID %v not found in database. Fetching remotely...", movie.ImdbID)
			return AddMovie(movie)
		}
		log.Printf("Error retrieving movie by ID: %v", err)
		return nil, err
	}
	var movieRatings MovieRatings
	err = json.Unmarshal([]byte(data), &movieRatings)
	if err != nil {
		fmt.Println("error deserailizing json data : ", err)
		return nil, err
	}
	log.Printf("Successfully retrieved movie: %v", movieRatings)
	return &movieRatings, nil
}

func GetMovieByIdRemote(id string, status bool) (*MovieRatings, error) {
	godotenv.Load(".env")
	api_key := os.Getenv("OMDB_API_KEY")
	searchString := fmt.Sprintf("http://www.omdbapi.com/?apikey=%v&i=%v&plot=full", api_key, id)
	res, err := http.Get(searchString)
	if err != nil {
		log.Printf("Error fetching movie from remote: %v", err)
		return nil, err
	}
	var movieRatings MovieRatings
	err = json.NewDecoder(res.Body).Decode(&movieRatings)
	if err != nil {
		log.Printf("Error decoding response for movie ID %v: %v", id, err)
		return nil, err
	}
	res.Body.Close()
	return CheckAndAddMovie(movieRatings, status)
}

func GetMovieById(id string) (*MovieRatings, error) {
	query := "SELECT row_to_json(movieratings) FROM (SELECT m.*, COALESCE((SELECT json_agg(r) FROM ratings r WHERE r.movieId = $1), '[]') AS ratings FROM movie m WHERE m.imdbID = $1) AS movieratings"
	row := db.DB.QueryRow(query, id)
	var data string
	err := row.Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Movie with ID %v not found in database. Fetching remotely...", id)
			return GetMovieByIdRemote(id, true)
		}
		log.Printf("Error retrieving movie by ID: %v", err)
		return nil, err
	}
	var movieRatings MovieRatings
	err = json.Unmarshal([]byte(data), &movieRatings)
	if err != nil {
		fmt.Println("error deserailizing json data : ", err)
		return nil, err
	}
	log.Printf("Successfully retrieved movie with id : %s", id)
	return &movieRatings, nil
}

func GetMoviesByQuery(searchQuery string, yearQuery int, page int) (*MovieRes, error) {
	godotenv.Load(".env")
	api_key := os.Getenv("OMDB_API_KEY")
	searchString := fmt.Sprintf("http://www.omdbapi.com/?apikey=%v&s=%v", api_key, searchQuery)
	if yearQuery > 1900 {
		searchString += fmt.Sprintf("&y=%v", yearQuery)
	}
	if page <= 0 {
		page = 1
	}
	searchString += fmt.Sprintf("&page=%v", page)
	client := &http.Client{}
	reqRemote, err := http.NewRequest("GET", searchString, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(reqRemote)
	if err != nil {
		return nil, err
	}
	var mres MovieRes
	err = json.NewDecoder(res.Body).Decode(&mres)
	if err != nil || mres.Response == "False" {
		return nil, err
	}
	totalres, err := strconv.Atoi(mres.TotalResults)
	if err != nil {
		return nil, err
	}
	mres.Pages = int(math.Ceil(float64(totalres) / 10.0))
	return &mres, nil
}
