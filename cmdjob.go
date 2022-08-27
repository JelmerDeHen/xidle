package xidle

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/JelmerDeHen/scrnsaver"
)

type CmdJob struct {
	Name string
	Args []string
	Cmd  *exec.Cmd

	CurrentArgs []string
	Output      *bytes.Buffer

	// Set to os.Interrupt (SIGINT) or os.Kill (SIGKILL)
	KillSignal os.Signal

	// When the outfile contains dynamic values such as timestamp it needs to be regenerated between execs
	// A parameter named "${OUTFILE}" will be replaced by the result of OutfileGenerator()
	OutfileGenerator func() string
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

// Idlemon.IdleLess callback
func (c *CmdJob) Run() {
	// Don't do anything if we are still running
	if c.Running() {
		return
	}

	// Replace variables in args
	c.CurrentArgs = c.replaceDynamicArgs(c.Args)

	log.Printf("CmdJob.Run(): %s %v\n", c.Name, strings.Join(c.CurrentArgs[:], " "))

	c.Cmd = exec.Command(c.Name, c.CurrentArgs...)

	c.Output = new(bytes.Buffer)
	c.Cmd.Stdout = c.Output
	c.Cmd.Stderr = c.Output

	go c.Cmd.Run()

	// Give time to fail
	time.Sleep(time.Second * 1)

	if !c.Running() {
		if c.Cmd.ProcessState != nil && c.Cmd.ProcessState.ExitCode() != 0 {
			log.Printf("CmdJob.Run(): Starting %s resulted in non-zero exit code after 1 second: errno=%d; output=%q\n", c.Name, c.Cmd.ProcessState.ExitCode(), c.Output.String())
			return
		}
	}
}

// Idlemon.IdleOver callback
// Some applications can be terminated gracefully by configuring KillSignal. ffmpeg for example listens to syscall.SIGINT
func (c *CmdJob) Kill() {
	if !c.Running() {
		return
	}

	log.Printf("CmdJob.Kill(): %s %v\n", c.Name, strings.Join(c.CurrentArgs[:], " "))

	if c.KillSignal == nil {
		c.Cmd.Process.Kill()
		return
	}
	c.Cmd.Process.Signal(c.KillSignal)
}

func (c *CmdJob) Running() bool {
	if c.Cmd == nil {
		return false
	}

	if c.Cmd != nil && c.Cmd.ProcessState != nil {
		//if c.Cmd != nil && c.Cmd.ProcessState != nil && c.Cmd.ProcessState.Exited() {
		return false
	}

	return true
}

// Deprecate?
func NewCmdJob(name string, arg ...string) *CmdJob {
	return &CmdJob{
		Name: name,
		Args: arg,
	}
}

func NewIdlemon(c *CmdJob) *Idlemon {
	// When user was present last minute then spawn the app
	// When user was idle for over 10 mins kill the app
	idlecmd := &Idlemon{
		IdleLessT: time.Minute,
		IdleOverT: time.Minute * 10,

		IdleLess: c.Run,
		IdleOver: c.Kill,
	}

	return idlecmd
}
