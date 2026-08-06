package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type SyslogServer struct {
	bindAddr   string
	bufferDir  string
	normalizer *SyslogNormalizer
	ingestFunc func([]LogEntry) error
	listener   *net.UDPConn
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

func NewSyslogServer(bindAddr, bufferDir string, ingestFunc func([]LogEntry) error) *SyslogServer {
	return &SyslogServer{
		bindAddr:   bindAddr,
		bufferDir:  bufferDir,
		normalizer: NewSyslogNormalizer(),
		ingestFunc: ingestFunc,
		stopCh:     make(chan struct{}),
	}
}

func (s *SyslogServer) Start() error {
	if err := os.MkdirAll(s.bufferDir, 0755); err != nil {
		return fmt.Errorf("create buffer dir: %w", err)
	}

	addr, err := net.ResolveUDPAddr("udp", s.bindAddr)
	if err != nil {
		return fmt.Errorf("resolve udp addr: %w", err)
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("listen udp: %w", err)
	}
	s.listener = listener

	s.wg.Add(1)
	go s.readLoop()

	log.Printf("[syslog] Server started on %s", s.bindAddr)
	return nil
}

func (s *SyslogServer) readLoop() {
	defer s.wg.Done()
	buf := make([]byte, 4096)
	for {
		select {
		case <-s.stopCh:
			return
		default:
		}

		s.listener.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, srcAddr, err := s.listener.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			log.Printf("[syslog] read error: %v", err)
			continue
		}

		rawMsg := string(buf[:n])
		entry := s.normalizer.Normalize(rawMsg, srcAddr.IP.String())
		if entry != nil {
			s.bufferEntry(*entry)
		}
	}
}

func (s *SyslogServer) bufferEntry(entry LogEntry) {
	entry.Time = time.Now().UTC().Format(time.RFC3339)

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[syslog] marshal entry: %v", err)
		return
	}

	filename := fmt.Sprintf("%s_%d.json", time.Now().UTC().Format("20060102150405"), time.Now().UTC().UnixNano())
	filePath := filepath.Join(s.bufferDir, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Printf("[syslog] write buffer: %v", err)
		return
	}

	s.flushBuffer()
}

func (s *SyslogServer) flushBuffer() {
	entries, err := s.readBufferDir()
	if err != nil || len(entries) == 0 {
		return
	}

	var slice []LogEntry
	for _, entry := range entries {
		slice = append(slice, entry)
	}

	if err := s.ingestFunc(slice); err != nil {
		log.Printf("[syslog] ingest failed: %v", err)
		return
	}

	for fname := range entries {
		fpath := filepath.Join(s.bufferDir, fname)
		os.Remove(fpath)
	}
}

func (s *SyslogServer) readBufferDir() (map[string]LogEntry, error) {
	result := make(map[string]LogEntry)

	files, err := filepath.Glob(filepath.Join(s.bufferDir, "*.json"))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		result[filepath.Base(f)] = entry
	}

	return result, nil
}

func (s *SyslogServer) Stop() {
	close(s.stopCh)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	log.Printf("[syslog] Server stopped")
}
