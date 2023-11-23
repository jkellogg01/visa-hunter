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

func newOrgQueue() *Queue[Organisation] {
	return &Queue[Organisation]{
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
	if q.Len < 1 {
		return nil, fmt.Errorf("no items to peek on queue")
	}

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

func SeedDB() {
	db := MustConnectDB()
	defer db.Close()

	file, err := os.Open("internal/database/data.csv")
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal("error reading file: ", err)
	}

	// "Organisation Name","Town/City","County","Type & Rating","Route"
	// 0 -> Organisation.Name
	// 1 -> Organisation.City
	// 2 -> Organisation.County
	// 3 (SPLIT)-> Job.Type, Job.Rating
	// 4 -> Job.VisaRoute

	jobs, err := parseJobs(data)
	if err != nil {
		log.Fatal("failed to parse job data")
	}

	orgs, err := parseOrgs(data, jobs)
	if err != nil {
		log.Fatal("failed to parse organisation data")
	}

	for _, job := range jobs {
		_, err := db.Exec(
			`insert into job (id, type, rating, visa_route) values (?, ?, ?, ?)`,
			job.ID, job.Type, job.Rating, job.VisaRoute,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	for orgs.Len > 0 {
		curr, err := orgs.Pop()
		if err != nil {
			log.Fatal("error fetching next organisation")
		}

		result, err := db.Exec(
			`insert into organisation (name, city, county) values (?, ?, ?)`,
			curr.Name, curr.City, curr.County,
		)
		if err != nil {
			log.Fatal("error inserting organisation into table")
		}

		orgID, err := result.LastInsertId()
		if err != nil {
			log.Fatal("error fething ID for inserted organisation")
		}
		curr.ID = orgID

		for _, jobID := range curr.Jobs {
			db.Exec(
				`insert into organisation_job (organisation_id, job_id) values (?, ?)`,
				curr.ID, jobID,
			)
		}
	}
}

func parseJobs(data [][]string) ([]*Job, error) {
	jobs := []*Job{}

	var nextJobId int64 = 0
	for _, row := range data {
		ratingIdx := strings.Index(row[3], "(")
		jobType := row[3][:ratingIdx]
		jobRating := row[3][ratingIdx+1 : len(row[3])-1]
		jobVisaRoute := row[4]
		jobRedundant := false
		for _, job := range jobs {
			if job.Type == jobType ||
				job.Rating == jobRating ||
				job.VisaRoute == jobVisaRoute {
				jobRedundant = true
				break
			}
		}
		if jobRedundant {
			continue
		}
		nextJobId++
		jobs = append(jobs, &Job{
			ID:        nextJobId,
			Type:      jobType,
			Rating:    jobRating,
			VisaRoute: jobVisaRoute,
		})
	}

	return jobs, nil
}

func parseOrgs(data [][]string, jobs []*Job) (*Queue[Organisation], error) {
	result := newOrgQueue()

	for _, row := range data {
		curr := &Organisation{
			Name:   row[0],
			City:   row[1],
			County: row[2],
			Jobs:   []int64{},
		}

		if result.Len < 1 || result.Head.Val == curr {
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
		jobRating := row[3][ratingIdx+1 : len(row[3])-1]
		jobVisaRoute := row[4]
		var jobID int64 = 0
		for _, job := range jobs {
			if job.Type == jobType &&
				job.Rating == jobRating &&
				job.VisaRoute == jobVisaRoute {
				jobID = job.ID
				break
			}
		}
		if jobID == 0 {
			return nil, fmt.Errorf("failed to find current job (skill issue)")
		}

		curr.Jobs = append(curr.Jobs, jobID)
	}

	return result, nil
}
