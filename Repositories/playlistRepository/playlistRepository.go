package playlistRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	db "github.com/Akhil4264/movieManager/connections"
)

var accessType map[string]bool = map[string]bool{
	"PUBLIC_READ" : true,
	"PUBLIC_WRITE" : true,
	"RESTRICTED" : true,
}

var userAccess map[string]bool = map[string]bool{
	"PUBLIC_READ" : true,
	"PUBLIC_WRITE" : true,
	"RESTRICTED_READ" : true,
	"RESTRICTED_WRITE" : true,
}

type Playlist struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	ShareAccess string    `json:"share_access"`
	OwnerID     int       `json:"ownerId"`
	CreatedOn   time.Time `json:"createdOn"`
	ModifiedAt  time.Time `json:"modifiedAt"`
}

type MoviePlaylist struct {
	MovieID   int `json:"movieId"`
	PlaylistID int `json:"playlistId"`
}

type UserPlaylist struct {
	UserId      int  `json:"userId"`
	PlaylistId  int  `json:"playlistId"`
	WriteAccess bool `json:"writeAccess"`
}

func AddPlaylist(playlist Playlist)(*Playlist,error){
	// add userid to owner id before sending here
	tx,err := db.DB.Begin()
	if(err != nil){
		fmt.Println("error starting transaction : ",err)
		return nil,err
	}
	query := "INSERT INTO playlist values($1,$2,$3,$4,$4) RETURNING id"
	err = tx.QueryRow(query,playlist.Name,playlist.ShareAccess,playlist.OwnerID,time.Now().Unix()).Scan(&playlist.Id)
	if(err != nil){
		fmt.Println("error adding playlist : ",err)
		_ = tx.Rollback()
		return nil,err
	}
	_,err = tx.Exec("INSERT INTO userPlaylist values($1,$2,$3)",playlist.OwnerID,playlist.Id,true)
	if(err != nil){
		fmt.Println("error adding userplaylist : ",err)
		_ = tx.Rollback()
		return nil,err
	}

	if err = tx.Commit(); err != nil{
		fmt.Printf("error commiting the transaction : ",err)
		return nil,err
	}

	return &playlist,nil
}

func UpdatePlayListById(id int,field string,val string) error {
	var v interface{}
	var err error
	switch field {
		case "Name":
			if(reflect.TypeOf(val).Kind() != reflect.String){
				return errors.New("invalid Value")
			}
			v = val
		case "ShareAccess":
			if(!accessType[val]){
				return errors.New("invalid value for access")
			}
			v = val
		case  "OwnerID" :
			v,err = strconv.ParseInt(val,10,0)
			if(err != nil){
				return err
			}
		case  "CreatedOn","ModifiedAt" :
			v,err = strconv.ParseInt(val,10,64)
			if(err != nil){
				return err
			}
		default :
			return errors.New("invalid Field")
	}
	// check if usser has write access
	_,err = db.DB.Exec("UPDATE playlist set $1=$2 where id=$3",field,v,id)
	if(err != nil){
		return err
	}
	return nil
}

func DeletePlayListById(id int) error {
	_,err := db.DB.Exec("DELETE from playlist where id=$1",id)
	if(err != nil){
		fmt.Println("error deleting a playlist",err)
		return err
	}
	return nil
	// wronggggggg
	// check if user has write access
	// also should remove the userplaylist relations and movieplaylist relations
}

func AddMovieToPlaylist(playlistId int,movieId string) error{
	// check if user has write access
	_,err := db.DB.Exec("INSERT INTO movieplaylist values($1,$2)",playlistId,movieId)
	if(err != nil){
		fmt.Println("error adding to movieplaylist : ",err)
		return err
	}
	return nil
}

func RemoveMovieFromPlaylist(playlistId int,movieId string) error{
	// check if user has write access
	_,err := db.DB.Exec("DELETE from movieplaylist where playlistId=$1 and movieid=$2",playlistId,movieId)
	if(err != nil){
		fmt.Println("error removing from movieplaylist : ",err)
		return err
	}
	return nil
}


func SharePlaylistToUser(playlistId int,userId int,writeAccess bool) error{
	// check if user has write access
	_,err :=db.DB.Exec("INSERT INTO userplaylist values($1,$2,$3)",playlistId,userId,writeAccess)
	if(err != nil){
		fmt.Println("error adding to userplaylist")
	}
	return nil
}


func checkAccessForPlaylist(playlistId int,userId int) (int,error){
	// check if it is public
	// if public --> everyone has readaccess
	// if private --> check for writeaccess, readaccess
	// 2 --> write 
	// 1 --> read 
	// 0 --> no access

	// change it to left join ... 

	var playlistAccess string
	err := db.DB.QueryRow("select ShareAccess from playlist where id=$1",playlistId).Scan(&playlistAccess)
	if(err != nil){
		switch err{
			case sql.ErrNoRows:
				return 0,nil
			default:
				return 0,err
		}

	}
	if(!accessType[playlistAccess]){
		return 0,errors.New("invalid access")
	}

	var userAccess bool
	err = db.DB.QueryRow("select writeAccess from userplaylist where id=$1 and userId=$2",playlistId,userId).Scan(&userAccess)
	if(err != nil){
		switch err {
		case sql.ErrNoRows:
			if(playlistAccess == "RESTRICTED"){
				return 0,nil
			}else if(playlistAccess == "PUBLIC_READ"){
				return 1,nil
			}else{
				return 2,nil
			}
		}
	}

	return 0,nil


	// err = db.DB.QueryRow("select p.shareAccess,t.writeAccess from (select * from playlist where playlistId=$2)as p left join (select * from userplaylist where userId=$1 and playlist=$2)as t on p.id = t.playlist_id ",userId,playlistId).Scan()

}


func GetPlayListById(PlaylistId int,userId int){

}




