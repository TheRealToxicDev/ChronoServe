package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level      LogLevel
	logger     *log.Logger
	file       *os.File
	maxSize    int64
	maxBackups int
	directory  string
	filename   string
	mu         sync.Mutex
}

// LoggerOptions holds configuration for logger initialization
type LoggerOptions struct {
	Level      LogLevel
	Directory  string
	Filename   string
	MaxSize    int // Size in MB
	MaxBackups int
}

func NewLogger(opts LoggerOptions) (*Logger, error) {
	logger := &Logger{
		level:      opts.Level,
		maxSize:    int64(opts.MaxSize) * 1024 * 1024, // Convert MB to bytes
		maxBackups: opts.MaxBackups,
		directory:  opts.Directory,
		filename:   opts.Filename,
	}

	if err := logger.initialize(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) initialize() error {
	if l.directory != "" {
		if err := os.MkdirAll(l.directory, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		logPath := filepath.Join(l.directory, l.filename)
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}

		l.file = file
		l.logger = log.New(io.MultiWriter(os.Stdout, file), "", 0)
	} else {
		l.logger = log.New(os.Stdout, "", 0)
	}

	return nil
}

func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level >= l.level {
		// Get caller information
		_, file, line, _ := runtime.Caller(2)

		// Format the message
		msg := fmt.Sprintf(format, v...)
		timestamp := time.Now().Format("2006-01-02 15:04:05.000")
		logEntry := fmt.Sprintf("%s [%s] %s:%d: %s\n",
			timestamp,
			level.String(),
			filepath.Base(file),
			line,
			msg,
		)

		// Write to log
		l.logger.Print(logEntry)

		// Check if rotation is needed
		if l.file != nil {
			if info, err := l.file.Stat(); err == nil && info.Size() > l.maxSize {
				l.rotate()
			}
		}
	}
}

func (l *Logger) rotate() {
	if l.file == nil {
		return
	}

	// Close current file
	l.file.Close()

	// Generate timestamp for backup
	timestamp := time.Now().Format("20060102-150405")
	currentPath := filepath.Join(l.directory, l.filename)
	backupPath := filepath.Join(l.directory, fmt.Sprintf("%s.%s", l.filename, timestamp))

	// Rename current file to backup
	os.Rename(currentPath, backupPath)

	// Create new log file
	l.initialize()

	// Clean old backups
	l.cleanOldBackups()
}

func (l *Logger) cleanOldBackups() {
	if l.maxBackups <= 0 {
		return
	}

	pattern := filepath.Join(l.directory, l.filename+".*")
	matches, _ := filepath.Glob(pattern)

	if len(matches) > l.maxBackups {
		// Sort files by modification time
		type backup struct {
			path    string
			modTime time.Time
		}
		backups := make([]backup, 0, len(matches))

		for _, path := range matches {
			if info, err := os.Stat(path); err == nil {
				backups = append(backups, backup{path, info.ModTime()})
			}
		}

		// Sort by modification time (oldest first)
		for i := 0; i < len(backups)-l.maxBackups; i++ {
			os.Remove(backups[i].path)
		}
	}
}

// Sync ensures all log entries are written to their destination
// This is especially important to call before program exit
func (l *Logger) Sync() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// For file-based logging, we need to sync the file
	if l.file != nil {
		if err := l.file.Sync(); err != nil {
			return fmt.Errorf("failed to sync log file: %w", err)
		}
	}

	return nil
}

// Close properly closes the logger
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// GetLogLevel converts string level to LogLevel
func GetLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// Logger methods for different log levels
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}
