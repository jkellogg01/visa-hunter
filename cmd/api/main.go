package main

import (
	"visa-hunter/internal/database"
)

// type ResponsePage struct {
// 	Data []database.Company
// 	Page int
// 	Next int
// }

func main() {
	database.SeedDB()

	// http.HandleFunc("/", handleIndex)

	// http.HandleFunc("/index", handleIndexPage)

	// http.HandleFunc("/search", handleSearchPage)

	// log.Print(http.ListenAndServe(":8080", nil))
}

// func handleIndex(w http.ResponseWriter, r *http.Request) {
// 	db := database.MustConnectDB()
// 	defer db.Close()

// 	tmpl := template.Must(template.ParseGlob("views/index.html"))

// 	cmpData, err := database.GetFiltered(db, map[string]string{}, 48, 0)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("tmpl.DefinedTemplates(): %v\n", tmpl.DefinedTemplates())

// 	tmpl.ExecuteTemplate(w, "main", ResponsePage{Data: cmpData, Page: 0, Next: 1})
// }

// func handleIndexPage(w http.ResponseWriter, r *http.Request) {
// 	db := database.MustConnectDB()
// 	defer db.Close()

// 	pageqp := r.URL.Query().Get("page")
// 	pagenum, err := strconv.Atoi(pageqp)
// 	if err != nil {
// 		log.Println("Invalid page number, returning page zero")
// 		pagenum = 0
// 	}

// 	tmpl := template.Must(template.ParseFiles("views/index.html"))

// 	cmpData, err := database.GetFiltered(db, map[string]string{}, 48, pagenum*48)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	tmpl.ExecuteTemplate(w, "index", ResponsePage{Data: cmpData, Page: pagenum, Next: pagenum + 1})
// }

// func handleSearchPage(w http.ResponseWriter, r *http.Request) {
// 	db := database.MustConnectDB()
// 	defer db.Close()

// 	pageqp := r.URL.Query().Get("page")
// 	pagenum, err := strconv.Atoi(pageqp)
// 	if err != nil {
// 		pagenum = 0
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		log.Fatal(err)
// 	}

// 	var params map[string]string
// 	for k, v := range r.Form {
// 		log.Println(k, v)
// 	}

// 	tmpl := template.Must(template.ParseFiles("views/index.html"))

// 	cmpData, err := database.GetFiltered(db, params, 48, pagenum*48)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	tmpl.ExecuteTemplate(w, "query", ResponsePage{
// 		Data: cmpData,
// 		Page: pagenum,
// 		Next: pagenum + 1,
// 	})
// }
