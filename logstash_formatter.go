package mrlog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// APIIP 로그를 전송하는 API의 IP
	APIIP = "host"
	// APPID 로그를 보내는 APP의 ID
	APPID = "appid"
	// LEVEL 로그 레벨
	LEVEL = "level"
	// TRACEID 로그 추적 ID
	TRACEID = "traceid"
	// UIP API를 사용한 USER IP
	UIP = "uip"
	// UID API를 사용한 USER ID
	UID = "uid"
	// CLIENTID 유저가 사용하는 Cliend id
	CLIENTID = "clientid"
	// METHOD 요청 들어온 Method
	METHOD = "method"
	// PATH 요청 들어온 Path
	PATH = "path"
	// HEADER Request의 Header
	HEADER = "header"
	// SESSIONKEY 세션 키
	SESSIONKEY = "sessionkey"
	// DEFAULTFIELD 기본 필드
	DEFAULTFIELD = "defaultField"
)

// LogFormat generates json in logstash format.
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
		if prefix != "" && strings.HasPrefix(k, prefix) {
			k = strings.TrimPrefix(k, prefix)
		}

		switch v := v.(type) {
		case error:
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
	fields["nano_sec"] = entry.Time.Format(time.RFC3339Nano)

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
		fields["appid"] = f.Type
	}

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
