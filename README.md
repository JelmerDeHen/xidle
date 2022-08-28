# xidle

Run functions based on idle state in Xorg.

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
		PollInterval: time.Second,
		IdleLess: func() {
			log.Println("Present")
		},
		IdleLessTimeout: time.Second * 3,
		IdleOver: func() {
			log.Println("User afk")
		},
		IdleOverTimeout: time.Second * 3,
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
- When not idle during last minute spawn the application
- When idle over 10 minutes kill the application

# Examples

- https://github.com/JelmerDeHen/pling
- https://github.com/JelmerDeHen/monarch/blob/main/core/cmd/x11grab_commands.go
- https://github.com/JelmerDeHen/monarch/blob/main/core/cmd/arecord_commands.go

