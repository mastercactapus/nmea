package nmea

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Type corresponds to a sentence type
type Type string

// ErrUnknownType is used when a sentence type is unknown or currently unsupported
var ErrUnknownType = errors.New("unknown sentence type")

// Supported NMEA sentence types
const (
	TypeGPRMC Type = "GPRMC"
	TypeGPGSA Type = "GPGSA"
)

// Sentence is a NMEA sentence
type Sentence interface {
	Type() Type
	String() string
}

// Raw is a NMEA sentence that has been broken up into its TypeName and Fields. Checksums are handled automatically
type Raw struct {
	TypeName string
	Fields   []string
}

func (r Raw) String() string {
	data := r.TypeName
	if r.Fields != nil && len(r.Fields) > 0 {
		data += "," + strings.Join(r.Fields, ",")
	}
	check := Checksum([]byte(data))

	return fmt.Sprintf("$%s*%02X", data, check)
}

// Checksum will calculate a NMEA checksum of data
func Checksum(p []byte) byte {
	var sum byte
	for _, b := range p {
		sum ^= b
	}
	return sum
}

// ParseRaw will return a Raw struct, validating checksum (if any) and separating individual fields and the type
func ParseRaw(line []byte) (*Raw, error) {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return nil, io.ErrShortBuffer
	}
	if line[0] != '$' {
		return nil, fmt.Errorf("expected '$' but got '%s'", string(line[0]))
	}
	line = line[1:]
	if len(line) >= 3 && line[len(line)-3] == '*' {
		checkStr := bytes.ToUpper(line[len(line)-2:])
		check := (checkStr[0]-48)*16 + (checkStr[1] - 48)
		line = line[:len(line)-3]
		if Checksum(line) != check {
			return nil, fmt.Errorf("checksum: expected 0x%02x but found 0x%02x", Checksum(line), check)
		}
	}
	fields := make([]string, 0, 20)
	for {
		i := bytes.IndexByte(line, ',')
		if i == -1 {
			fields = append(fields, string(line))
			break
		}
		fields = append(fields, string(line[:i]))
		line = line[i+1:]
	}
	return &Raw{TypeName: fields[0], Fields: fields[1:]}, nil
}

// Parse will return a struct for the line type. If type is unknown, ErrUnknownType will be returned.
func Parse(line []byte) (Sentence, error) {
	r, err := ParseRaw(line)
	if err != nil {
		return nil, err
	}
	switch Type(r.TypeName) {
	case TypeGPRMC:
		s := new(GPRMC)
		return s, s.Parse(r)
	default:
		return nil, ErrUnknownType
	}
}
