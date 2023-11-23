package database

type Organisation struct {
	ID     int
	Name   string
	City   string
	County string
	Jobs   []int
}

type Job struct {
	ID        int
	Type      string
	Rating    string
	VisaRoute string
}
