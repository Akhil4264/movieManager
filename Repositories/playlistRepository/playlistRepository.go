package playlistRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	// "time"

	"github.com/Akhil4264/movieManager/Repositories/movieRepository"
	db "github.com/Akhil4264/movieManager/connections"
)

var AccessType map[string]bool = map[string]bool{
	"PUBLIC_READ":  true,
	"PUBLIC_WRITE": true,
	"RESTRICTED":   true,
}

// var userAccess map[string]bool = map[string]bool{
// 	"PUBLIC_READ":      true,
// 	"PUBLIC_WRITE":     true,
// 	"RESTRICTED_READ":  true,
// 	"RESTRICTED_WRITE": true,
// }

type Playlist struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	ShareAccess string    `json:"shareaccess"`
	OwnerID     int       `json:"ownerId"`
	MovieCount 	int 		`json:"movieCount,omitempty"`
	CreatedOn   int64      `json:"createdOn"`
	ModifiedAt  int64      `json:"modifiedAt"`
}

type MoviePlaylist struct {
	MovieID    int `json:"movieId"`
	PlaylistID int `json:"playlistId"`
}

type UserPlaylist struct {
	UserId      int  `json:"userId"`
	PlaylistId  int  `json:"playlistId"`
	WriteAccess bool `json:"writeAccess"`
}

type PlayListWithMovies struct{
	Playlist
	movieList []movieRepository.Movie
}

func CreatePlaylist(playlist Playlist) (*Playlist, error) {
	// add userid to owner id before sending here
	tx, err := db.DB.Begin()
	if err != nil {
		fmt.Println("error starting transaction : ", err)
		return nil, err
	}
	query := "INSERT INTO playlist(name,shareaccess,ownerId,createdOn,modifiedAt) values($1,$2,$3,$4,$5) RETURNING id"
	err = tx.QueryRow(query, playlist.Name, playlist.ShareAccess, playlist.OwnerID, playlist.CreatedOn,playlist.ModifiedAt).Scan(&playlist.Id)
	if err != nil {
		fmt.Println("error adding playlist : ", err)
		_ = tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec("INSERT INTO userPlaylist values($1,$2,$3)", playlist.OwnerID, playlist.Id, true)
	if err != nil {
		fmt.Println("error adding userplaylist : ", err)
		_ = tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		fmt.Printf("error commiting the transaction : ", err)
		return nil, err
	}

	return &playlist, nil
}

func UpdatePlayListById(id int, field string, val string) error {
	var v interface{}
	var err error
	switch field {
	case "name":
		if reflect.TypeOf(val).Kind() != reflect.String {
			return errors.New("invalid Value")
		}
		v = val
	case "shareaccess":
		if !AccessType[val] {
			return errors.New("invalid value for access")
		}
		v = val
	case "ownerid":
		v, err = strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid Field")
	}
	query := fmt.Sprintf("UPDATE playlist SET %s=$1,modifiedat=$2 WHERE id=$3",field)
	_, err = db.DB.Exec(query, v,time.Now().Unix(),id)
	if err != nil {
		return err
	}
	return nil
}

func AddMovieToPlaylist(playlistId int, movieId string) error {
	tx,err := db.DB.Begin()
	if(err != nil){
		return err
	}
	err = tx.QueryRow("SELECT imdbId from movie where imdbId=$1",movieId).Scan(&movieId)
	if(err != nil){
		switch err {
		case sql.ErrNoRows:
			_, err = movieRepository.GetMovieByIdRemote(movieId,true)
			if(err != nil){
				fmt.Println("error fetching remote movieId")
				err = tx.Rollback()
				return err
			}
		default:
			_ = tx.Rollback()
			return err
		}
	}
	_, err =  tx.Exec("INSERT INTO movieplaylist values($1,$2)", movieId, playlistId)
	if err != nil {
		fmt.Println("error adding to movieplaylist : ", err)
		_=tx.Rollback()
		return err
	}
	_, err =  tx.Exec("UPDATE playlist set modifiedat=$1 WHERE id=$2", time.Now().Unix(), playlistId)
	if err != nil {
		fmt.Println("error updating modified time of playlist : ", err)
		_=tx.Rollback()
		return err
	}
	if err = tx.Commit(); err!=nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		_=tx.Rollback()
		return err
	}
	return nil
}

func RemoveMovieFromPlaylist(playlistId int, movieId string) error {
	// check if user has write access
	_, err := db.DB.Exec("DELETE from movieplaylist where playlist_id=$1 and movie_id=$2", playlistId, movieId)
	if err != nil {
		fmt.Println("error removing from movieplaylist : ", err)
		return err
	}
	_, err =  db.DB.Exec("UPDATE playlist set modifiedat=$1 WHERE id=$2", time.Now().Unix(), playlistId)
	if err != nil {
		fmt.Println("error updating modified time of playlist : ", err)
		return err
	}
	return nil
}

func SharePlaylistToUser(playlistId int, userId int, writeAccess bool) error {
	_, err := db.DB.Exec("INSERT INTO userplaylist values($1,$2,$3)", userId,playlistId,writeAccess)
	if err != nil {
		fmt.Println("error adding to userplaylist")
		return err
	}
	return nil
}

func UpdateAccessToUser(playlistId int, userId int, writeAccess bool) error {
	// check if user has write access
	_, err := db.DB.Exec("Update userplaylist set writeaccess=$3 where user_id=$1 and playlist_id=$2", userId,playlistId,writeAccess)
	if err != nil {
		fmt.Println("error updating access for user")
		return err
	}
	return nil
}

func CheckAccessForPlaylist(playlistId int, userId int) (int, error) {
	// 3 --> owner
	// 2 --> write
	// 1 --> read
	// 0 --> no access
	var shareAccess string
	var ownerId int
	var user_id int
	var playlist_id int
	var writeaccess bool
	row := db.DB.QueryRow("Select playlist.shareaccess,playlist.ownerid,coalesce(t.user_id,-1) as userid,coalesce(t.playlist_id,-1) as playlistid,coalesce(writeaccess,false) as writeaccess from playlist left join (select * from userplaylist where user_Id=$1 and playlist_Id=$2) as t on t.playlist_Id = playlist.id where playlist.id=$2",userId,playlistId)
	err := row.Scan(&shareAccess,&ownerId,&user_id,&playlist_id,&writeaccess)
	if(err != nil){
		switch err {
			case sql.ErrNoRows:
				return 0,err
			default :
				return -1,err
		}
	}
	if(ownerId == user_id){
		return 3,nil
	}
	if(writeaccess){
		return 2,nil
	}
	if((user_id == -1 || playlist_id == -1)&& shareAccess=="PUBLIC_READ"){
		return 1,nil
	}
	if((user_id == -1 || playlist_id == -1)&& shareAccess=="PUBLIC_WRITE"){
		return 2,nil
	}
	if((user_id == -1 || playlist_id == -1)&& shareAccess=="RESTRICTED"){
		return 0,nil
	}
	if(!writeaccess){
		return 1,nil
	}
	return -1,errors.New("invalid request")
}

func MovieExistsInPlaylist(playlistId int,movieId string)(bool,error){
	err :=  db.DB.QueryRow("select playlist_id from movieplaylist where movie_id=$1 and playlist_id=$2", movieId, playlistId).Scan(&playlistId)
	if(err != nil){
		switch err {
			case sql.ErrNoRows:
				return false,nil
			default:
				return false,err
		}
	}
	return true,nil
}

func PlaylistSharedToUser(playlistId int,userId int)(*UserPlaylist,error){
	var userplaylist UserPlaylist
	err :=  db.DB.QueryRow("select * from userplaylist where user_id=$1 and playlist_id=$2", userId, playlistId).Scan(&userplaylist.UserId,&userplaylist.PlaylistId,&userplaylist.WriteAccess)
	if(err != nil){
		switch err {
			case sql.ErrNoRows:
				return nil,nil
			default:
				return nil,err
		}
	}
	return &userplaylist,nil
}

func RemoveAccessToUser(playlistId int,userId int) error{
	_,err := db.DB.Exec("DELETE from userplaylist where user_id=$1 and playlist_id=$2",userId,playlistId)
	return err
}


func CopyPlayList(playlistId int,userId int)(*Playlist,error){
	row,err := db.DB.Query("select playlist.name,movieplaylist.movieId from playlist left join movieplaylist on movieplaylist.playlist_id=playlist.id where playlist_id=$1",playlistId)
	if(err != nil){
		return nil,err
	}
	var playlistName string
	var movieId string
	err = row.Scan(&playlistName,&movieId)
	if(err != nil){
		return nil,err
	}
	playlist := Playlist{Name: playlistName,OwnerID: userId,CreatedOn: time.Now().Unix(),ModifiedAt: time.Now().Unix(),ShareAccess: "RESTRICTED"}
	newplaylist,err := CreatePlaylist(playlist)
	if(err != nil){
		return nil,err
	}
	count := 0;
	for(row.Next()){
		err = row.Scan(&playlistName,&movieId)
		if(err != nil){
			return nil,err
		}
		err= AddMovieToPlaylist(newplaylist.Id,movieId)
		if(err != nil){
			return nil,err
		}
		count+=1
	}
	return newplaylist,nil
}

func DeletePlayListById(playlistId int) error{
	_,err := db.DB.Exec("Delete from playlist where id=$1",playlistId)
	if(err != nil){
		return err
	}
	return nil
}

func GetPlayListById(playlistId int) (*PlayListWithMovies,error){
	var playlistMovies PlayListWithMovies
	err := db.DB.QueryRow("SELECT COALESCE((SELECT json_agg(jsonb_build_object('playlist_id', mp.playlist_id, 'movie_id', mp.movie_id, 'title', movie.title, 'year', movie.year, 'poster', movie.poster, 'type', movie.type)) FROM movieplaylist AS mp LEFT JOIN movie ON mp.movie_id = movie.Imdbid WHERE mp.playlist_id = $1), '[]') AS movies, playlist.id, playlist.name, playlist.shareAccess, playlist.createdon FROM playlist WHERE playlist.id = $1",playlistId).Scan(&playlistMovies.movieList,playlistMovies.Id,playlistMovies.Name,playlistMovies.ShareAccess,playlistMovies.CreatedOn)
	if(err != nil){
		switch err {
			case sql.ErrNoRows:
				return nil,errors.New("playlist Doesn't exists")
			default:
				return nil,err
		}
	}
	return &playlistMovies,nil
}

func GetPlayLists(userId int)(*[]Playlist,error){
	var playListArray []Playlist
	rows,err := db.DB.Query("SELECT p.* FROM userplaylist JOIN (SELECT playlist.*, COUNT(movieplaylist.*) AS moviecount FROM playlist LEFT JOIN movieplaylist ON playlist.id = movieplaylist.playlist_id GROUP BY playlist.id) AS p ON p.id = userplaylist.playlist_id WHERE user_id = $1",userId)
	if(err != nil){
		return nil,err
	}

	for(rows.Next()){
		var playlist Playlist
		err = rows.Scan(&playlist.Id,&playlist.Name,&playlist.ShareAccess,&playlist.OwnerID,&playlist.CreatedOn,playlist.ModifiedAt,playlist.MovieCount)
		if(err != nil){
			return nil,err
		}
		playListArray = append(playListArray, playlist)
	}
	return &playListArray,nil
}

