package logger

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"strings"
)

// TextFormatter 文本格式
type TextFormatter struct {
	ColorPrint bool // 是否打印颜色
}

// colors
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

const TimeFormat = "2006-01-02T15:04:05.000Z07:00"

func (f *TextFormatter) formatLevel(level logrus.Level) string {
	levelTxt := fmt.Sprintf("%-7s", strings.ToUpper(level.String())) // align level
	if f.ColorPrint {
		var levelColor int
		switch level {
		case logrus.DebugLevel, logrus.TraceLevel:
			levelColor = gray
		case logrus.WarnLevel:
			levelColor = yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = red
		default:
			levelColor = blue
		}
		levelTxt = fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, levelTxt)
	}
	return levelTxt
}

// Format implements logurs.TextFormatter
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var msg string
	// format level
	msg += f.formatLevel(entry.Level)
	// format timestamp
	msg += fmt.Sprintf(" [%s]", entry.Time.Format(TimeFormat))
	// append message
	if entry.Message != "" {
		msg += fmt.Sprintf(" %s", entry.Message)
	}
	// format entry data
	if data := f.formatEntryData(entry.Data); len(data) != 0 {
		dataStr, err := jsoniter.MarshalToString(data)
		if err != nil {
			return nil, err
		}
		msg += fmt.Sprintf(" | %s", dataStr)
	}
	// end the message with \n
	msg += "\n"
	return []byte(msg), nil
}

func (f *TextFormatter) formatEntryData(data logrus.Fields) map[string]interface{} {
	if len(data) == 0 {
		return nil
	}
	dataMap := make(map[string]interface{})
	for k, v := range data {
		switch v := v.(type) {
		case error:
			dataMap[k] = v.Error()
		case fmt.Stringer:
			dataMap[k] = v.String()
		default:
			dataMap[k] = v
		}
	}
	return dataMap
}
