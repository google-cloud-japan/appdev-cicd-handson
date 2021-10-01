package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type user struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func main() {
	var dberr error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)
	db, dberr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if dberr != nil {
		log.Fatal("failed to connect database")
	}
	http.HandleFunc("/users/", usersHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	res, err := queryUsers(html.EscapeString(r.URL.Path[7:]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint:errcheck
		return
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint:errcheck
		return
	}
	w.Write(bytes) //nolint:errcheck
}

func queryUsers(code string) (users []user, err error) {
	where := user{}
	if code != "" {
		where.Code = code
	}
	result := db.Where(where).Find(&users)
	return users, result.Error
}
