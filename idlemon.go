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
	if im.IdleLess != nil && im.IdleLessT == 0 {
		log.Printf("Set Idlemon.IdleLessT to duration greater than 0\n")
		return
	}
	if im.IdleOver != nil && im.IdleOverT == 0 {
		log.Printf("Set Idlemon.IdleOverT to duration greater than 0\n")
		return
	}

	// Prevents execution each poll
	if im.IdleOverT < im.IdleLessT {
		log.Printf("Set Idlemon.IdleOverT to duration greater than idlemon.idleLessT\n")
		return
	}

	// Configure default poll time to 1 second when this is not configured
	if im.PollT == 0 {
		im.PollT = time.Second
	}

	if im.IdleLess == nil || im.IdleOver == nil {
		log.Println("Define callbacks")
		return
	}

	for {
		if !scrnsaver.HasXorg() {
			log.Printf("Idlemon.Run(): User has no Xorg session\n")
			break
		}

		// If user has X session check if user is idle
		info, err := scrnsaver.GetXScreenSaverInfo()
		if err != nil {
			log.Printf("scrnsaver.GetXScreenSaverInfo(): %s\n", err)
			continue
		}

		if im.Poll != nil {
			im.Poll()
		}

		if im.IdleLess != nil {
			// When idle less than duration
			if info.Idle < im.IdleLessT {
				if im.Dbg {
					log.Printf("User present: info.Idle=%vs < %v\n", info.Idle.Seconds(), im.IdleLessT.Seconds())
				}
				im.IdleLess()
			}
		}

		if im.IdleOver != nil {
			// When idle over duration
			if info.Idle > im.IdleOverT {
				if im.Dbg {
					log.Printf("User idle: info.Idle=%vs > %v\n", info.Idle.Seconds(), im.IdleOverT.Seconds())
				}
				im.IdleOver()
			}
		}

		time.Sleep(im.PollT)
	}
}
