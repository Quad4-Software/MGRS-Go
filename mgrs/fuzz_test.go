package mgrs

import (
	"bytes"
	"math"
	"testing"
)

func fuzzSanitizeAscii(b []byte) []byte {
	const maxChars = 80
	if len(b) > maxChars {
		b = b[:maxChars]
	}
	dst := make([]byte, 0, len(b))
	for _, x := range b {
		switch {
		case x == ' ' || (x >= '0' && x <= '9') || (x >= 'A' && x <= 'Z') || (x >= 'a' && x <= 'z'):
			dst = append(dst, x)
		default:
			dst = append(dst, 'X')
		}
	}
	return dst
}

func FuzzDecodeRandomAsciiNoPanic(f *testing.F) {
	f.Add([]byte("19TDJ3858897366"))
	f.Add([]byte(`03u xd 98281 74621`))
	f.Add([]byte("garbage-string"))
	f.Fuzz(func(t *testing.T, payload []byte) {
		buf := bytes.ToUpper(fuzzSanitizeAscii(payload))
		_, _ = Decode(string(bytes.TrimSpace(buf)), false)
	})
}

func clampLatForFuzz(lat float64) float64 {
	switch {
	case lat < -79.999:
		lat = -79.999
	case lat > 83.999:
		lat = 83.999
	default:
	}
	return lat
}

func wrapLonDeg180(lon float64) float64 {
	l := math.Mod(lon+180, 360)
	if l < 0 {
		l += 360
	}
	return l - 180
}

func FuzzEncodeDecodeIEEE(f *testing.F) {
	seeds := [][2]uint64{
		{math.Float64bits(52.658), math.Float64bits(5.892)},
		{math.Float64bits(-33.9243), math.Float64bits(18.4241)},
		{math.Float64bits(71.1699), math.Float64bits(25.7836)},
		{math.Float64bits(0), math.Float64bits(174.764)},
	}
	for _, s := range seeds {
		f.Add(s[0], s[1])
	}
	const tolDeg = 1e-2
	const pairs = DefaultDigitPairs

	f.Fuzz(func(t *testing.T, latBits, lonBits uint64) {
		lat := math.Float64frombits(latBits)
		lon := math.Float64frombits(lonBits)
		if math.IsNaN(lat) || math.IsNaN(lon) || math.IsInf(lat, 0) || math.IsInf(lon, 0) {
			t.Skip()
		}
		lat = clampLatForFuzz(lat)
		lon = wrapLonDeg180(lon)

		s, encErr := Encode(lat, lon, pairs)
		if encErr != nil {
			t.Fatalf("unexpected encode failure lat=%g lon=%g err=%v", lat, lon, encErr)
		}
		pt, decErr := Decode(s, false)
		if decErr != nil {
			t.Fatal(decErr)
		}
		cosClamp := math.Max(math.Abs(math.Cos(lat*degRad)), 1e-3)
		lonTol := tolDeg / cosClamp
		if math.Abs(pt.Lat-lat) > tolDeg || math.Abs(pt.Lon-lon) > lonTol {
			t.Fatalf(
				"roundtrip drift MGRS %q in lat %+e lon %+e from lat=%g lon=%g",
				s, pt.Lat-lat, pt.Lon-lon, lat, lon,
			)
		}
	})
}
