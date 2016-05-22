package nmea

import (
	"fmt"
	"strconv"
)

func parseFieldInt(val, name string) (int, error) {
	if val == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %s", name, err)
	}
	return i, nil
}

func parseFieldFloat(val, name string) (float64, error) {
	if val == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %s", name, err)
	}
	return f, nil
}

func parseFieldCoord(val, dirStr, typeName string) (Coord, error) {
	var dir CoordDirection
	if val == "" && dirStr != "" {
		return 0, fmt.Errorf("got direction for %s, but no %s value", typeName, typeName)
	} else if val == "" {
		return 0, nil
	}

	if typeName == "latitude" {
		switch dirStr {
		case "N":
			dir = CoordDirectionNorth
		case "S":
			dir = CoordDirectionSouth
		default:
			return 0, fmt.Errorf("invalid or missing direction for %s", typeName)
		}
	} else {
		switch dirStr {
		case "E":
			dir = CoordDirectionEast
		case "W":
			dir = CoordDirectionWest
		default:
			return 0, fmt.Errorf("invalid or missing direction for %s", typeName)
		}
	}

	c, err := ParseCoord(val, dir)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %s", typeName, err)
	}
	return c, nil
}
