package kvbench

import (
	"github.com/brianvoe/gofakeit/v6"
)

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

type TestObj struct {
	ID        uint64
	UID       string `objectbox:"index"`
	FirstName string
	LastName  string
	Updated   int64
}

func NewTestObj() *TestObj {
	return &TestObj{
		UID:       gofakeit.UUID(),
		FirstName: gofakeit.Name(),
		LastName:  gofakeit.LastName(),
		Updated:   gofakeit.Int64(),
	}
}
