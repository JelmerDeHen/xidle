
Run command when user is not idle:

```go
package main

import (
  "github.com/JelmerDeHen/xidle"
)

func main() {
	runner := xidle.NewCmdJob("sleep", "1337")
	idlemon := xidle.NewIdlemon(runner)
	idlemon.Run()
}
```

Example with manually defined callbacks and durations:

```
package main

import (
	"log"
	"time"

  "github.com/JelmerDeHen/xidle"
)

func main() {
	idlemon := &xidle.Idlemon{
		Poll: func() {
			log.Println("Polling")
		},
		PollT: time.Second,
		IdleLess: func() {
			log.Println("Present")
		},
		IdleLessT: time.Second * 3,
		IdleOver: func() {
			log.Println("User afk")
		},
		IdleOverT: time.Second * 3,
	}
	idlemon.Run()
}
```
