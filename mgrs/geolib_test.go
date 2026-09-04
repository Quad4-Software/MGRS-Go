package mgrs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestGeoLibEncodeMatchesGeoConvert(t *testing.T) {
	path, err := exec.LookPath("GeoConvert")
	if err != nil {
		t.Skip("GeoConvert not on PATH")
	}
	cases := []struct {
		lat, lon float64
		prec     int
	}{
		{52.658, 5.892, 5},
		{-85, 77.85, 5},
		{85, 0, 5},
		{60, 4, 5},
		{44.226933333333335, -69.76891944444445, 5},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("%.4f_%.4f", tc.lat, tc.lon), func(t *testing.T) {
			var in bytes.Buffer
			fmt.Fprintf(&in, "%.12f %.12f\n", tc.lat, tc.lon)
			cmd := exec.Command(path, "-m", fmt.Sprintf("-p%d", tc.prec))
			cmd.Stdin = &in
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("GeoConvert: %v", err)
			}
			want := strings.TrimSpace(string(out))
			want = strings.ReplaceAll(want, " ", "")
			got, err := Encode(tc.lat, tc.lon, tc.prec)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("got %q want GeoConvert %q", got, want)
			}
		})
	}
}
