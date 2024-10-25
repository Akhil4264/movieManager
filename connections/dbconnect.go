package dbconnection

import (
	"fmt"
	"os"
	"log"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func Connect(){
	err := godotenv.Load(".env")

	if(err != nil){
		panic(err)
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")

	dataSource := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", dbUser, dbName, dbPass, dbHost)

	db,err = sqlx.Connect("postgres",dataSource)

    if err != nil {
        log.Fatalln(err)
    }
  
    defer db.Close()
    if err := db.Ping(); err != nil {
        log.Fatal("Error  connecting to DB : ",err)
    } else {
        log.Println("############Connection to DB established successfully############")
    }
}
