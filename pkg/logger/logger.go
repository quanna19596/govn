package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level       string
	Filename    string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
	Environment string
}

func NewLogger(config LoggerConfig) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(config.Level)

	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)

	var writer io.Writer

	if config.Environment == "development" {
		writer = PrettyJSONWriter{Writer: os.Stdout}
	} else {
		writer = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize, // megabytes
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge, //days
			Compress:   config.Compress,
		}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()
	return &logger
}

type PrettyJSONWriter struct {
	Writer io.Writer
}

func (pjw PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, p, "", "  ")
	if err != nil {
		return pjw.Writer.Write(p)
	}

	return pjw.Writer.Write(prettyJSON.Bytes())
}
