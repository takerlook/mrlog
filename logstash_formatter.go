package logrustash

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// LogFormat generates json in logstash format.
// Logstash site: http://logstash.net/
type LogFormat struct {
	// if not empty use for logstash type field.
	Type string
	// TimestampFormat sets the format used for timestamps.
	TimestampFormat string
}

// Format formats log message.
func (f *LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	return f.FormatWithPrefix(entry, "")
}

// FormatWithPrefix removes prefix from keys and formats log message.
func (f *LogFormat) FormatWithPrefix(entry *logrus.Entry, prefix string) ([]byte, error) {
	fields := make(logrus.Fields)
	for k, v := range entry.Data {
		// Remove the prefix when sending the fields to logstash
		if prefix != "" && strings.HasPrefix(k, prefix) {
			k = strings.TrimPrefix(k, prefix)
		}

		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/377
			fields[k] = v.Error()
		default:
			fields[k] = v
		}
	}

	timeStampFormat := f.TimestampFormat

	if timeStampFormat == "" {
		timeStampFormat = time.RFC3339Nano
	}

	fields["@timestamp"] = entry.Time.Format(timeStampFormat)

	// set message field
	v, ok := entry.Data["message"]
	if ok {
		fields["fields.message"] = v
	}
	fields["message"] = entry.Message

	// set level field
	v, ok = entry.Data["level"]
	if ok {
		fields["fields.level"] = v
	}
	fields["level"] = entry.Level.String()

	// set type field
	if f.Type != "" {
		v, ok = entry.Data["type"]
		if ok {
			fields["fields.type"] = v
		}
		fields["api_id"] = f.Type
	}

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
