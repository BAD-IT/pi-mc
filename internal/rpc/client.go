package rpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type PiRpcClient struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	events  chan Event
	errChan chan error
}

func NewPiRpcClient(binPath string, args ...string) (*PiRpcClient, error) {
	cmd := exec.Command(binPath, args...)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	cmd.Stderr = os.Stderr
	
	return &PiRpcClient{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		events:  make(chan Event, 100),
		errChan: make(chan error, 1),
	}, nil
}

func (c *PiRpcClient) Start() error {
	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("start cmd: %w", err)
	}
	
	go c.readLoop()
	return nil
}

func (c *PiRpcClient) Stop() error {
	if c.cmd.Process != nil {
		return c.cmd.Process.Kill()
	}
	return nil
}

func (c *PiRpcClient) SendCommand(cmd Command) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	
	// Write JSON-Lines
	_, err = c.stdin.Write(append(b, '\n'))
	return err
}

func (c *PiRpcClient) Events() <-chan Event {
	return c.events
}

func (c *PiRpcClient) Errors() <-chan error {
	return c.errChan
}

func (c *PiRpcClient) readLoop() {
	defer close(c.events)
	
	scanner := bufio.NewScanner(c.stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			c.errChan <- fmt.Errorf("parse event: %w. line: %s", err, line)
			continue
		}
		
		c.events <- event
	}
	
	if err := scanner.Err(); err != nil {
		c.errChan <- fmt.Errorf("scanner error: %w", err)
	}
	
	c.cmd.Wait()
}
