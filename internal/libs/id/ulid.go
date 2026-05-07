package id

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	source  rand.Source = rand.NewSource(time.Now().UnixNano())
	entropy *rand.Rand  = rand.New(source)
)

func New() string {
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return id.String()
}
