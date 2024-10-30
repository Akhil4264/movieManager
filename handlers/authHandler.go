package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"log"
	"os"
	"sync"
	"sync/atomic"

	userRepository "github.com/Akhil4264/movieManager/Repositories/userRepository"
	authmiddleware "github.com/Akhil4264/movieManager/middlewares/authmiddleware"
	cors "github.com/Akhil4264/movieManager/middlewares/corsmiddleware"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)



func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loggedUser userRepository.User
	err := json.NewDecoder(r.Body).Decode(&loggedUser)
	if err != nil {
		log.Printf("Error decoding JSON for login: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	log.Printf("Login attempt for email: %s", loggedUser.Email)

	foundUser, err := userRepository.FindUserByEmail(loggedUser.Email)
	if err != nil {
		log.Printf("Error finding user by email (%s): %v", loggedUser.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if (foundUser == nil) {
		log.Printf("User not found for email: %s", loggedUser.Email)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loggedUser.Password))
	if err != nil {
		log.Printf("Password mismatch for email: %s - %v", loggedUser.Email, err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	token, err := authmiddleware.GenToken(foundUser.Id)
	if err != nil {
		log.Printf("Error generating token for user ID %d: %v", foundUser.Id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	log.Printf("User logged in successfully: %s", loggedUser.Email)
	json.NewEncoder(w).Encode(map[string]string{
		"msg": "You have logged in successfully .... ",
	})
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req userRepository.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userFound atomic.Bool
	var goerr atomic.Bool
	var wg sync.WaitGroup

	wg.Add(1)
	go func(userFound *atomic.Bool) {
		defer wg.Done()
		u, err := userRepository.FindUserByEmail(req.Email)
		if (err != nil) {
			goerr.Store(true)
			return
		}
		if (u != nil) {
			userFound.Store(true)
		}
	}(&userFound)

	wg.Add(1)
	go func(userFound *atomic.Bool) {
		defer wg.Done()
		u, err := userRepository.FindUserByName(req.Username)
		if err != nil {
			goerr.Store(true)
			return
		}
		if (u != nil ){
			userFound.Store(true)
		}
	}(&userFound)

	wg.Wait()

	if goerr.Load() {
		fmt.Println("error accessing db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userFound.Load() {
		fmt.Println("User already exists with the provided email or username.")
		w.WriteHeader(http.StatusConflict)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password),10)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Password = string(hashedPass)

	res, err := userRepository.AddUser(req)
	if err != nil {
		fmt.Println("Error adding user to the database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := authmiddleware.GenToken(res.Id)
	if err != nil {
		fmt.Println("Error generating token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"msg": "You have registered successfully...",
	})
}


func GithubCallback(w http.ResponseWriter,r *http.Request){
	var err error = godotenv.Load(".env")
	if(err != nil){
		panic(err)
	}
	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	cors.EnableCors(&w,r)
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	code := data["code"]
	reqMap:= map[string]any{
			"client_id" : client_id,
			"client_secret" : client_secret,
			"code" : code,
			"redirect_uri" : "http://localhost:3000",
	}
	reqBody,err := json.Marshal(reqMap)
	if(err != nil){
		panic(err)
	}
	reqReader := bytes.NewBuffer(reqBody)
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", reqReader)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resptr, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, err := io.ReadAll(resptr.Body)
	if err != nil {
		panic(err)
	}
	var oauthResponse map[string]string
	err = json.Unmarshal(body, &oauthResponse)
	if err != nil {
		panic(err)
	}
	for k,v := range oauthResponse {
		fmt.Println(k," ",v)
	}
}


