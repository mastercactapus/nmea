package nmea

import (
	"fmt"
	"math"
	"strconv"
)

// CoordDirection specifies a hemesphere
type CoordDirection bool

const (
	CoordDirectionNorth CoordDirection = true
	CoordDirectionEast  CoordDirection = true
	CoordDirectionSouth CoordDirection = false
	CoordDirectionWest  CoordDirection = false
)

// LatString will return the direction string for latitude (N or S)
func (c CoordDirection) LatString() string {
	if c == CoordDirectionNorth {
		return "N"
	}
	return "S"
}

// LongString will return the direction string for longitude (E or W)
func (c CoordDirection) LongString() string {
	if c == CoordDirectionEast {
		return "E"
	}
	return "W"
}

// a Coord represents a geographic coordinate
type Coord float64

// CoordFromDD will return a Coord from Decimal Degrees and direction
func CoordFromDD(deg float64, dir CoordDirection) Coord {
	if !bool(dir) {
		deg = -deg
	}
	return Coord(deg)
}

// CoordFromDMS will return a Coord from Degrees, Minutes, Seconds and direction
func CoordFromDMS(deg, min, sec float64, dir CoordDirection) Coord {
	deg += min/60 + sec/3600
	if !bool(dir) {
		deg = -deg
	}
	return Coord(deg)
}

// CoordFromDDM will return a Coord from Degrees, Decimal Minutes, and direction
func CoordFromDDM(deg, min float64, dir CoordDirection) Coord {
	deg += min / 60
	if !bool(dir) {
		deg = -deg
	}
	return Coord(deg)
}

// DD will return Decimal Degrees and direction
func (c Coord) DD() (deg float64, dir CoordDirection) {
	if c >= 0 {
		dir = true
	}
	return float64(c), dir
}

// DMS will return Degrees, Minutes, Seconds and direction
func (c Coord) DMS() (deg, min, sec float64, dir CoordDirection) {
	val := float64(c)
	if val < 0 {
		val = -val
	} else {
		dir = true
	}

	deg = math.Floor(val)
	min = math.Floor(60 * (val - deg))
	sec = 3600 * ((val - deg) - min/60)
	return deg, min, sec, dir
}

// DDM will return Degrees, Decimal Minutes, and direction
func (c Coord) DDM() (deg, min float64, dir CoordDirection) {
	val := float64(c)
	if val < 0 {
		val = -val
	} else {
		dir = true
	}
	deg = math.Floor(val)
	min = 60 * (val - deg)
	return deg, min, dir
}

func ParseCoord(c string, dir CoordDirection) (Coord, error) {
	if len(c) < 3 {
		deg, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return 0, err
		}
		return CoordFromDD(deg, dir), nil
	}
	deg, err := strconv.ParseFloat(c[:2], 64)
	if err != nil {
		return 0, err
	}
	min, err := strconv.ParseFloat(c[2:], 64)
	if err != nil {
		return 0, err
	}
	return CoordFromDDM(deg, min, dir), nil
}

// StringLat will return the coordinate as a NMEA-formatted string
func (c Coord) String() string {
	deg, min, _ := c.DDM()
	return fmt.Sprintf("%02.0f%02f", deg, min)
}

// Direction will return the direction of the coordinate
func (c Coord) Direction() CoordDirection {
	return c >= 0
}
