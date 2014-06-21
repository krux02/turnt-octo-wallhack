package helpers

import (
	"code.google.com/p/go.exp/fsnotify"
	"log"
	"time"
)

const Duration = 500 * time.Millisecond

func Listen(file string, observer chan string) {
	log.Println("listen")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	err = watcher.Watch(file)
	if err != nil {
		panic(err)
	}

	for true {
		select {
		case event := <-watcher.Event:
			log.Println(event.Name)

			timer := time.NewTimer(Duration)
			// wait half a second and clear event spamming

			for b := true; b; {
				select {
				case event = <-watcher.Event:
					log.Println("lala", event.Name)
					timer.Reset(Duration)
				case <-timer.C:
					log.Println("timeout", event.Name)
					b = false
				}
			}
			observer <- event.Name
		case err = <-watcher.Error:
			panic(err)
		}
	}
}
