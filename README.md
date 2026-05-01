# MGRS-Go

Golang library for WGS84 **MGRS** over the UTM latitude band **`[-80°, 84°)`**. Grid-zone lettering and Norway/Svalbard UTM zone widening follow NGA-oriented practice (e.g. MSP GEOTRANS / GeographicLib). **Polar UPS** MGRS is not implemented.

**Module:** `git.quad4.io/Go-Libs/MGRS-Go`

## Requirements

Go **1.26+**

## Usage

```go
import "git.quad4.io/Go-Libs/MGRS-Go/mgrs"

s, err := mgrs.Encode(52.658, 5.892, mgrs.DefaultDigitPairs)
payload, err := mgrs.EncodeBytes(52.658, 5.892, mgrs.DefaultDigitPairs) // avoids string alloc
dst, err := mgrs.AppendEncode(nil, lat, lon, mgrs.DefaultDigitPairs)       // buffers many refs
buf := make([]byte, mgrs.MaxEncodedLen)
n, err := mgrs.EncodeTo(buf, 52.658, 5.892, mgrs.DefaultDigitPairs) // zero-alloc path
payload = buf[:n]

pt, err := mgrs.Decode(s, false /* SW corner of cell */)
z, err := mgrs.LongitudeZone(52.658, 5.892)
```

Runnable demo:

```bash
go run ./example -lat 52.658 -lon 5.892
go run ./example -decode 19TDJ3858897366
go run ./example -bench 500000 -lat 52.658 -lon 5.892
```

## Verification

When the **PROJ** `proj` binary is on `PATH`, tests compare UTM easting/northing to `+proj=utm +datum=WGS84` within a sub-metre envelope. Published MGRS literals are also decoded against known lat/lon benchmarks.

## Testing

| Goal | Command |
|------|---------|
| Default | `go test ./...` |
| Data races | `go test ./... -race` |
| Bounded fuzz | `go test ./mgrs -fuzz=Fuzz -fuzztime=5s` |
| Bench + allocs | `go test ./mgrs -bench=. -benchmem` |

## License

Copyright **2026 Quad4** — [LICENSE](LICENSE) (0BSD).
