package nmea

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const gprmcStr = "$GPRMC,232158.000,A,1445.1076,N,02315.4367,W,0.27,232.04,190516,,,D*79"
const gprmc3339 = "2016-05-19T23:21:58Z"

func TestGPRMC_Parse(t *testing.T) {
	r, err := ParseRaw([]byte(gprmcStr))
	if err != nil {
		t.Fatal(err)
	}

	g := new(GPRMC)
	err = g.Parse(r)
	assert.Nil(t, err)

	assert.Equal(t, gprmc3339, g.Time.Format(time.RFC3339Nano), "timestamp")
	assert.Equal(t, TypeGPRMC, g.Type(), "type")
	assert.Equal(t, 14.751793333333334, float64(g.Latitude), "latitude")
	assert.Equal(t, -7.257278333333333, float64(g.Longitude), "longitude")
	assert.Equal(t, 0.27, g.Speed, "speed")
	assert.Equal(t, 232.04, g.TrueCourse, "true course")
	assert.Zero(t, g.Variation, "variation")
	assert.Equal(t, GPRMCFixDifferential, g.FixType)
}
