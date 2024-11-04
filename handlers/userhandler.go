package handlers

import (
	"encoding/json"
	// "errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Akhil4264/movieManager/Repositories/userRepository"
	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	usersList, err := userRepository.GetAllUsers(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(usersList)
}

func GetUsersByQuery(w http.ResponseWriter, r *http.Request) {

	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	query := vars["Query"]
	query = strings.ToLower(query)
	usersList, err := userRepository.GetUsersByQuery(user.Id, query)
	if err != nil {
		fmt.Println("error accesing users by query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(usersList)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	_, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	userid := vars["userId"]
	id, err := strconv.Atoi(userid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	foundUser, err := userRepository.FindUserById(id)
	if err != nil {
		fmt.Println("error accesing users by query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(foundUser)
}

// func CreateUser(w http.ResponseWriter , r *http.Request){
// 	fmt.Fprintf(w,"Creating an user...")
// }

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	user, err := userRepository.GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	userid := vars["userId"]
	id, err := strconv.Atoi(userid)
	if(user.Id != id){
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var field string
	var val interface{}

	var data map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("error decoding the req body : ",err)
		return
	}

	field, ok := data["field"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("field is invlaid")
		return
	}

	if v, ok := data["val"]; ok {
		switch v.(type) {
		case string:
			val = v.(string)
		case float64: 
			if v, ok := v.(float64); ok && float64(int64(v)) == v {
				val = int64(v)
			} else {
				val = v
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("val must be a string, int, or int64")
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("val is required")
		return
	}
	err = userRepository.UpdateUserFieldById(id,field,val)
	if(err != nil){
		fmt.Println("error updating field : ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// func DeleteUserById(w http.ResponseWriter, r *http.Request) {
// 	user, err := userRepository.GetUserFromRequest(r)
// 	if err != nil {
// 		w.WriteHeader(http.StatusForbidden)
// 		return
// 	}
// 	vars := mux.Vars(r)
// 	userid := vars["userId"]
// 	id, err := strconv.Atoi(userid)
// 	if(user.Id != id){
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
// 	// err = userRepository.
// }
