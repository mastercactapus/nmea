package nmea

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const rawSentence = `$GPVTG,230.17,T,,M,0.38,N,0.70,K,D*33`

var rawFields = []string{"230.17", "T", "", "M", "0.38", "N", "0.70", "K", "D"}

func TestParseRaw(t *testing.T) {
	r, err := ParseRaw([]byte(rawSentence))
	assert.Nil(t, err)
	assert.Equal(t, "GPVTG", r.TypeName)
	assert.EqualValues(t, rawFields, r.Fields)
}

func TestRaw_String(t *testing.T) {
	r := &Raw{
		TypeName: "GPVTG",
		Fields:   rawFields,
	}
	assert.Equal(t, rawSentence, r.String())
}

func ExampleParse_string() {
	res, err := Parse([]byte("$GPRMC,232158.000,A,1445.1076,N,02315.4367,W,0.27,232.04,190516,,,D*79"))
	if err != nil {
		panic(err)
	}
	if res.Type() != TypeGPRMC {
		panic("bad type")
	}

	fmt.Println(res.(*GPRMC).Time.String())
	// Output: 2016-05-19 23:21:58 +0000 UTC
}
