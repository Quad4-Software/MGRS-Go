# MGRS-Go

WGS84 MGRS encode and decode for the full globe: UTM in `[-80°, 84°)` and polar UPS outside that belt.

Grid-zone lettering, Norway/Svalbard zone widening, UPS bands `A/B/Y/Z`, and upper-bound ~nanometre nudges follow GeographicLib / MSP GEOTRANS practice.

**Module:** `github.com/Quad4-Software/MGRS-Go`

## Requirements

- Go 1.26.5 or later

## Install

```bash
go get github.com/Quad4-Software/MGRS-Go@latest
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/Quad4-Software/MGRS-Go/mgrs"
)

func main() {
	ref, err := mgrs.Encode(52.658, 5.892, mgrs.DefaultDigitPairs)
	if err != nil {
		panic(err)
	}
	fmt.Println(ref)

	pt, err := mgrs.Decode(ref, false) // southwest corner of the cell
	if err != nil {
		panic(err)
	}
	fmt.Println(pt.Lat, pt.Lon)

	// Polar UPS
	polar, err := mgrs.Encode(-85, 77.85, mgrs.DefaultDigitPairs)
	if err != nil {
		panic(err)
	}
	fmt.Println(polar) // BHP4301516908

	parts, err := mgrs.DecodeParts(ref, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(parts.Zone, string([]byte{parts.Band, parts.Col, parts.Row}), parts.Easting, parts.Northing)

	pretty, err := mgrs.FormatSpaced(ref)
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty)
}
```

### Grid metres (UTM / UPS)

```go
g, err := mgrs.LatLonToGrid(52.658, 5.892)
ref, err := mgrs.EncodeGrid(g, mgrs.DefaultDigitPairs)

g2 := mgrs.Grid{Zone: 38, North: true, Easting: 444000, Northing: 3688000}
ref2, err := mgrs.EncodeGrid(g2, 2) // 38SMB4488
```

`LongitudeZone` remains UTM-only (zones 1–60). Use `LatLonToGrid` for polar points (`ZoneUPS`).

### Allocation-conscious encode

| API | Typical use | Heap allocs per call |
|-----|-------------|----------------------|
| `EncodeTo` / `EncodeGridTo` | Caller-owned buffer, reused | 0 |
| `AppendEncode` / `AppendEncodeGrid` | Append into a slice with capacity, reset `len` each time | 0 |
| `EncodeBytes` / `EncodeGridBytes` | New `[]byte` each call | 1 |
| `Encode` / `EncodeGrid` | New `string` each call | 1 |

```go
buf := make([]byte, mgrs.MaxEncodedLen)
n, err := mgrs.EncodeTo(buf, 52.658, 5.892, mgrs.DefaultDigitPairs)
refBytes := buf[:n]

dst := make([]byte, 0, mgrs.MaxEncodedLen)
dst, err = mgrs.AppendEncode(dst[:0], 52.658, 5.892, mgrs.DefaultDigitPairs)
```

### Example CLI

```bash
go run ./example -lat 52.658 -lon 5.892
go run ./example -lat -85 -lon 77.85
go run ./example -decode 19TDJ3858897366
go run ./example -decode BHP4301516908 -pretty
go run ./example -bench 500000 -lat 52.658 -lon 5.892
```

## Verification

CI installs proven external tools and **requires** them (`MGRS_REQUIRE_ORACLES=1`):

| Oracle | What it checks |
|--------|----------------|
| Checked-in goldens (`mgrs/testdata/golden_mgrs.jsonl`) | Exact MGRS strings from GeographicLib docs, Blue Marble, Norway/Svalbard, UPS |
| PROJ `proj` | UTM easting/northing within 0.5 m |
| PROJ `cs2cs` | UPS stereographic metres within 1 m |
| GeographicLib `GeoConvert` | Exact MGRS encode strings and centre-cell decode lat/lon |

Locally, those CLI tests skip if the tools are missing unless you export `MGRS_REQUIRE_ORACLES=1`.

```bash
# Debian/Ubuntu
sudo apt-get install -y proj-bin geographiclib-tools
export MGRS_REQUIRE_ORACLES=1
go test ./mgrs -count=1 -run 'Proj|GeoLib|Golden'
```

Default-precision encode is about 90–120 ns/op on a fast desktop CPU. Measure with:

```bash
go test ./mgrs -bench=. -benchmem
```

## Testing

| Goal | Command |
|------|---------|
| Default | `go test ./...` |
| Goldens | `go test ./mgrs -run Golden -v` |
| PROJ metres | `go test ./mgrs -run Proj -v` |
| GeoConvert | `go test ./mgrs -run GeoLib -v` |
| Required oracles | `MGRS_REQUIRE_ORACLES=1 go test ./mgrs -run 'Proj|GeoLib|Golden'` |
| Race detector | `go test ./... -race` |
| Bounded fuzz | `go test ./mgrs -fuzz=FuzzEncodeDecodeIEEE -fuzztime=5s` |
| Bench + allocs | `go test ./mgrs -bench=. -benchmem` |

## License

Copyright 2026 Quad4. See [LICENSE](LICENSE) (0BSD).
