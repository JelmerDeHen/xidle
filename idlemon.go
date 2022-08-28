package xidle

import (
	"log"
	"time"

	"github.com/JelmerDeHen/scrnsaver"
)

// Define callbacks:
// Poll() is executed each poll, PollInterval defines time to sleep between polls
// IdleLess() is executed when idle time < IdleLessTimeout
// IdleOver() is executed when idle time > IdleOverTimeout
type Idlemon struct {
	Poll            func()
	PollInterval    time.Duration
	IdleLess        func()
	IdleLessTimeout time.Duration
	IdleOver        func()
	IdleOverTimeout time.Duration

	Dbg bool
}

// Run defined callbacks based on how long user is idle
func (im *Idlemon) Run() {
	if im.IdleLess != nil && im.IdleLessTimeout == 0 {
		log.Println("Idlemon.Run(): Set Idlemon.IdleLessTimeout to duration greater than 0")
		return
	}
	if im.IdleOver != nil && im.IdleOverTimeout == 0 {
		log.Println("Idlemon.Run(): Set Idlemon.IdleOverTimeout to duration greater than 0")
		return
	}

	// Prevents execution each poll
	if im.IdleOverTimeout < im.IdleLessTimeout {
		log.Println("Idlemon.Run(): Set Idlemon.IdleOverTimeout to duration greater than idlemon.idleLessTimeout")
		return
	}

	// Configure default poll time to 1 second when this is not configured
	if im.PollInterval == 0 {
		im.PollInterval = time.Second
	}

	for {
		if !scrnsaver.HasXorg() {
			log.Println("Idlemon.Run(): No Xorg session")
			break
		}

		// If user has X session check if user is idle
		info, err := scrnsaver.GetXScreenSaverInfo()
		if err != nil {
			log.Printf("Idlemon.Run(): scrnsaver.GetXScreenSaverInfo(): %s\n", err)
			continue
		}

		if im.Poll != nil {
			im.Poll()
		}

		if im.IdleLess != nil {
			// When idle less than duration
			if info.Idle < im.IdleLessTimeout {
				if im.Dbg {
					log.Printf("User present: info.Idle=%vs < %v\n", info.Idle.Seconds(), im.IdleLessTimeout.Seconds())
				}
				im.IdleLess()
			}
		}

		if im.IdleOver != nil {
			// When idle over duration
			if info.Idle > im.IdleOverTimeout {
				if im.Dbg {
					log.Printf("User idle: info.Idle=%vs > %v\n", info.Idle.Seconds(), im.IdleOverTimeout.Seconds())
				}
				im.IdleOver()
			}
		}

		time.Sleep(im.PollInterval)
	}
}
