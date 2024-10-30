package userRepository

import (
	// "fmt"
	"database/sql"
	"errors"
	"reflect"

	// "fmt"
	"time"

	db "github.com/Akhil4264/movieManager/connections"
)

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	CreatedOn int64  `json:"created_on"`
	LastLogin int64  `json:"last_login"`
}

var SignedUsers []User = make([]User, 0)

func AddUser(user User) (*User, error) {
	user.CreatedOn = time.Now().Unix()
	user.LastLogin = user.CreatedOn
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
	row := db.DB.QueryRow("SELECT id,username,email,created_on,last_login FROM User WHERE id =$1", id)
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

func GetAllUsers() ([]User, error) {
	usersList := make([]User, 0)
	row, err := db.DB.Query("SELECT id,username,email FROM Users")
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

func GetUsersByQuery(query string) (*[]User, error) {
	var usersList []User
	rows,err := db.DB.Query("SELECT id,username,email from users where username like %$1%",query)
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
	_,err := db.DB.Exec("UPDATE users SET $1=$2 where id=$3",field,val,id)
	if(err != nil){
		return err
	}
	return nil
}
