package nmea

import (
	"fmt"
	"strconv"
	"time"
)

// GPGGAFix is the fix type/quality for a GPGGA sentence
type GPGGAFix string

// Fix types for GPGGA
const (
	GPGGAFixInvalid    GPGGAFix = "0"
	GPGGAFixGPS        GPGGAFix = "1" // (SPS)
	GPGGAFixDGPS       GPGGAFix = "2"
	GPGGAFixPPS        GPGGAFix = "3"
	GPGGAFixRTK        GPGGAFix = "4" // Real Time Kinematic
	GPGGAFixFRTK       GPGGAFix = "5" // Float RTK
	GPGGAFixEstimated  GPGGAFix = "6" // dead reckoning
	GPGGAFixManual     GPGGAFix = "7" // Manual input mode
	GPGGAFixSimulation GPGGAFix = "8"
)

// GPGGA contains essential fix data including 3D location and accuracy data
type GPGGA struct {
	Time        time.Time // Time the fix was taken
	Latitude    Coord
	Longitude   Coord
	FixType     GPGGAFix      // type/quality of the fix
	Satellites  int           // number of satellites used
	HDOP        float64       // horizontal dilution of precision
	Altitude    float64       // altitude
	GeoIDHeight float64       // geoid height. if this is missing; altitude is suspect
	DGPSUpdate  time.Duration // time since last DGPS update
	DGPSID      string        // DGPS station ID
}

// Type will return TypeGPGGA to fulfill the Sentence interface
func (g GPGGA) Type() Type {
	return TypeGPGGA
}

// String will provide a NMEA formatted string. Date information from the Time field is ignored.
func (g GPGGA) String() string {
	return Raw{
		TypeName: string(TypeGPGGA),
		Fields: []string{
			g.Time.Format(timeFormat),
			g.Latitude.String(),
			g.Latitude.Direction().LatString(),
			g.Longitude.String(),
			g.Longitude.Direction().LongString(),
			string(g.FixType),
			strconv.Itoa(g.Satellites),
			strconv.FormatFloat(g.HDOP, 'f', -1, 64),
			strconv.FormatFloat(g.Altitude, 'f', -1, 64),
			"M",
			strconv.FormatFloat(g.GeoIDHeight, 'f', -1, 64),
			"M",
			strconv.Itoa(int(g.DGPSUpdate.Seconds())),
			g.DGPSID,
		},
	}.String()
}

// Parse will parse GPGGA data from a raw sentence struct
func (g *GPGGA) Parse(r *Raw) error {
	if r.TypeName != string(TypeGPGGA) {
		return fmt.Errorf("wrong type for GPGGA '%s'", r.TypeName)
	}
	if r.Fields == nil || len(r.Fields) < 14 {
		return fmt.Errorf("not enough fields, need at least 14")
	}

	var err error
	if r.Fields[0] != "" {
		g.Time, err = time.ParseInLocation(timeFormat, r.Fields[0], time.UTC)
		if err != nil {
			return fmt.Errorf("parse time: %s", err)
		}
	} else {
		g.Time = time.Time{}
	}

	g.Latitude, err = parseFieldCoord(r.Fields[1], r.Fields[2], "latitude")
	if err != nil {
		return err
	}
	g.Longitude, err = parseFieldCoord(r.Fields[3], r.Fields[4], "longitude")
	if err != nil {
		return err
	}

	switch GPGGAFix(r.Fields[5]) {
	case GPGGAFixInvalid, GPGGAFix(""):
		g.FixType = GPGGAFixInvalid
	case GPGGAFixGPS, GPGGAFixDGPS, GPGGAFixPPS, GPGGAFixRTK, GPGGAFixFRTK,
		GPGGAFixEstimated, GPGGAFixManual, GPGGAFixSimulation:

		g.FixType = GPGGAFix(r.Fields[5])
	default:
		return fmt.Errorf("invalid fix type: %s", r.Fields[5])
	}

	g.Satellites, err = parseFieldInt(r.Fields[6], "Satellites")
	if err != nil {
		return err
	}

	g.HDOP, err = parseFieldFloat(r.Fields[7], "HDOP")
	if err != nil {
		return err
	}

	g.Altitude, err = parseFieldFloat(r.Fields[8], "Altitude")
	if err != nil {
		return err
	}
	if r.Fields[8] != "" && r.Fields[9] != "" && r.Fields[9] != "M" {
		return fmt.Errorf("unknown unit for Altitude: %s", r.Fields[9])
	}

	g.GeoIDHeight, err = parseFieldFloat(r.Fields[10], "GeoIDHeight")
	if err != nil {
		return err
	}
	if r.Fields[10] != "" && r.Fields[11] != "" && r.Fields[11] != "M" {
		return fmt.Errorf("unknown unit for GeoIDHeight: %s", r.Fields[11])
	}

	if r.Fields[12] != "" {
		g.DGPSUpdate, err = time.ParseDuration(r.Fields[12] + "s")
		if err != nil {
			return fmt.Errorf("parse DGPSUpdate: %s", err)
		}
	} else {
		g.DGPSUpdate = time.Duration(0)
	}

	g.DGPSID = r.Fields[13]
	return nil
}
