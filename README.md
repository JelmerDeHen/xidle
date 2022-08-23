# Usage

Example with callbacks and durations:

```sh
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

# Run application

Run application when user is present and kill app when afk.

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

Default behavior:
- When application is running for over 1 hour, kill and respawn application (rotate outfile)
- When not idle during last minute spawn the application
- When idle over 10 minutes kill the application

