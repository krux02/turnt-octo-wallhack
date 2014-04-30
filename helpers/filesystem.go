package helpers

import (
	"code.google.com/p/go.exp/fsnotify"
	"time"
)

func Listen(file string, observer chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	err = watcher.WatchFlags(file, fsnotify.FSN_MODIFY)
	if err != nil {
		panic(err)
	}

	for true {
		select {
		case event := <-watcher.Event:
			// wait half a second and clear event spamming
			t := time.Tick(500 * time.Millisecond)
			for b := true; b; {
				select {
				case event = <-watcher.Event:
				case <-t:
					b = false
				}
			}
			observer <- event.Name
		case err = <-watcher.Error:
			panic(err)
		}
	}
}
