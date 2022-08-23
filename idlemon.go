package xidle

import (
	"log"
	"time"

	"github.com/JelmerDeHen/scrnsaver"
)

// Define callbacks:
// Poll() is executed each poll; PollT defines time to sleep between polls
// IdleLess() is executed when idle time < IdleLessT
// IdleOver() is executed when idle time > IdleOverT
type Idlemon struct {
	Poll      func()
	PollT     time.Duration
	IdleLess  func()
	IdleLessT time.Duration
	IdleOver  func()
	IdleOverT time.Duration

	Dbg bool
}

// Run defined callbacks based on how long user is idle
func (im *Idlemon) Run() {
	if im.Poll == nil || im.IdleLess == nil || im.IdleOver == nil {
		log.Println("Define callbacks")
		return
	}

	for {
		// Get idle time
		info, err := scrnsaver.GetXScreenSaverInfo()
		if err != nil {
			panic(err)
		}

		im.Poll()

		// When idle less than duration
		if info.Idle < im.IdleLessT {
			if im.Dbg {
				log.Printf("User present: info.Idle=%vs < %v\n", info.Idle.Seconds(), im.IdleLessT.Seconds())
			}
			im.IdleLess()
		}

		// When idle over duration
		if info.Idle > im.IdleOverT {
			if im.Dbg {
				log.Printf("User idle: info.Idle=%vs > %v\n", info.Idle.Seconds(), im.IdleOverT.Seconds())
			}
			im.IdleOver()
		}

		time.Sleep(im.PollT)
	}
}
