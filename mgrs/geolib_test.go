package mgrs

import (
	"bytes"
	"fmt"
	"math"
	"os/exec"
	"strings"
	"testing"
)

// geoConvertPrec maps our digit-pair count to GeoConvert -p.
// GeographicLib uses digits_per_coord = 5 + prec for MGRS.
func geoConvertPrec(digitPairs int) int {
	return digitPairs - 5
}

func TestGeoLibEncodeMatchesGeoConvert(t *testing.T) {
	path := requireTool(t, "GeoConvert")
	cases := []struct {
		name     string
		lat, lon float64
		pairs    int
	}{
		{"europe", 52.658, 5.892, 5},
		{"antarctic-ups", -85, 77.85, 5},
		{"arctic-ups", 85, 0, 5},
		{"norway-zone", 60, 4, 5},
		{"svalbard-zone", 75, 12, 5},
		{"blue-marble-point", 44.226933333333335, -69.76891944444445, 5},
		{"equator", 0, 0, 5},
		{"melbourne", -37.8136, 144.9631, 5},
		{"baghdad-1km", 33.33424, 44.40363, 2},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var in bytes.Buffer
			fmt.Fprintf(&in, "%.12f %.12f\n", tc.lat, tc.lon)
			cmd := exec.Command(path, "-m", fmt.Sprintf("-p%d", geoConvertPrec(tc.pairs)))
			cmd.Stdin = &in
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("GeoConvert: %v", err)
			}
			want := strings.ReplaceAll(strings.TrimSpace(string(out)), " ", "")
			got, err := Encode(tc.lat, tc.lon, tc.pairs)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("got %q want GeoConvert %q", got, want)
			}
		})
	}
}

func TestGeoLibDecodeCentreMatchesGeoConvert(t *testing.T) {
	path := requireTool(t, "GeoConvert")
	refs := []string{
		"31UFU9559138152",
		"BHP4301516908",
		"38SMB4488",
		"19TDJ3858797365",
	}
	for _, ref := range refs {
		ref := ref
		t.Run(ref, func(t *testing.T) {
			cmd := exec.Command(path, "-g", "-p", "9")
			cmd.Stdin = strings.NewReader(ref + "\n")
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("GeoConvert decode: %v out=%s", err, out)
			}
			fields := strings.Fields(string(out))
			if len(fields) < 2 {
				t.Fatalf("unexpected GeoConvert output %q", out)
			}
			var lat, lon float64
			if _, err := fmt.Sscanf(fields[0], "%f", &lat); err != nil {
				t.Fatal(err)
			}
			if _, err := fmt.Sscanf(fields[1], "%f", &lon); err != nil {
				t.Fatal(err)
			}
			pt, err := Decode(ref, true)
			if err != nil {
				t.Fatal(err)
			}
			const tolDeg = 2e-5
			if math.Abs(pt.Lat-lat) > tolDeg || math.Abs(pt.Lon-lon) > tolDeg {
				t.Fatalf("decode centre vs GeoConvert\nlib  %.9f %.9f\ngeolib %.9f %.9f",
					pt.Lat, pt.Lon, lat, lon)
			}
		})
	}
}
