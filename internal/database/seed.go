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
	orgs := NewOrgQueue()
	for _, row := range data {
		currOrg := &Organisation{
			Name:   row[0],
			City:   row[1],
			County: row[2],
			Jobs:   []*Job{},
		}
		if orgs.Len < 1 || orgs.Head.Val != currOrg {
			orgs.Push(currOrg)
		}
		ratIdx := strings.Index(row[3], "(")
		jobType := row[3][:ratIdx]
		jobRating := row[3][ratIdx:]
		currJob := &Job{
			Type:      jobType,
			Rating:    jobRating,
			VisaRoute: row[4],
		}
		orgs.Head.Val.Jobs = append(orgs.Head.Val.Jobs, currJob)
	}

	for {
		if orgs.Tail == nil {
			break
		}

	}
}
