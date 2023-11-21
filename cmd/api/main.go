package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"visa-hunter/internal/database"
)

type ResponsePage struct {
	Data []database.Company
	Page int
	Next int
}

func main() {
	// database.SeedDB()

	http.HandleFunc("/", handleIndex)

	http.HandleFunc("/index", handleIndexPage)

	log.Print(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	db := database.MustConnectDB()
	defer db.Close()

	pageqp := r.URL.Query().Get("page")
	pagenum, err := strconv.Atoi(pageqp)
	if err != nil {
		log.Println("page number absent or invalid, returning page zero")
		pagenum = 0
	}

	tmpl := template.Must(template.ParseFiles("views/index.html"))

	cmpData, err := database.GetFiltered(db, map[string]string{}, 48, pagenum*48)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(w, ResponsePage{Data: cmpData, Page: pagenum, Next: pagenum + 1})
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	db := database.MustConnectDB()
	defer db.Close()

	pageqp := r.URL.Query().Get("page")
	pagenum, err := strconv.Atoi(pageqp)
	if err != nil {
		log.Println("Invalid page number, returning page zero")
		pagenum = 0
	}

	tmpl := template.Must(template.ParseFiles("views/index.html"))

	cmpData, err := database.GetFiltered(db, map[string]string{}, 48, pagenum*48)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "results", ResponsePage{Data: cmpData, Page: pagenum, Next: pagenum + 1})
}
