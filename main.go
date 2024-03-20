package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Comment struct {
	Author  string    `json:"author"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

type Comments []Comment

const commentsFile = "data/comments.json"

func loadComments() (Comments, error) {
	if _, err := os.Stat(commentsFile); os.IsNotExist(err) {
		return Comments{}, nil
	}

	data, err := os.ReadFile(commentsFile)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	err = json.Unmarshal(data, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func writeComments(comments Comments) error {
	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	err = os.WriteFile(commentsFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func handleComments(writer http.ResponseWriter, request *http.Request) {
	comments, err := loadComments()
	if err != nil {
		log.Print("loadComments:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")

	switch request.Method {
	case http.MethodGet:
		writer.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(writer).Encode(comments)
		if err != nil {
			log.Print("json.NewEncoder.Encode:", err)
			writer.WriteHeader(http.StatusInternalServerError)
		}

	case http.MethodPost:
		var comment Comment
		err := json.NewDecoder(request.Body).Decode(&comment)
		if err != nil {
			log.Print("json.NewDecoder.Decode:", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		comments = append(comments, comment)

		err = writeComments(comments)
		if err != nil {
			log.Print("writeComments:", err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func main() {
	http.HandleFunc("/comments", handleComments)

	http.HandleFunc("/slides", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request,
			"https://docs.google.com/presentation/d/1Y5bDlnVyX_L36x1M0qfs8Y1EbV68_ijIyb-4otr_wkQ/edit?usp=sharing",
			http.StatusFound,
		)
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
