// proc offers control over arbitrary processes.
// proc depends heavily on os/exec.
// If/when proc needs to do resource control & isolation
// its probable that os/exec will need to be replaced with
// either the os package itself or syscall.
package proc

// Fingers crossed that its easier than rewriting os.ForkExec

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"
)

const outputReadRefreshInterval = time.Millisecond * 10

// Proc is a BitBox process.
// Stderr and Stdout are dumped to temp files on disk.
type Proc struct {
	// TODO: there is 99% chance we need a synchronization primitive at this level
	stdout *os.File
	stderr *os.File
	cmd    *exec.Cmd
}

// NewProc constructs and begins process.
// Stdout & Stderr of the new process are pointed at temp files.
// The tempfiles are acessable through the coresponding members.
func NewProc(cmdName string, args ...string) (Proc, error) {

	var err error
	var cmdPath string
	if cmdPath, err = exec.LookPath(cmdName); err != nil {
		return Proc{}, err
	}
	cmd := exec.Command(cmdPath, args...)

	var stdout *os.File
	if stdout, err = ioutil.TempFile("", ""); err != nil {
		return Proc{}, err
	}
	cmd.Stdout = stdout

	var stderr *os.File
	if stderr, err = ioutil.TempFile("", ""); err != nil {
		return Proc{}, err
	}
	cmd.Stderr = stderr

	if err = cmd.Start(); err != nil {
		return Proc{}, err
	}

	go func() {
		// Wait on the cmd to make sure resources get released
		cmd.Wait()
	}()

	return Proc{
		stdout: stdout,
		stderr: stderr,
		cmd:    cmd,
	}, nil
}

// Kill causes the running process to exit and closes the related resources.
// There is no guarantee that the process has actually exited when Kill returns.
// See documentation for os.Process.Kill()
func (p Proc) Kill() error {
	if err := p.cmd.Process.Kill(); err != nil {
		// processes which have ended return errors
		return err
	}
	// Only release resources if no err has been returned.
	// TODO: should we return the error or try to close both and return some conglomarate error?
	if err := p.stderr.Close(); err != nil {
		return err
	}
	if err := p.stdout.Close(); err != nil {
		return err
	}
	return nil
}

// Status returns the status of the process.
func (p Proc) Status() ProcStatus {
	if p.cmd.ProcessState == nil {
		return Running
	}

	if p.cmd.ProcessState.ExitCode() != 0 {
		return Stopped
	}
	return Exited
}

// ProcStatus is the status of a process.
type ProcStatus int

const (
	// Running indicates that the process is running.
	Running ProcStatus = iota
	// Exited indicates that the process returned a non-zero exit code.
	Exited
	// Stopped indicates that the process returned no exit code.
	Stopped
)

func (ps ProcStatus) String() string {
	return [...]string{"Running", "Exited", "Stopped"}[ps]
}

// Query streams output from the process to the returned channel.
// The Stdout and Stderr files are opened for reads and polled until
// a third routine finds that the process has exited.
// The third routine cancels the context of the pollReads.
// After the read routines finish the third routine sends the ExitCode and closes the channel.
func (p Proc) Query() (<-chan ProcOutput, error) {
	ctx, cancel := context.WithCancel(context.Background())
	stream := make(chan ProcOutput)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	if err := pollRead(ctx, p.stdout.Name(), wg, stream, newProcOutput_Stdout); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to setup poll read on stdout: %w", err)
	}
	if err := pollRead(ctx, p.stderr.Name(), wg, stream, newProcOutput_Stderr); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to setup poll read on stderr: %w", err)
	}

	go func() {
		// Throw away the error from Wait() because either:
		// 1. Wait() was already called (cool, we just use it to know that the process completed)
		// 2. The process returned a non-zero exit code (cool, we will return any exit code)
		p.cmd.Wait()
		cancel()
		wg.Wait()
		stream <- &ProcOutput_ExitCode{
			ExitCode: uint32(p.cmd.ProcessState.ExitCode()),
		}
		close(stream)
	}()
	return stream, nil
}

func pollRead(
	ctx context.Context,
	fileName string,
	wg *sync.WaitGroup,
	stream chan<- ProcOutput,
	packageOutput func([]byte) ProcOutput,
) error {
	flags := os.O_RDONLY | os.O_SYNC
	file, err := os.OpenFile(fileName, flags, 0600)
	if err != nil {
		return err
	}

	bufFile := bufio.NewReader(file)
	ticker := time.NewTicker(outputReadRefreshInterval)

	go func() {
	PollLoop:
		for {
			select {
			case <-ticker.C:
				// ReadLoop
				for {
					b, err := ioutil.ReadAll(bufFile)
					if err != nil {
						// TODO: should we log the error somehow?
						break PollLoop
					}
					if len(b) == 0 {
						break // Yes, only exit the inner loop. This is the only path back to the PollLoop
					}
					stream <- packageOutput(b)
				}
			case <-ctx.Done():
				// ReadLoop
				for {
					b, err := ioutil.ReadAll(bufFile)
					if err != nil {
						// TODO: should we log the error somehow?
						break PollLoop
					}
					if len(b) == 0 {
						break PollLoop
					}
					stream <- packageOutput(b)
				}
			}
		}
		wg.Done()
	}()
	return nil
}

// ProcOutput is any output from a process.
type ProcOutput interface {
	isProcOutput()
}

var _ ProcOutput = (*ProcOutput_Stdout)(nil)

// ProcOutput_Stdout is any output from the process which was written to Stdout.
type ProcOutput_Stdout struct {
	// Stdout is a series of characters sent to Stdout by a process.
	Stdout string
}

func newProcOutput_Stdout(b []byte) ProcOutput {
	return &ProcOutput_Stdout{
		Stdout: (string)(b),
	}
}

func (*ProcOutput_Stdout) isProcOutput() {}

var _ ProcOutput = (*ProcOutput_Stderr)(nil)

// ProcOutput_Stderr is any output from the process which was written to Stderr.
type ProcOutput_Stderr struct {
	// Stderr is a series of characters sent to Stderr by a process.
	Stderr string
}

func newProcOutput_Stderr(b []byte) ProcOutput {
	return &ProcOutput_Stderr{
		Stderr: (string)(b),
	}
}

func (*ProcOutput_Stderr) isProcOutput() {}

var _ ProcOutput = (*ProcOutput_ExitCode)(nil)

// ProcOutput_ExitCode is any output from the process which was written to Stderr.
type ProcOutput_ExitCode struct {
	// ExitCode is the exit code of a process.
	ExitCode uint32
}

func (*ProcOutput_ExitCode) isProcOutput() {}
