package main

import (
	"bearguard/action"
	"time"
)

func main() {
	action.Startup()
	for {
		time.Sleep(1 * time.Hour)
	}
}
