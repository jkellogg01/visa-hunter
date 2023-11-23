package database

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type OrgNode struct {
	Val  *Organisation
	Next *OrgNode
}

type OrgQueue struct {
	Len  int
	Head *OrgNode
	Tail *OrgNode
}

func New() *OrgQueue {
	return &OrgQueue{
		Len:  0,
		Head: nil,
		Tail: nil,
	}
}

func (q *OrgQueue) Push(org *Organisation) {
	node := &OrgNode{Val: org, Next: nil}
	q.Len++
	if q.Len <= 1 {
		q.Head = node
		q.Tail = node
		return
	}
	q.Head.Next = node
	q.Head = node
	return
}

func (q *OrgQueue) Pop() (*Organisation, error) {
	if q.Len < 1 {
		return nil, fmt.Errorf("no items to pop from queue")
	}
	q.Len--
	result := q.Tail
	q.Tail = result.Next
	result.Next = nil
	return result.Val, nil
}

func (q *OrgQueue) Peek() (*Organisation, error) {
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

}
