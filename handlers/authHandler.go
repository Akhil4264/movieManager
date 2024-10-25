package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// "io/ioutil"
	"io"

	cors "github.com/Akhil4264/movieManager/middlewares"
	"github.com/joho/godotenv"
)



func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login Page")
}
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Signup Page")
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


