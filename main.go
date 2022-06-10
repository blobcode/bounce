package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"git.mills.io/prologic/bitcask"
	"github.com/bmizerany/pat"
	"github.com/teris-io/shortid"
	"gopkg.in/ini.v1"
)

type RequestBody struct {
	Url string
}

func redirectHandler(db bitcask.Bitcask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// query for id
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

// handle link creation
func newHandler(db bitcask.Bitcask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// init
		id, _ := shortid.Generate()
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
			db.PutWithTTL([]byte(id), []byte(site), 24*time.Hour)
			w.Write(json)
		}
	}
}

func main() {
	// parse config
	cfg, err := ini.Load("bounce.ini")
	if err != nil {
		fmt.Printf("failed to read config: %v", err)
		os.Exit(1)
	}

	port := cfg.Section("config").Key("port").String()
	dbpath := cfg.Section("config").Key("path").String()

	// db setup
	db, _ := bitcask.Open(dbpath)
	defer db.Close()

	// http setup
	m := pat.New()
	m.Get("/r/:id", http.HandlerFunc(redirectHandler(*db)))
	m.Post("/new", http.HandlerFunc(newHandler(*db)))

	http.Handle("/", m)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
