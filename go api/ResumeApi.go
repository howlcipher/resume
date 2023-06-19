package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Profile1  string `json:"profile1"`
	Profile2  string `json:"profile2"`
}

type Skill struct {
	Coding_Skill string `json:"coding_skill"`
	Experience   int64  `json:"experience"`
}

func main() {
	db, err := sql.Open("sqlite3", "resume.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/person", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve personal information
		personRows, err := db.Query("SELECT * FROM Personal_Information")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer personRows.Close()

		var people []Person
		for personRows.Next() {
			var person Person
			err := personRows.Scan(&person.Firstname, &person.Lastname, &person.Phone, &person.Email, &person.Profile1, &person.Profile2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			people = append(people, person)
		}

		// Retrieve skills
		skillRows, err := db.Query("SELECT * FROM Skill")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer skillRows.Close()

		var skills []Skill
		for skillRows.Next() {
			var skill Skill
			err := skillRows.Scan(&skill.Coding_Skill, &skill.Experience)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			skills = append(skills, skill)
		}

		// Combine personal information and skills
		type Response struct {
			People []Person `json:"people"`
			Skills []Skill  `json:"skills"`
		}

		response := Response{
			People: people,
			Skills: skills,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
