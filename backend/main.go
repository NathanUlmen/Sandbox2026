package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "data.sqlite"
	}

	db, err := openSQLite(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("../web"))
	mux.Handle("/", fs)
	mux.HandleFunc("/api/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload struct {
			Markdown string `json:"markdown"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		html, err := renderMarkdownWithPandoc(payload.Markdown)
		if err != nil {
			http.Error(w, "pandoc conversion failed", http.StatusInternalServerError)
			return
		}

		id, err := insertPost(db, payload.Markdown, html, time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			http.Error(w, "db insert failed", http.StatusInternalServerError)
			return
		}

		response := struct {
			ID   int64  `json:"id"`
			HTML string `json:"html"`
		}{
			ID:   id,
			HTML: html,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("listening on :" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}


