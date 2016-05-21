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
	// Time is the time/date of the fix
	Time time.Time

	// Active will be true if the unit reports the fix as valid
	Active bool

	Latitude  Coord
	Longitude Coord

	// Speed in knots
	Speed float64

	// TrueCourse -- track made good in degrees True
	TrueCourse float64

	// Var is the magnetic variation
	Variation Coord

	// FixType is the type of fix the receiver has
	FixType GPRMCFix
}

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
			strconv.FormatFloat(g.Speed, 'f', 5, 64),
			strconv.FormatFloat(g.TrueCourse, 'f', 5, 64),
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
			return fmt.Errorf("parse timestamp:", err)
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

	var dir CoordDirection
	if r.Fields[2] != "" {
		switch r.Fields[3] {
		case "N":
			dir = CoordDirectionNorth
		case "S":
			dir = CoordDirectionSouth
		default:
			return fmt.Errorf("invalid or missing direction for latitude")
		}
		g.Latitude, err = ParseCoord(r.Fields[2], dir)
		if err != nil {
			return fmt.Errorf("parse latitude: %s", err)
		}
	} else if r.Fields[3] != "" {
		return fmt.Errorf("got direction for latitude, but no latitude value")
	} else {
		g.Latitude = 0
	}

	if r.Fields[4] != "" {
		switch r.Fields[5] {
		case "E":
			dir = CoordDirectionEast
		case "W":
			dir = CoordDirectionWest
		default:
			return fmt.Errorf("invalid or missing direction for longitude")
		}
		g.Longitude, err = ParseCoord(r.Fields[4], dir)
		if err != nil {
			return fmt.Errorf("parse longitude: %s", err)
		}
	} else if r.Fields[5] != "" {
		return fmt.Errorf("got direction for longitude, but no longitude value")
	} else {
		g.Longitude = 0
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

	if r.Fields[8] != "" {
		t, err := time.ParseInLocation(dateFormat, r.Fields[8], time.UTC)
		if err != nil {
			return fmt.Errorf("parse date: %s", err)
		}
		// add date & time
		g.Time = t.Add(time.Hour*time.Duration(g.Time.Hour()) + time.Minute*time.Duration(g.Time.Minute()) + time.Second*time.Duration(g.Time.Second()) + time.Duration(g.Time.Nanosecond()))
	}

	if r.Fields[9] != "" {
		switch r.Fields[10] {
		case "E":
			dir = CoordDirectionEast
		case "W":
			dir = CoordDirectionWest
		default:
			return fmt.Errorf("invalid or missing direction for variation")
		}
		g.Variation, err = ParseCoord(r.Fields[9], dir)
		if err != nil {
			return fmt.Errorf("parse variation: %s", err)
		}
	} else if r.Fields[10] != "" {
		return fmt.Errorf("got direction for variation, but no variation value")
	} else {
		g.Variation = 0
	}

	if len(r.Fields) >= 12 && r.Fields[11] != "" {
		switch GPRMCFix(r.Fields[11]) {
		case GPRMCFixAutonomous, GPRMCFixDifferential:
			if !g.Active {
				return fmt.Errorf("fix type mismatched with status")
			}
			g.FixType = GPRMCFix(r.Fields[11][0])
		case GPRMCFixEstimated, GPRMCFixNotValid, GPRMCFixSimulator:
			if !g.Active {
				return fmt.Errorf("fix type mismatched with status")
			}
			g.FixType = GPRMCFix(r.Fields[11][0])
		default:
			return fmt.Errorf("unknown fix type value: %s", r.Fields[11])
		}
	} else {
		g.FixType = GPRMCFixUnspecified
	}

	return nil
}
