package main

import (
	"fmt"
	"runtime"
)

type TerminalServer struct {
	container string
	wsHandler func([]byte) []byte
}

func NewTerminalServer() *TerminalServer {
	return &TerminalServer{}
}

func (t *TerminalServer) Start(container string) error {
	t.container = container

	if runtime.GOOS == "windows" {
		return t.startWindows()
	}
	return t.startLinux()
}

func (t *TerminalServer) startLinux() error {
	return fmt.Errorf("terminal server not yet implemented")
}

func (t *TerminalServer) startWindows() error {
	return fmt.Errorf("windows terminal not yet implemented")
}

func (t *TerminalServer) HandleWebSocket(msg []byte) []byte {
	return msg
}

func (t *TerminalServer) Stop() error {
	return nil
}

type PTYServer struct {
	container string
	wsPath    string
}

func NewPTYServer(container, wsPath string) *PTYServer {
	return &PTYServer{
		container: container,
		wsPath:    wsPath,
	}
}

func (p *PTYServer) Start() error {
	if runtime.GOOS == "windows" {
		return p.startConPTY()
	}
	return p.startLinuxPTY()
}

func (p *PTYServer) startLinuxPTY() error {
	return fmt.Errorf("pty server not yet implemented")
}

func (p *PTYServer) startConPTY() error {
	return fmt.Errorf("conpty not yet implemented")
}

func (p *PTYServer) WriteToPTY(data []byte) (int, error) {
	return len(data), nil
}

func (p *PTYServer) ReadFromPTY(buf []byte) (int, error) {
	return 0, nil
}
