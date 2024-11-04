package userRepository

import (
	// "fmt"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	// "fmt"
	"reflect"

	// "fmt"
	"time"

	db "github.com/Akhil4264/movieManager/connections"
	"github.com/Akhil4264/movieManager/middlewares/authmiddleware"
)

type User struct {
	Id        int    `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email,omitempty"`
	CreatedOn int64  `json:"created_on,omitempty"`
	LastLogin int64  `json:"last_login,omitempty"`
}

var SignedUsers []User = make([]User, 0)

func AddUser(user User) (*User, error) {
	user.CreatedOn = time.Now().Unix()
	user.Username = strings.ToLower(user.Username)
	user.LastLogin = user.CreatedOn
	fmt.Println(user.CreatedOn)
	row := db.DB.QueryRow("INSERT INTO Users (username, password, email, created_on, last_login) VALUES ($1, $2, $3, $4, $5) RETURNING id,username,email,created_on,last_login", user.Username, user.Password, user.Email, user.CreatedOn, user.LastLogin)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedOn, &user.LastLogin)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}
	return &user, err
}

func FindUserById(id int) (*User, error) {
	var user User
	row := db.DB.QueryRow("SELECT id,username,email,created_on,last_login FROM users WHERE id =$1", id)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedOn, &user.LastLogin)
	if(err != nil){
		switch err {
			case sql.ErrNoRows:
				return nil, nil
			default:
				return nil, err
		}
	}
	return &user,nil
}

func GetUserFromRequest(r *http.Request)(*User,error){
	claims,err := authmiddleware.HandleClaims(r)
	if(err != nil || claims == nil){
		return nil,err
	}
	user,err := FindUserById(int((*claims)["userId"].(float64)))
	if(err != nil){
		return nil,err
	}
	if(user == nil){
		return nil,nil
	}
	return user,nil
}

func FindUserByEmail(mailId string) (*User, error) {
	var user User
	row := db.DB.QueryRow("SELECT id,username,password,email,created_on,last_login FROM Users WHERE email =$1", mailId)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.CreatedOn, &user.LastLogin)
	if err == nil {
		return &user, nil
	}
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func FindUserByName(username string) (*User, error) {
	var user User
	row := db.DB.QueryRow("SELECT id,username,email,created_on,last_login FROM Users WHERE username =$1", username)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedOn, &user.LastLogin)
	if err == nil {
		return &user, nil
	}
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func GetAllUsers(user *User) ([]User, error) {
	var query string
	var id int
	if(user != nil){
		query = "SELECT id,username,email FROM Users WHERE id != $1"
		id = user.Id
	}else{
		query = "SELECT id,username,email FROM Users"
	}
	usersList := make([]User, 0)
	row, err := db.DB.Query(query,id)
	if err != nil {
		return usersList, err
	}
	for row.Next() {
		var user User
		err = row.Scan(&user.Id, &user.Username, &user.Email)
		if err != nil {
			return make([]User, 0), err
		}
		usersList = append(usersList, user)
	}
	return usersList, nil
}

func GetUsersByQuery(id int,query string) (*[]User, error) {
	var usersList []User
	q := fmt.Sprintf("SELECT id, username, email FROM users WHERE username LIKE '%%%s%%' AND id != $1", query)
	rows,err := db.DB.Query(q,id)
	if(err != nil){
		return nil,err
	}

	for(rows.Next()){
		var user User
		err  = rows.Scan(&user.Id,&user.Username,&user.Email)
		if(err != nil){
			return nil,err
		}
		usersList = append(usersList, user)
	}
	return &usersList,err
}

func UpdateUserFieldById(id int,field string,val interface{}) error{
	switch field {
		case "username","password","email":
			if(reflect.TypeOf(val).Kind() != reflect.String){
				return errors.New("invalid datatype")
			}
		case "last_login":
			if(reflect.TypeOf(val).Kind() != reflect.Int64){
				return errors.New("invalid datatype")
			}
		default : 
			return errors.New("invlaid Field")
	}
	query := fmt.Sprintf("UPDATE users SET %s=$1 where id=$2",field)
	_,err := db.DB.Exec(query,val,id)
	if(err != nil){
		return err
	}
	return nil
}
