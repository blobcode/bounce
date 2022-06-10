package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"git.mills.io/prologic/bitcask"
	"github.com/bmizerany/pat"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type RequestBody struct {
	Url string
}

func redirectHandler(db bitcask.Bitcask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get(":id")
		val, _ := db.Get([]byte(id))

		if val != nil && id != "" {
			http.Redirect(w, r, string(val), http.StatusTemporaryRedirect)
		} else {
			// 404
			w.WriteHeader(404)
			w.Write([]byte("404"))
		}
	}
}

func newHandler(db bitcask.Bitcask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// init
		id, _ := gonanoid.New()
		var b RequestBody

		// decode body
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		site := b.Url

		w.Header().Set("Content-Type", "application/json")
		res := make(map[string]string)
		res["id"] = id
		json, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("500"))
		}

		u, err := url.Parse(site)

		if err != nil || u.Scheme == "" || u.Host == "" {
			w.WriteHeader(500)
			w.Write([]byte("500"))
		} else {
			db.Put([]byte(id), []byte(site))
			w.Write(json)
		}
	}
}

func main() {
	// db setup
	db, _ := bitcask.Open("/tmp/db")
	defer db.Close()
	// test value
	db.Put([]byte("test"), []byte("https://google.com"))

	// http setup
	m := pat.New()
	m.Get("/r/:id", http.HandlerFunc(redirectHandler(*db)))
	m.Post("/new", http.HandlerFunc(newHandler(*db)))

	// Register this pat with the default serve mux so that other packages
	// may also be exported. (i.e. /debug/pprof/*)
	http.Handle("/", m)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}