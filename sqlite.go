package main

import (
	"fmt"
	"time"
)

func test() {

	c := make(chan func())
	c <- func() {

		ticker := time.NewTicker(time.Second)
		select {
		case <-ticker.C:
			fmt.Printf("hello")
		}

	}

	for f := range c {

		f()
	}

}
