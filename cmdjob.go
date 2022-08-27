package xidle

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/JelmerDeHen/scrnsaver"
)

type CmdJob struct {
	name string
	arg  []string
	cmd  *exec.Cmd

	// When the outfile contains dynamic values such as timestamp it needs to be regenerated between execs
	// A parameter named "${OUTFILE}" will be replaced by the result of OutfileGenerator()
	OutfileGenerator func() string

	retries int
}

// Update dynamic variables in args between runs
// ${OUTFILE} by output of c.OutfileGenerator()
// ${RESOLUTION} by out screen resolution when Xorg is running or empty string
func (c *CmdJob) replaceDynamicArgs(args []string) []string {
	// Prepare replacements
	var resolution string
	if scrnsaver.HasXorg() {
		resolution = scrnsaver.GetResolution()
	}

	var outfile string
	if c.OutfileGenerator != nil {
		outfile = c.OutfileGenerator()
	}

	replacer := strings.NewReplacer(
		"${OUTFILE}", outfile,
		"${RESOLUTION}", resolution,
	)

	// Replace
	newArgs := make([]string, len(args))
	for i, arg := range args {
		newArgs[i] = replacer.Replace(arg)
	}

	return newArgs
}

func (c *CmdJob) Spawn() {
	// Don't do anything if we are still running
	if c.Running() {
		return
	}

	// Replace variables in args
	args := c.replaceDynamicArgs(c.arg)

	log.Printf("CmdJob.Spawn(): %s %v\n", c.name, strings.Join(args[:], " "))

	c.cmd = exec.Command(c.name, args...)

	var stdcombined bytes.Buffer
	c.cmd.Stdout = &stdcombined
	c.cmd.Stderr = &stdcombined

	go c.cmd.Run()

	// Give time to fail
	time.Sleep(time.Second * 1)

	if !c.Running() {
		if c.cmd.ProcessState != nil && c.cmd.ProcessState.ExitCode() != 0 {
			log.Printf("CmdJob.Spawn(): Starting %s resulted in non-zero exit code after 1 second: errno=%d; output=%q\n", c.name, c.cmd.ProcessState.ExitCode(), stdcombined.String())
			return
		}
	}
}

func (c *CmdJob) Kill() {
	if !c.Running() {
		return
	}
	log.Printf("CmdJob.Kill(): %s %v\n", c.name, c.arg)
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
		name: name,
		arg:  arg,
	}
}

func NewIdlemon(runner *CmdJob) *Idlemon {
	// When user was present last minute then spawn the app
	// When user was idle for over 10 mins kill the app
	idlecmd := &Idlemon{
		IdleLessT: time.Minute,
		IdleOverT: time.Minute * 10,

		IdleLess: runner.Spawn,
		IdleOver: runner.Kill,
	}

	return idlecmd
}
