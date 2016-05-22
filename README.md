# nmea

[![go-doc](https://godoc.org/github.com/mastercactapus/nmea?status.svg)](https://godoc.org/github.com/mastercactapus/nmea)

A NMEA parsing library for go

See [gps_exporter](https://github.com/mastercactapus/gps_exporter) for example usage.

## Installation

Get the latest version with:

```bash
go get -u github.com/mastercactapus/nmea`
```

## Supported

Currently parsing and serialization are supported for the following sentence types:

- [GPRMC](https://godoc.org/github.com/mastercactapus/nmea#GPRMC)
- [GPGSA](https://godoc.org/github.com/mastercactapus/nmea#GPGSA)
- [GPGGA](https://godoc.org/github.com/mastercactapus/nmea#GPRMC)

## Example Usage

An example of parsing the timestamp from a GPRMC sentence:

```go

import "github.com/mastercactapus/nmea"

res, err := nmea.Parse([]byte("$GPRMC,232158.000,A,1445.1076,N,02315.4367,W,0.27,232.04,190516,,,D*79"))
if err != nil {
    panic(err)
}
if res.Type() != nmea.TypeGPRMC {
    panic("bad type")
}

fmt.Println(res.(*nmea.GPRMC).Time.String())

// Output: 2016-05-19 23:21:58 +0000 UTC
```



