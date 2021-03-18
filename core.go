package bitbox

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jmbarzee/bitbox/proc"
)

// Core offers the central functionality of BitBox.
// Core supports basic process control and interaction.
type Core struct {
	m         sync.Mutex
	processes map[uuid.UUID]proc.Proc
}

// Start initiates a process.
func (c *Core) Start(cmd string, args ...string) (uuid.UUID, error) {

	id := uuid.New()
	proc, err := proc.NewProc(cmd, args...)
	if err != nil {
		return uuid.UUID{}, c.newError("Start", err)
	}

	c.m.Lock()
	c.processes[id] = proc // Chance of colision (16 byte id, so roughly 2^128 chance)
	c.m.Unlock()
	return id, nil
}

// Stop halts a process.
func (c *Core) Stop(id uuid.UUID) error {
	var p proc.Proc
	var err error

	if p, err = c.findProcess(id); err != nil {
		c.newError("Stop", err)
	}
	if err = p.Kill(); err != nil {
		c.newError("Stop", err)
	}
	return nil
}

// Status returns the status of the process.
	return Stopped, errors.New("unimplemented")
func (c *Core) Status(id uuid.UUID) (proc.ProcStatus, error) {
}

// Query streams the output/result of a process.
func (c *Core) Query(id uuid.UUID) (chan<- ProcOutput, error) {
	return nil, errors.New("unimplemented")
}

func (c *Core) findProcess(id uuid.UUID) (proc.Proc, error) {
	c.m.Lock()
	p, ok := c.processes[id]
	if !ok {
		return proc.Proc{}, fmt.Errorf("could not find specified process %v", id)
	}
	c.m.Lock()
	return p, nil
}

func (*Core) newError(action string, err error) error {
	return fmt.Errorf("could not %v process: %w", action, err)
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
