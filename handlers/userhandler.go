package handlers

import (
	"fmt"
	"net/http"
)

func GetAllUsers(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"gathering list of all users...")
}

func GetUsersByQuery(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"Accessing users by matching string...")
}

func GetUserById(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"Accessing user by id...")
}


func CreateUser(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"Creating an user...")
}

func UpdateUserById(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"updating user by id...")
}

func DeleteUserById(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"deleteing user by id...")
}

