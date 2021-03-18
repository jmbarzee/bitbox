package bitbox

import (
	"errors"

	"github.com/google/uuid"
)

// Core offers the central functionality of BitBox.
// Core supports basic process control and interaction.
type Core struct {
}

// Start initiates a process.
func (c *Core) Start(cmd string, params []string) (uuid.UUID, error) {
	return uuid.UUID{}, errors.New("unimplemented")
}

// Stop halts a process.
func (c *Core) Stop(id uuid.UUID) error {
	return errors.New("unimplemented")
}

// Status returns the status of the process.
func (c *Core) Status(id uuid.UUID) (ProcStatus, error) {
	return Stopped, errors.New("unimplemented")
}

// Query streams the output/result of a process.
func (c *Core) Query(id uuid.UUID) (chan<- ProcOutput, error) {
	return nil, errors.New("unimplemented")
}

// ProcStatus is the status of a process.
type ProcStatus int

const (
	// Stopped indicates that the process is not running.
	Stopped ProcStatus = iota
	// Running indicates that the process is running.
	Running
)

func (ps ProcStatus) String() string {
	return [...]string{"Stopped", "Running"}[ps]
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

func (*ProcOutput_Stdout) isProcOutput() {}

var _ ProcOutput = (*ProcOutput_Stderr)(nil)

// ProcOutput_Stderr is any output from the process which was written to Stderr.
type ProcOutput_Stderr struct {
	// Stderr is a series of characters sent to Stderr by a process.
	Stderr string
}

func (*ProcOutput_Stderr) isProcOutput() {}

var _ ProcOutput = (*ProcOutput_ExitCode)(nil)

// ProcOutput_ExitCode is any output from the process which was written to Stderr.
type ProcOutput_ExitCode struct {
	// ExitCode is the exit code of a process.
	ExitCode uint32
}

func (*ProcOutput_ExitCode) isProcOutput() {}
