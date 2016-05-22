package nmea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const epsilon = 0.000001

func TestCoordFromDDM(t *testing.T) {
	assert.InEpsilon(t, 12.065, float64(CoordFromDDM(12, 3.9, CoordDirectionNorth)), epsilon)
	assert.InEpsilon(t, -12.065, float64(CoordFromDDM(12, 3.9, CoordDirectionSouth)), epsilon)
	assert.InEpsilon(t, 12.065, float64(CoordFromDDM(12, 3.9, CoordDirectionEast)), epsilon)
	assert.InEpsilon(t, -12.065, float64(CoordFromDDM(12, 3.9, CoordDirectionWest)), epsilon)
}
func TestCoordFromDD(t *testing.T) {
	assert.InEpsilon(t, 10.1, float64(CoordFromDD(10.1, CoordDirectionNorth)), epsilon)
	assert.InEpsilon(t, -10.1, float64(CoordFromDD(10.1, CoordDirectionSouth)), epsilon)
	assert.InEpsilon(t, 10.1, float64(CoordFromDD(10.1, CoordDirectionEast)), epsilon)
	assert.InEpsilon(t, -10.1, float64(CoordFromDD(10.1, CoordDirectionWest)), epsilon)
}
func TestCoordFromDMS(t *testing.T) {
	assert.InEpsilon(t, 12.065, float64(CoordFromDMS(12, 3, 54, CoordDirectionNorth)), epsilon)
	assert.InEpsilon(t, -12.065, float64(CoordFromDMS(12, 3, 54, CoordDirectionSouth)), epsilon)
	assert.InEpsilon(t, 12.065, float64(CoordFromDMS(12, 3, 54, CoordDirectionEast)), epsilon)
	assert.InEpsilon(t, -12.065, float64(CoordFromDMS(12, 3, 54, CoordDirectionWest)), epsilon)
}
func TestCoord_Direction(t *testing.T) {
	assert.Equal(t, CoordDirectionNorth, Coord(1).Direction())
	assert.Equal(t, CoordDirectionSouth, Coord(-1).Direction())
	assert.Equal(t, CoordDirectionEast, Coord(1).Direction())
	assert.Equal(t, CoordDirectionWest, Coord(-1).Direction())
}
func TestCoord_DD(t *testing.T) {
	deg, dir := Coord(12.065).DD()
	assert.InEpsilon(t, 12.065, deg, epsilon)
	assert.Equal(t, CoordDirectionNorth, dir)
	assert.Equal(t, CoordDirectionEast, dir)
	deg, dir = Coord(-12.065).DD()
	assert.InEpsilon(t, 12.065, deg, epsilon)
	assert.Equal(t, CoordDirectionSouth, dir)
	assert.Equal(t, CoordDirectionWest, dir)
}
func TestCoord_DDM(t *testing.T) {
	deg, min, dir := Coord(12.065).DDM()
	assert.InEpsilon(t, 12, deg, epsilon)
	assert.InEpsilon(t, 3.9, min, epsilon)
	assert.Equal(t, CoordDirectionNorth, dir)
	assert.Equal(t, CoordDirectionEast, dir)
	deg, min, dir = Coord(-12.065).DDM()
	assert.InEpsilon(t, 12, deg, epsilon)
	assert.InEpsilon(t, 3.9, min, epsilon)
	assert.Equal(t, CoordDirectionSouth, dir)
	assert.Equal(t, CoordDirectionWest, dir)
}
func TestCoord_DMS(t *testing.T) {
	deg, min, sec, dir := Coord(12.065).DMS()
	assert.InEpsilon(t, 12, deg, epsilon)
	assert.InEpsilon(t, 3, min, epsilon)
	assert.InEpsilon(t, 54, sec, epsilon)
	assert.Equal(t, CoordDirectionNorth, dir)
	assert.Equal(t, CoordDirectionEast, dir)
	deg, min, sec, dir = Coord(-12.065).DMS()
	assert.InEpsilon(t, 12, deg, epsilon)
	assert.InEpsilon(t, 3, min, epsilon)
	assert.InEpsilon(t, 54, sec, epsilon)
	assert.Equal(t, CoordDirectionSouth, dir)
	assert.Equal(t, CoordDirectionWest, dir)
}
func TestCoord_Parse(t *testing.T) {
	c, err := ParseCoord("1203.9", CoordDirectionNorth)
	assert.Nil(t, err)
	assert.InEpsilon(t, 12.065, float64(c), epsilon)
	c, err = ParseCoord("1203.9", CoordDirectionSouth)
	assert.Nil(t, err)
	assert.InEpsilon(t, -12.065, float64(c), epsilon)
}
func TestCoord_String(t *testing.T) {
	assert.Equal(t, "1203.9", Coord(12.065).String())
	assert.Equal(t, "1203.9", Coord(-12.065).String())
}

func TestCoordDirection_LatString(t *testing.T) {
	assert.Equal(t, "N", CoordDirectionNorth.LatString())
	assert.Equal(t, "S", CoordDirectionSouth.LatString())
}
func TestCoordDirection_LongString(t *testing.T) {
	assert.Equal(t, "E", CoordDirectionEast.LongString())
	assert.Equal(t, "W", CoordDirectionWest.LongString())
}
