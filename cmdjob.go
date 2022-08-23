package xidle

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Helper for managing apps with Idlemon
type CmdJob struct {
	name string
	arg  []string

	cmd         *exec.Cmd
	stdcombined bytes.Buffer

	// state
	start time.Time
	etime time.Duration

	etimeMax time.Duration

	Dbg bool

	// When the outfile contains dynamic values such as timestamp it needs to be regenerated between execs
	// A parameter named "${OUTFILE}" will be replaced by the result of OutfileGenerator()
	OutfileGenerator func() string

	retries int
}

func (c *CmdJob) Spawn() {
	if c.Running() {
		return
	}

	//if c.Dbg {
	log.Printf("Spawn(): %s %v\n", c.name, c.arg)
	//}

	// Kill before exec
	if c.Running() {
		c.Kill()
	}

	// Copy to prevent changing original args
	args := make([]string, len(c.arg))
	copy(args, c.arg)

	for i, v := range args {
		if v == "${OUTFILE}" {
			args[i] = c.OutfileGenerator()
		}
	}

	c.cmd = exec.Command(c.name, args...)

	c.cmd.Stdout = &c.stdcombined
	c.cmd.Stderr = &c.stdcombined

	go c.cmd.Run()

	// Give time to fail
	time.Sleep(time.Second * 1)

	if !c.Running() {
		err := fmt.Errorf("%s terminated in <1s: errno=%d\n%s", c.name, c.cmd.ProcessState.ExitCode(), c.stdcombined.String())
		panic(err)
	}

	c.start = time.Now()

  fmt.Println(c.stdcombined.String())
}

func (c *CmdJob) Poll() {
	// Elapsed time
	c.etime = time.Since(c.start)

	// Kill process to rotate outfile
	if c.etime > c.etimeMax {
		if c.Dbg {
			log.Printf("Process elapsed time limit reached: c.etime=%vs > c.etimeMax%vs\n", c.etime.Seconds(), c.etimeMax.Seconds())
		}
		c.Kill()
	}

	//fmt.Printf("%+v\n", c)
}

func (c *CmdJob) Kill() {
	if !c.Running() {
		return
	}
	//if c.Dbg {
	log.Printf("Kill(): %s %v\n", c.name, c.arg)
	//}
	c.cmd.Process.Kill()
}

func (c *CmdJob) Running() bool {
	if c.cmd == nil {
		return false
	}

	if c.cmd != nil && c.cmd.ProcessState != nil {
		//if c.cmd != nil && c.cmd.ProcessState != nil && c.cmd.ProcessState.Exited() {
		return false
	}

	return true
}

func NewCmdJob(name string, arg ...string) *CmdJob {
	return &CmdJob{
		name:     name,
		arg:      arg,
		start:    time.Now(),
		etimeMax: time.Hour,
		// testing
		//etimeMax: time.Second * 20,

	}
}

func NewIdlemon(runner *CmdJob) *Idlemon {
	idlecmd := &Idlemon{
		PollT:     time.Second,
		IdleLessT: time.Minute,
		IdleOverT: time.Minute * 10,

		Poll:     runner.Poll,
		IdleLess: runner.Spawn,
		IdleOver: runner.Kill,
	}

	// testing
	//idlecmd.IdleLessT = time.Second * 1
	//idlecmd.IdleOverT = time.Second * 3

	return idlecmd
}
