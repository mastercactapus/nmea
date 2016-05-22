package nmea

import (
	"fmt"
	"strconv"
)

// GPGSAFix is the fix type for a GPGSA sentence
type GPGSAFix string

// Fix types for GPGSA
const (
	GPGSAFixNoFix GPGSAFix = "1"
	GPGSAFix2D    GPGSAFix = "2"
	GPGSAFix3D    GPGSAFix = "3"
)

// GPGSA is used to communicate dilution of precision and active satellites
type GPGSA struct {
	AutoSelection bool     // specifies if selection of 2D vs 3D fix is automatic or manual
	FixType       GPGSAFix // the type of fix the receiver has
	Satellites    []string // PRNs of satellites used for fix. Maximum of 12
	PDOP          float64  // dilution of precision
	HDOP          float64  // horizontal dilution of precision
	VDOP          float64  // vertical dilution of precision
}

// Type returns TypeGPGSA to fulfill the Sentence interface
func (g GPGSA) Type() Type {
	return TypeGPGSA
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

	switch GPGSAFix(r.Fields[1]) {
	case GPGSAFix2D, GPGSAFix3D, GPGSAFixNoFix:
		g.FixType = GPGSAFix(r.Fields[1])
	case GPGSAFix(""):
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
	g.PDOP, err = parseFieldFloat(r.Fields[14], "PDOP")
	if err != nil {
		return err
	}

	g.HDOP, err = parseFieldFloat(r.Fields[15], "HDOP")
	if err != nil {
		return err
	}

	g.VDOP, err = parseFieldFloat(r.Fields[16], "VDOP")
	if err != nil {
		return err
	}

	return nil
}
