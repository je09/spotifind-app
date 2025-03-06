package main

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

var (
	logLocation = map[string]string{
		"darwin":  "/Library/Logs/Spotifind/spotifind.log",
		"linux":   "/.spotifind/spotifind.log",
		"windows": "\\AppData\\Roaming\\spotifind\\spotifind.log",
	}
)

var (
	// LogFileLocation needs to print location of the log in case of u
	LogFileLocation = ""
)

// Logger provides methods that call both console and file loggers.
type Logger struct {
	log *slog.Logger
}

func NewLogger() *Logger {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Error getting home directory", "error", err)
		os.Exit(1)
	}

	if _, ok := logLocation[runtime.GOOS]; !ok {
		panic("Unsupported OS")
	}

	path := homeDir + logLocation[runtime.GOOS]
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		slog.Error("Error creating log directory", "error", err)
		os.Exit(1)
	}
	LogFileLocation = path

	// Create a lumberjack logger for log rotation.
	rotatingLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}
	defer rotatingLogger.Close()

	// Create a multi-writer for console and file logging.
	multiWriter := io.MultiWriter(rotatingLogger, os.Stdout)

	// Create a logger that writes to both stdout and the log file.
	logger := slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("Saving logs to " + path)

	return &Logger{
		log: logger,
	}
}

func (l *Logger) Print(message string) {
	l.log.Info("PRINT: " + message)
}

func (l *Logger) Trace(message string) {
	l.log.Debug("TRACE: " + message)
}

func (l *Logger) Debug(message string) {
	l.log.Debug("DEBUG: " + message)
}

func (l *Logger) Info(message string) {
	l.log.Info("INFO: " + message)
}

func (l *Logger) Warning(message string) {
	l.log.Warn("WARNING: " + message)
}

func (l *Logger) Error(message string) {
	l.log.Error("ERROR: " + message)
}

func (l *Logger) Fatal(message string) {
	l.log.Error("FATAL: " + message)
	os.Exit(1)
}
