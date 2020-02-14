package fnlog

import (
	"encoding/json"
	"fmt"
	"time"
)

// JSONFormatter - log with json format
type JSONFormatter struct {
	Timeformat string
}

// Message - json message
func (p *JSONFormatter) Message(level LogLevel, fieldMap fields, args ...interface{}) string {
	s := p.defaultLog(level)
	if fieldMap != nil {
		s = p.logWithField(s, fieldMap)
	}

	if args != nil {
		s += fmt.Sprintf(`"message":"%+v",`, args...)
	}

	s = s[:len(s)-1] + "}\n"

	return s
}

func (p *JSONFormatter) defaultLog(level LogLevel) string {
	needCaller := false
	s := "{"

	if level != accessLevel {
		needCaller = true
	}

	s += `"severity":"` + levelText[level] + `",`
	s += fmt.Sprintf(`"timestamp":"%v",`, time.Now().Format(p.Timeformat))

	if needCaller {
		s += fmt.Sprintf(`"caller":"%v",`, GetCaller())
	}

	return s
}

func (p *JSONFormatter) logWithField(s string, f fields) string {
	for k, v := range f {
		switch v.(type) {
		case int, int8, int16, int32, int64:
			s += fmt.Sprintf(`"%v":%v,`, k, v)
		case uint, uint8, uint16, uint32, uint64:
			s += fmt.Sprintf(`"%v":%v,`, k, v)
		case float32, float64:
			s += fmt.Sprintf(`"%v":%v,`, k, v)
		case error:
			s += fmt.Sprintf(`"%v":"%v",`, k, v.(error).Error())
		case string:
			s += fmt.Sprintf(`"%v":"%v",`, k, v)
		default:
			mar, err := json.Marshal(v)
			if err != nil {
				break
			}
			s += fmt.Sprintf(`"%v":%v,`, k, string(mar))
		}
	}

	return s
}
