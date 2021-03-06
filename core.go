package bitbox

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jmbarzee/bitbox/proc"
)

// Core offers the central functionality of BitBox.
// Core supports basic process control and interaction.
type Core struct {
	sync.RWMutex
	processes map[uuid.UUID]*proc.Proc
}

// Start initiates a process.
func (c *Core) Start(cmd string, args ...string) (uuid.UUID, error) {

	id := uuid.New()
	proc, err := proc.NewProc(cmd, args...)
	if err != nil {
		return uuid.UUID{}, c.newError("Start", err)
	}

	c.Lock()
	c.processes[id] = proc // Chance of colision (16 byte id, so roughly 2^128 chance)
	c.Unlock()
	return id, nil
}

// Stop halts a process.
func (c *Core) Stop(id uuid.UUID) error {
	var p *proc.Proc
	var err error

	if p, err = c.findProcess(id); err != nil {
		return c.newError("Stop", err)
	}
	if err = p.Stop(); err != nil {
		return c.newError("Stop", err)
	}
	return nil
}

// Status returns the status of the process.
func (c *Core) Status(id uuid.UUID) (proc.ProcStatus, error) {
	var p *proc.Proc
	var err error

	if p, err = c.findProcess(id); err != nil {
		return proc.Exited, c.newError("Status", err)
	}

	return p.Status(), nil
}

// Query streams the output/result of a process.
func (c *Core) Query(id uuid.UUID) (<-chan proc.ProcOutput, error) {
	var p *proc.Proc
	var err error

	if p, err = c.findProcess(id); err != nil {
		return nil, c.newError("Query", err)
	}

	return p.Query()
}

func (c *Core) findProcess(id uuid.UUID) (*proc.Proc, error) {
	c.RLock()
	defer c.RUnlock()
	p, ok := c.processes[id]
	if !ok {
		return nil, fmt.Errorf("could not find specified process %v", id)
	}
	return p, nil
}

func (*Core) newError(action string, err error) error {
	return fmt.Errorf("could not %v process: %w", action, err)
}
