package main

import "time"

type Author struct {
	FirstName string
	LastName  string
	Email     string
	Birthdate time.Time
	Added     time.Time
}
