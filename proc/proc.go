package proc

import (
	"io/ioutil"
	"os"
	"os/exec"
)

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
func NewProc(cmdName string, args ...string) (*Proc, error) {

	var err error
	var cmdPath string
	if cmdPath, err = exec.LookPath(cmdName); err != nil {
		return nil, err
	}
	cmd := exec.Command(cmdPath, args...)

	var stdout *os.File
	if stdout, err = ioutil.TempFile("", ""); err != nil {
		return nil, err
	}
	cmd.Stdout = stdout

	var stderr *os.File
	if stderr, err = ioutil.TempFile("", ""); err != nil {
		return nil, err
	}
	cmd.Stderr = stderr

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	return &Proc{
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
