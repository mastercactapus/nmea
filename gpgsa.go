package nmea

import (
	"fmt"
	"strconv"
)

// GPGSAFix is the fix type for a GPGSA sentence
type GPGSAFix string

const (
	GPGSAFixNoFix GPGSAFix = "1"
	GPGSAFix2D    GPGSAFix = "2"
	GPGSAFix3D    GPGSAFix = "3"
)

// GPGSA is used to communicate dilution of precision and active satellites
type GPGSA struct {
	// AutoSelection specifies if selection of 2D vs 3D fix is automatic or manual
	AutoSelection bool

	// Fix type is the type of fix the receiver has
	FixType GPGSAFix

	// PRNs of satellites used for fix. Maximum of 12
	Satellites []string

	// PDOP is the dilution of precision
	PDOP float64

	// HDOP is the horizontal dilution of precision
	HDOP float64

	// VDOP is the vertical dilution of precision
	VDOP float64
}

// String will provide a NMEA formatted string. If more than 12 Satellites are present, only the first 12 will be serialized
func (g GPGSA) String() string {
	r := &Raw{TypeName: string(TypeGPGSA)}
	r.Fields = make([]string, 17)

	if g.AutoSelection {
		r.Fields[0] = "A"
	} else {
		r.Fields[0] = "M"
	}

	r.Fields[1] = string(g.FixType)

	if g.Satellites != nil {
		copy(r.Fields[2:14], g.Satellites)
	}

	r.Fields[14] = strconv.FormatFloat(g.PDOP, 'f', -1, 64)
	r.Fields[15] = strconv.FormatFloat(g.HDOP, 'f', -1, 64)
	r.Fields[16] = strconv.FormatFloat(g.VDOP, 'f', -1, 64)

	return r.String()
}

// Parse will parse GPGSA data from a raw sentence struct
func (g *GPGSA) Parse(r *Raw) error {
	if r.TypeName != string(TypeGPGSA) {
		return fmt.Errorf("wrong type for GPGSA '%s'", r.TypeName)
	}
	if r.Fields == nil || len(r.Fields) < 17 {
		return fmt.Errorf("not enough fields, need at least 17")
	}

	switch r.Fields[0] {
	case "", "M":
		g.AutoSelection = false
	case "A":
		g.AutoSelection = true
	default:
		return fmt.Errorf("invalid selection type: %s", r.Fields[0])
	}

	switch r.Fields[1] {
	case string(GPGSAFix2D), string(GPGSAFix3D), string(GPGSAFixNoFix):
		g.FixType = GPGSAFix(r.Fields[1])
	case "":
		g.FixType = GPGSAFixNoFix
	default:
		return fmt.Errorf("invalid fix type: %s", r.Fields[1])
	}

	g.Satellites = make([]string, 0, 12)
	for _, sat := range r.Fields[2:14] {
		if sat == "" {
			continue
		}
		g.Satellites = append(g.Satellites, sat)
	}

	var err error
	g.PDOP, err = strconv.ParseFloat(r.Fields[14], 64)
	if err != nil {
		return fmt.Errorf("parse PDOP: %s", err)
	}
	g.HDOP, err = strconv.ParseFloat(r.Fields[15], 64)
	if err != nil {
		return fmt.Errorf("parse HDOP: %s", err)
	}
	g.VDOP, err = strconv.ParseFloat(r.Fields[16], 64)
	if err != nil {
		return fmt.Errorf("parse VDOP: %s", err)
	}

	return nil
}
