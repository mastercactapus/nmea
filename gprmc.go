package nmea

import (
	"fmt"
	"strconv"
	"time"
)

const timeFormat = "150405"
const dateFormat = "020106"

// GPRMCFix (fix type for GPRMC) was added in NMEA version 2.3 and indicates the kind of fix the receiver currently has. Only Autonomous and Differential represent a valid (Active) signal
type GPRMCFix string

// Fix types for GPRMC
const (
	GPRMCFixUnspecified  GPRMCFix = ""
	GPRMCFixAutonomous   GPRMCFix = "A"
	GPRMCFixDifferential GPRMCFix = "D"
	GPRMCFixEstimated    GPRMCFix = "E"
	GPRMCFixNotValid     GPRMCFix = "N"
	GPRMCFixSimulator    GPRMCFix = "S"
)

// GPRMC represents a GPRMC type NMEA sentence
type GPRMC struct {
	Time       time.Time // the time/date of the fix
	Active     bool      // true if the unit reports the fix as valid/active (Void otherwise)
	Latitude   Coord
	Longitude  Coord
	Speed      float64  // Speed in knots
	TrueCourse float64  // track made good in degrees True
	Variation  Coord    //  magnetic variation
	FixType    GPRMCFix // type of fix the receiver has
}

// Type returns TypeGPRMC to fulfill the Sentence interface
func (g GPRMC) Type() Type {
	return TypeGPRMC
}

// String will return a NMEA formatted string-representation of the GPRMC data
func (g GPRMC) String() string {
	t := g.Time.Format(timeFormat + dateFormat)
	stat := 'V'
	if g.Active {
		stat = 'A'
	}

	return Raw{
		TypeName: string(TypeGPRMC),
		Fields: []string{
			t[:6],
			string(stat),
			g.Latitude.String(),
			g.Latitude.Direction().LatString(),
			g.Longitude.String(),
			g.Longitude.Direction().LongString(),
			strconv.FormatFloat(g.Speed, 'f', -1, 64),
			strconv.FormatFloat(g.TrueCourse, 'f', -1, 64),
			t[6:],
			g.Variation.String(),
			g.Variation.Direction().LongString(),
			string(g.FixType),
		},
	}.String()
}

// Parse will parse GPRMC data from a raw sentence struct
func (g *GPRMC) Parse(r *Raw) error {
	if r.TypeName != string(TypeGPRMC) {
		return fmt.Errorf("wrong type for GPRMC '%s'", r.TypeName)
	}
	if r.Fields == nil || len(r.Fields) < 11 {
		return fmt.Errorf("not enough fields, need at least 11")
	}
	var err error
	if r.Fields[0] != "" {
		g.Time, err = time.ParseInLocation(timeFormat, r.Fields[0], time.UTC)
		if err != nil {
			return fmt.Errorf("parse timestamp: %s", err)
		}
	} else {
		g.Time = time.Time{}
	}

	if r.Fields[1] != "" {
		switch r.Fields[1] {
		case "A":
			g.Active = true
		case "V":
			g.Active = false
		default:
			return fmt.Errorf("invalid status value: %s", r.Fields[1])
		}
	} else {
		g.Active = false
	}

	g.Latitude, err = parseFieldCoord(r.Fields[2], r.Fields[3], "latitude")
	if err != nil {
		return err
	}
	g.Longitude, err = parseFieldCoord(r.Fields[4], r.Fields[5], "longitude")
	if err != nil {
		return err
	}

	if r.Fields[6] != "" {
		g.Speed, err = strconv.ParseFloat(r.Fields[6], 64)
		if err != nil {
			return fmt.Errorf("parse speed: %s", err)
		}
	} else {
		g.Speed = 0
	}

	if r.Fields[7] != "" {
		g.TrueCourse, err = strconv.ParseFloat(r.Fields[7], 64)
		if err != nil {
			return fmt.Errorf("parse true course: %s", err)
		}
	} else {
		g.TrueCourse = 0
	}

	var t time.Time
	if r.Fields[8] != "" {
		t, err = time.ParseInLocation(dateFormat, r.Fields[8], time.UTC)
		if err != nil {
			return fmt.Errorf("parse date: %s", err)
		}
		// add date & time
		g.Time = t.Add(time.Hour*time.Duration(g.Time.Hour()) + time.Minute*time.Duration(g.Time.Minute()) + time.Second*time.Duration(g.Time.Second()) + time.Duration(g.Time.Nanosecond()))
	}

	g.Variation, err = parseFieldCoord(r.Fields[9], r.Fields[10], "variation")
	if err != nil {
		return err
	}

	if len(r.Fields) >= 12 && r.Fields[11] != "" {
		switch GPRMCFix(r.Fields[11]) {
		case GPRMCFixAutonomous, GPRMCFixDifferential, GPRMCFixEstimated,
			GPRMCFixNotValid, GPRMCFixSimulator:

			g.FixType = GPRMCFix(r.Fields[11])
		default:
			return fmt.Errorf("unknown fix type value: %s", r.Fields[11])
		}
	} else {
		g.FixType = GPRMCFixUnspecified
	}

	return nil
}
