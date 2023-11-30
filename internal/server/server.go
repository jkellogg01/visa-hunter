package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	favicon := http.FileServer(http.Dir("./public/favicon.ico"))
	http.Handle("/favicon.ico", favicon)

	http.HandleFunc("/", handleIndex)
	// http.HandleFunc("/search", handleSearch)

	// paginated queries will take a hash of the query instead of a query in the response body
	http.HandleFunc("/organisations", handleIndexPage)
	// http.HandleFunc("/search/:cursor", handleSearchPage)

	http.HandleFunc("/organisation", handleOrgDetail)

	http.HandleFunc("/seed", handleSeed)

	return http.ListenAndServe(port, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmpls, err := template.ParseGlob("./views/**/*.html")
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(tmpls.DefinedTemplates())

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

	tmpls.ExecuteTemplate(w, "index.html", templateData)
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmpls, err := template.ParseGlob("./views/**/*.html")
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(tmpls.DefinedTemplates())

	pageCursor, err := strconv.Atoi(r.URL.Query().Get("cursor"))
	if err != nil {
		log.Fatal(err)
	}

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
		WHERE
			organisation.id > ?
		GROUP BY
			organisation.id
		LIMIT 48;
	`, pageCursor)
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

	tmpls.ExecuteTemplate(w, "companies", templateData)
}

// func handleSearch(w http.ResponseWriter, r *http.Request) {

// }

// func handleSearchPage(w http.ResponseWriter, r *http.Request) {

// }

func handleOrgDetail(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	defer db.Close()

	var detail struct {
		ID     int64
		Name   string
		City   string
		County string
		Jobs   []struct {
			Type      string
			Rating    string
			VisaRoute string
		}
	}

	orgID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal("atoi error: ", err)
	}
	detail.ID = int64(orgID)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		row := db.QueryRow(`
		SELECT name, city, county FROM organisation WHERE id = ?
		`, detail.ID)

		err := row.Scan(&detail.Name, &detail.City, &detail.County)
		if err != nil {
			log.Fatal()
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rows, err := db.Query(`
		SELECT
			type, 
			rating, 
			visa_route
		FROM
			job
		WHERE
			id in(
			select job_id from organisation_job where organisation_id = ?
			)
		`, detail.ID)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var job struct {
				Type      string
				Rating    string
				VisaRoute string
			}
			if err := rows.Scan(
				&job.Type,
				&job.Rating,
				&job.VisaRoute,
			); err != nil {
				log.Fatal(err)
			}
			detail.Jobs = append(detail.Jobs, job)
		}

		wg.Done()
	}()

	wg.Wait()

	tmpls, err := template.ParseGlob("./views/**/*.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpls.ExecuteTemplate(w, "job-card-full.html", detail)
}

func handleSeed(w http.ResponseWriter, r *http.Request) {
	database.SeedDB()
}
