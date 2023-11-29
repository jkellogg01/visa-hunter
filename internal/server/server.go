package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"visa-hunter/internal/database"
)

type ResultPage struct {
	Data   []database.Organisation
	Cursor int64
}

func Start(port string) error {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/search", handleSearch)

	// paginated queries will take a hash of the query instead of a query in the response body
	http.HandleFunc("/:cursor", handleIndexPage)
	http.HandleFunc("/search/:cursor", handleSearchPage)

	http.HandleFunc("/seed", handleSeed)

	return http.ListenAndServe(port, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`
	SELECT
		organisation.*,
		GROUP_CONCAT(organisation_job.job_id)
	FROM
		organisation
		LEFT JOIN organisation_job ON organisation_job.organisation_id = organisation.id
	GROUP BY
		organisation.id
	LIMIT 48;
	`)
	if err != nil {
		log.Fatal(err)
	}

	var (
		results []database.Organisation
		lastID  int64
	)
	for rows.Next() {
		var (
			row        database.Organisation
			comSepJobs string
		)

		err := rows.Scan(&row.ID, &row.Name, &row.City, &row.County, &comSepJobs)
		if err != nil {
			log.Fatal(err)
		}

		for _, job := range strings.Split(comSepJobs, ",") {
			jobID, err := strconv.Atoi(job)
			if err != nil {
				log.Println(err)
			}
			row.Jobs = append(row.Jobs, int64(jobID))
		}
		lastID = row.ID
		results = append(results, row)
	}

	templateData := ResultPage{
		Data:   results,
		Cursor: lastID,
	}

	// Have to basically rewrite the front end bc of how I changed the data structure so I'll go change that and return to actually ship that data to the client
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {

}

func handleSearch(w http.ResponseWriter, r *http.Request) {

}

func handleSearchPage(w http.ResponseWriter, r *http.Request) {

}

func handleSeed(w http.ResponseWriter, r *http.Request) {
	database.SeedDB()
}
