package movieRepository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	db "github.com/Akhil4264/movieManager/connections"
	"github.com/joho/godotenv"
)

type Movie struct {
	ID         int       `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Year       string    `json:"year,omitempty"`
	Rated      string    `json:"rated,omitempty"`
	Released   string `json:"released,omitempty"`
	Runtime    string    `json:"runtime,omitempty"`
	Genre      string    `json:"genre,omitempty"`
	Director   string    `json:"director,omitempty"`
	Writer     string    `json:"writer,omitempty"`
	Actors     string    `json:"actors,omitempty"`
	Plot       string    `json:"plot,omitempty"`
	Language   string    `json:"language,omitempty"`
	Country    string    `json:"country,omitempty"`
	Awards     string    `json:"awards,omitempty"`
	Poster     string    `json:"poster,omitempty"`
	Metascore  string    `json:"metascore,omitempty"`
	ImdbRating string    `json:"imdbRating,omitempty"`
	ImdbVotes  string    `json:"imdbVotes,omitempty"`
	ImdbID     string    `json:"imdbID"`
	Type       string    `json:"type,omitempty"`
	DVD        string    `json:"dvd,omitempty"`
	BoxOffice  string    `json:"boxOffice,omitempty"`
	Production string    `json:"production,omitempty"`
	Website    string    `json:"website,omitempty"`
	Response   string     `json:"response,omitempty"`
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

func AddMovie(movie Movie) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	v := reflect.ValueOf(&movie).Elem()
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

	query := fmt.Sprintf("INSERT INTO movie (%s) VALUES (%s)",strings.Join(columns, ", "),strings.Join(placeholders, ", "))
	fmt.Println("insert movie query string : ",query)
	_, err := db.DB.Exec(query, values...)
	if err != nil {
		log.Printf("Error inserting movie: %v", err)
		return err
	}

	log.Printf("Successfully added movie: %v", movie)
	return nil
}

func CheckAndAddMovie(movie Movie, status bool) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if status {
		log.Printf("Status is true; adding movie: %v", movie)
		return AddMovie(movie)
	}

	row := db.DB.QueryRow("SELECT * from movies where ImdbID = $1", movie.ImdbID)
	v := reflect.ValueOf(&movie).Elem()
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Addr().Interface()
	}

	err := row.Scan(values...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Movie not found in database. Adding movie: %v", movie)
			return AddMovie(movie)
		}
		log.Printf("Error checking for existing movie: %v", err)
		return err
	}

	log.Printf("Movie already exists: %v", movie)
	return nil
}

func GetMovieByIdRemote(id string, status bool) (*Movie, error) {
	godotenv.Load(".env")
	api_key := os.Getenv("OMDB_API_KEY")
	searchString := fmt.Sprintf("http://www.omdbapi.com/?apikey=%v&i=%v", api_key, id)
	res, err := http.Get(searchString)
	if err != nil {
		log.Printf("Error fetching movie from remote: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	var movie Movie
	err = json.NewDecoder(res.Body).Decode(&movie)
	if err != nil {
		log.Printf("Error decoding response for movie ID %v: %v", id, err)
		return nil, err
	}

	err = AddMovie(movie)
	if err != nil {
		log.Printf("Error adding movie after fetching: %v", err)
		return &movie, err
	}

	log.Printf("Successfully fetched and added movie: %v", movie)
	return &movie, nil
}

func GetMovieById(id string) (*Movie, error) {
	row := db.DB.QueryRow("SELECT * from movie where ImdbID = $1", id)
	var movie Movie
	v := reflect.ValueOf(&movie).Elem()
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Addr().Interface()
	}

	err := row.Scan(values...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Movie with ID %v not found in database. Fetching remotely...", id)
			return GetMovieByIdRemote(id, true)
		}
		log.Printf("Error retrieving movie by ID: %v", err)
		return nil, err
	}

	log.Printf("Successfully retrieved movie: %v", movie)
	return &movie, nil
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
