// Package models
package models

import "encoding/gob"

// Reservation holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

func init() {
	gob.Register(Reservation{})
}
