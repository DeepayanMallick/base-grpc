package logging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type FluentdFormatter struct {
	TimestampFormat string
}

func (f *FluentdFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		if err, ok := v.(error); ok {
			v = err.Error()
		}
		data[k] = v
	}
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	data["time"] = entry.Time.Format(timestampFormat)
	data["message"] = entry.Message
	data["severity"] = entry.Level.String()

	if entry.HasCaller() {
		data["func"] = entry.Caller.Function
		data["file"] = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON: %v", err)
	}
	return append(serialized, '\n'), nil
}

func prefixFieldClashes(data logrus.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}
	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}
	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}
