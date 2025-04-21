package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	db := connectDB()
	defer db.Close()

	createTableIfNotExists(db)

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Use POST", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10MB max memory
		if err != nil {
			http.Error(w, "ParseMultipartForm error: "+err.Error(), http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file: "+err.Error(), http.StatusBadRequest)
			fmt.Println("FormFile error:", err)
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Read error: "+err.Error(), http.StatusInternalServerError)
			fmt.Println("Read error:", err)
			return
		}

		fmt.Println("Received file:", header.Filename)

		_, err = db.Exec(`INSERT INTO audio_files (filename, file_data) VALUES ($1, $2)`, header.Filename, fileBytes)
		if err != nil {
			http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
			fmt.Println("DB insert error:", err)
			return
		}

		fmt.Fprintln(w, "File uploaded successfully")
	})

	fmt.Println("Running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
