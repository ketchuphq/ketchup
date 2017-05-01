package memlogger

import (
	"fmt"

	"github.com/octavore/nagax/logger"
)

var _ logger.Logger = &MemoryLogger{}

type MemoryLogger struct {
	Infos    []string
	Warnings []string
	Errors   []string
}

func (m *MemoryLogger) Reset() {
	m.Infos = []string{}
	m.Warnings = []string{}
	m.Errors = []string{}
}

func (m *MemoryLogger) Count() int {
	return len(m.Infos) + len(m.Warnings) + len(m.Errors)
}

func (m *MemoryLogger) Info(args ...interface{}) {
	m.Infos = append(m.Infos, fmt.Sprint(args...))
}

func (m *MemoryLogger) Infof(format string, args ...interface{}) {
	m.Infos = append(m.Infos, fmt.Sprintf(format, args...))
}

func (m *MemoryLogger) Warning(args ...interface{}) {
	m.Warnings = append(m.Warnings, fmt.Sprint(args...))
}

func (m *MemoryLogger) Warningf(format string, args ...interface{}) {
	m.Warnings = append(m.Warnings, fmt.Sprintf(format, args...))
}

func (m *MemoryLogger) Error(args ...interface{}) {
	m.Errors = append(m.Errors, fmt.Sprint(args...))
}

func (m *MemoryLogger) Errorf(format string, args ...interface{}) {
	m.Errors = append(m.Errors, fmt.Sprintf(format, args...))
}
