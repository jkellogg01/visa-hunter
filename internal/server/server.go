package server

import (
	"html/template"
	"log"
	"net/http"
	"visa-hunter/internal/database"
)

type ResultPage struct {
	Data []struct {
		ID     int64
		Name   string
		City   string
		County string
		Jobs   int
	}
	Cursor int64
}

func Start(port string) error {
	http.HandleFunc("/", handleIndex)
	// http.HandleFunc("/search", handleSearch)

	// paginated queries will take a hash of the query instead of a query in the response body
	http.HandleFunc("/organisations", handleIndexPage)
	// http.HandleFunc("/search/:cursor", handleSearchPage)

	http.HandleFunc("/organisation/:id", handleOrgDetail)

	http.HandleFunc("/seed", handleSeed)

	return http.ListenAndServe(port, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	tmpls, err := template.ParseGlob("./views/**/*.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tmpls.DefinedTemplates())

	rows, err := db.Query(`
	SELECT
		organisation.*,
		(
			SELECT
				COUNT(*)
			FROM
				organisation_job
			WHERE
				organisation_job.organisation_id = organisation.id) AS jobs
		FROM
			organisation
		GROUP BY
			organisation.id
		LIMIT 48;
	`)
	if err != nil {
		log.Fatal(err)
	}

	var results []struct {
		ID     int64
		Name   string
		City   string
		County string
		Jobs   int
	}

	for rows.Next() {
		var row struct {
			ID     int64
			Name   string
			City   string
			County string
			Jobs   int
		}

		err := rows.Scan(&row.ID, &row.Name, &row.City, &row.County, &row.Jobs)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, row)
	}

	templateData := ResultPage{
		Data:   results,
		Cursor: results[len(results)-1].ID,
	}

	// Have to basically rewrite the front end bc of how I changed the data structure so I'll go change that and return to actually ship that data to the client

	tmpls.ExecuteTemplate(w, "index.html", templateData)
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {

}

// func handleSearch(w http.ResponseWriter, r *http.Request) {

// }

// func handleSearchPage(w http.ResponseWriter, r *http.Request) {

// }

func handleOrgDetail(w http.ResponseWriter, r *http.Request) {

}

func handleSeed(w http.ResponseWriter, r *http.Request) {
	database.SeedDB()
}
