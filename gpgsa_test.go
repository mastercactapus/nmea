package nmea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const gpgsaStr = "$GPGSA,A,3,03,06,19,24,12,28,01,17,,,,,1.39,1.10,0.84*00"

func TestGPGSA_Parse(t *testing.T) {
	r, err := ParseRaw([]byte(gpgsaStr))
	if err != nil {
		t.Fatal(err)
	}

	g := new(GPGSA)
	err = g.Parse(r)
	assert.Nil(t, err)

	assert.True(t, g.AutoSelection, "auto selection")
	assert.Equal(t, g.FixType, GPGSAFix3D, "fix type")
	assert.Equal(t, g.PDOP, 1.39, "PDOP")
	assert.Equal(t, g.HDOP, 1.1, "HDOP")
	assert.Equal(t, g.VDOP, 0.84, "VDOP")
	assert.EqualValues(t, g.Satellites, []string{"03", "06", "19", "24", "12", "28", "01", "17"})
}

func TestGPGSA_String(t *testing.T) {
	str := GPGSA{
		AutoSelection: false,
		FixType:       GPGSAFix2D,
		PDOP:          1.5,
		HDOP:          2.8,
		VDOP:          6.2,
		Satellites:    []string{"01", "02", "09"},
	}.String()

	assert.Equal(t, "$GPGSA,M,2,01,02,09,,,,,,,,,,1.5,2.8,6.2*3F", str)
}
