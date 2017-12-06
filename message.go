package main

import (
	"time"
)

type message struct {
	Author string
	Body   string
	Time   time.Time
}
