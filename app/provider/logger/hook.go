package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// StdoutHook 标准输出hook
type StdoutHook struct {
	Writer    io.Writer
	Level     logrus.Level
	Formatter logrus.Formatter
}

// NewStdoutHook 创建标准输出hook
func NewStdoutHook(level logrus.Level, fmtter logrus.Formatter) *StdoutHook {
	return &StdoutHook{
		Writer:    os.Stdout,
		Level:     level,
		Formatter: fmtter,
	}
}

// Levels 日志等级
func (hook *StdoutHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, len(logrus.AllLevels))
	for _, l := range logrus.AllLevels {
		if l <= hook.Level {
			levels = append(levels, l)
		}
	}
	return levels
}

// Fire 记录日志
func (hook *StdoutHook) Fire(entry *logrus.Entry) error {
	msg, err := hook.Formatter.Format(entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to format entry: %s", err)
		return err
	}
	if _, err := fmt.Fprintf(hook.Writer, "%s", msg); err != nil {
		return err
	}
	return nil
}
