package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	apiHostName = "api"
)

func main() {
	http.HandleFunc("/", handler)

	if host, found := os.LookupEnv("API_HOST"); found {
		apiHostName = host
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	res, err := apiUsers(r.Context(), html.EscapeString(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint:errcheck
		return
	}
	if len(res) == 0 {
		res = []apiModel{{Name: "World"}}
	}
	fmt.Fprintf(w, "Hello, %s!\n", res[0].Name)
}

type apiModel struct {
	Name string `json:"name"`
}

func apiUsers(ctx context.Context, params string) (result []apiModel, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users%s", apiHostName, params), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result = []apiModel{}
	e := json.Unmarshal(bytes, &result)
	return result, e
}
