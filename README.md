# Sniffy

Extends fsnotify.v1 to include
- Recursive directory watcher 
- Event filter
- Event Transmitter (demultiplexer).

### Install
`github.com/alihammad-gist/sniffy`

### Example

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alihammad-gist/sniffy"
)

func main() {
	// setting criteria for transmitters
	t1 := sniffy.Transmitter(
		sniffy.OpFilter(sniffy.Create),
		sniffy.TooSoonFilter(time.Second/2),
	)

	t2 := sniffy.Transmitter(
		sniffy.ExtFilter(".jsx"),
		sniffy.TooSoonFilter(time.Second/2),
	)

	w, err := sniffy.NewWatcher(t1, t2)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case e := <-t1.Events:
				fmt.Println("t1:", e.Name, e.Op)
			case e := <-t2.Events:
				fmt.Println("t2:", e.Name, e.Op)
			case err := <-w.Errors:
				log.Println(err)
			}
		}
	}()

	w.AddDir("t")

	done := make(chan bool)
	<-done
}
```

### Testing
So far tested on linux. Run `go test`

