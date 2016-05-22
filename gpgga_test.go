package nmea

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const gpggaStr = "$GPGGA,232200.000,1445.1076,N,02315.4370,W,2,08,1.10,310.5,M,-31.9,M,0000,0000*54"

func TestGPGGA_Parse(t *testing.T) {
	r, err := ParseRaw([]byte(gpggaStr))
	if err != nil {
		t.Fatal(err)
	}

	g := new(GPGGA)
	err = g.Parse(r)
	assert.Nil(t, err)

	assert.Equal(t, "232200", g.Time.Format(timeFormat), "timestamp")
	assert.Equal(t, TypeGPGGA, g.Type(), "type")
	assert.Equal(t, 14.751793333333334, float64(g.Latitude), "latitude")
	assert.Equal(t, -7.2572833333333335, float64(g.Longitude), "longitude")
	assert.Equal(t, GPGGAFixDGPS, g.FixType)
	assert.Equal(t, 8, g.Satellites, "satellites")
	assert.Equal(t, 1.1, g.HDOP, "HDOP")
	assert.Equal(t, 310.5, g.Altitude, "Altitude")
	assert.Equal(t, -31.9, g.GeoIDHeight, "GeoIDHeight")
	assert.Equal(t, time.Duration(0), g.DGPSUpdate, "DGPSUpdate")
	assert.Equal(t, "0000", g.DGPSID, "DGPSID")
}

func TestGPGGA_String(t *testing.T) {
	tm, err := time.ParseInLocation("1/2/06 15:04:05", "1/2/03 4:05:06", time.UTC)
	if err != nil {
		panic(err)
	}
	str := GPGGA{
		Time:        tm,
		Latitude:    Coord(12.065),
		Longitude:   Coord(-12.065),
		FixType:     GPGGAFixSimulation,
		Satellites:  4,
		HDOP:        2.3,
		Altitude:    10.4,
		GeoIDHeight: 12.3,
		DGPSUpdate:  time.Minute,
		DGPSID:      "bob",
	}.String()

	assert.Equal(t, "$GPGGA,040506,1203.9,N,1203.9,W,8,4,2.3,10.4,M,12.3,M,60,bob*07", str)
}
