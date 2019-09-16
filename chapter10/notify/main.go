package main

import (
	"log"

	"gopkg.in/fsnotify.v1"
)

func main() {
	counter := 0
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				counter++

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
			if counter > 3 {
				done <- true
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("waiting...")
	<-done
}