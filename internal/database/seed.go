package database

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

type Node[T any] struct {
	Val  *T
	Next *Node[T]
}

type Queue[T any] struct {
	Len  int
	Head *Node[T]
	Tail *Node[T]
}

func newOrgQueue() *Queue[SeedOrganisation] {
	return &Queue[SeedOrganisation]{
		Len:  0,
		Head: nil,
		Tail: nil,
	}
}

func (q *Queue[T]) Push(value *T) {
	node := &Node[T]{Val: value, Next: nil}
	q.Len++

	if q.Len <= 1 {
		q.Head = node
		q.Tail = node
		return
	}

	q.Head.Next = node
	q.Head = node
}

func (q *Queue[T]) Pop() (*T, error) {
	if q.Len < 1 || q.Tail == nil {
		q.Head = nil
		q.Tail = nil
		return nil, fmt.Errorf("no items to pop from queue")
	}
	q.Len--

	result := q.Tail
	q.Tail = q.Tail.Next
	result.Next = nil
	return result.Val, nil
}

func (q *Queue[T]) Peek() (*T, error) {
	if q.Len < 1 {
		return nil, fmt.Errorf("no items to peek on queue")
	}
	return q.Tail.Val, nil
}

// Would like to move some of the actual querying to the query files but
// I have some other data structure things i would want to change first
func SeedDB() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("error connecting to db:", err)
	}
	defer db.Close()

	file, err := os.Open("internal/database/data.csv")
	if err != nil {
		log.Fatal("error opening file:", err)
	}
	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal("error reading file:", err)
	}

	// "Organisation Name","Town/City","County","Type & Rating","Route"
	// 0 -> Organisation.Name
	// 1 -> Organisation.City
	// 2 -> Organisation.County
	// 3 (SPLIT)-> Job.Type, Job.Rating
	// 4 -> Job.VisaRoute

	jobs, err := parseJobs(data)
	if err != nil {
		log.Fatal("failed to parse job data", err)
	}

	orgs, err := parseOrgs(data, jobs)
	if err != nil {
		log.Fatal("failed to parse organisation data", err)
	}

	for _, job := range jobs {
		_, err := db.Exec(
			`insert into job (id, type, rating, visa_route) values (?, ?, ?, ?)`,
			job.ID, job.Type, job.Rating, job.VisaRoute,
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("inserted into job table:", job)
	}

	for orgs.Len > 0 {
		curr, err := orgs.Pop()
		if err != nil {
			log.Fatal("tried to pop nil value from queue end: ", err)
			break
		}

		result, err := db.Exec(
			`insert into organisation (name, city, county) values (?, ?, ?)`,
			curr.Name, curr.City, curr.County,
		)
		if err != nil {
			log.Fatal("error inserting organisation into table", err)
		}

		orgID, err := result.LastInsertId()
		if err != nil {
			log.Fatal("error fething ID for inserted organisation", err)
		}
		curr.ID = orgID
		log.Println("inserted into organisation table:", curr)

		for _, jobID := range curr.Jobs {
			_, err := db.Exec(
				`insert into organisation_job (organisation_id, job_id) values (?, ?)`,
				curr.ID, jobID,
			)
			if err != nil {
				log.Fatal("failed to insert through-table row", err)
			}
			log.Println("inserted into through-table:", curr.ID, jobID)
		}
	}
}

func parseJobs(data [][]string) ([]*SeedJob, error) {
	jobs := []*SeedJob{}

	var nextJobId int64 = 1
	for i, row := range data {
		log.Println("evaluating row: ", row[3:])
		if i == 0 {
			continue
		}
		ratingIdx := strings.Index(row[3], "(")
		jobType := row[3][:ratingIdx]
		jobRating := row[3][ratingIdx:]
		jobVisaRoute := row[4]
		jobRedundant := false
		for _, job := range jobs {
			if job.Type == jobType &&
				job.Rating == jobRating &&
				job.VisaRoute == jobVisaRoute {
				jobRedundant = true
				break
			}
		}
		if jobRedundant {
			// log.Println("job was redundant, did not insert")
			continue
		}

		jobs = append(jobs, &SeedJob{
			ID:        nextJobId,
			Type:      jobType,
			Rating:    jobRating,
			VisaRoute: jobVisaRoute,
		})
		nextJobId++
		log.Println("inserted job: ", row[3:])
	}

	return jobs, nil
}

func parseOrgs(data [][]string, jobs []*SeedJob) (*Queue[SeedOrganisation], error) {
	result := newOrgQueue()

	for i, row := range data {
		if i == 0 {
			continue
		}

		curr := &SeedOrganisation{
			Name:   row[0],
			City:   row[1],
			County: row[2],
			Jobs:   []int64{},
		}

		if result.Len < 1 || result.Head.Val.Name != curr.Name {
			result.Push(curr)
		} else {
			curr = result.Head.Val
		}

		// TODO: make this way more readable (extract function?)
		// This finds the job for the current row in the "jobs" slice
		// and pulls out the ID to append it to the current company's
		// array of job IDs
		ratingIdx := strings.Index(row[3], "(")
		jobType := row[3][:ratingIdx]
		jobRating := row[3][ratingIdx:]
		jobVisaRoute := row[4]
		var jobID int64 = 0
		for _, job := range jobs {
			log.Println(job, row[3:])
			if job.Type == jobType &&
				job.Rating == jobRating &&
				job.VisaRoute == jobVisaRoute {
				jobID = job.ID
				break
			}
		}
		if jobID == 0 {
			return nil, fmt.Errorf("failed to find current job (skill issue)")
		} else if contains(curr.Jobs, jobID) {
			continue
		}
		curr.Jobs = append(curr.Jobs, jobID)
		log.Println("current job in org: ", curr)
	}

	return result, nil
}

func contains(haystack []int64, needle int64) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}
	return false
}
