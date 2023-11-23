package database

type Organisation struct {
	ID     int64
	Name   string
	City   string
	County string
	Jobs   []int64
}

type Job struct {
	ID        int64
	Type      string
	Rating    string
	VisaRoute string
}
