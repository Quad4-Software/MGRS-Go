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
n, err := mgrs.EncodeTo(buf, 52.658, 5.892, mgrs.DefaultDigitPairs)
refBytes := buf[:n] // EncodeTo: no new heap memory per call if buf is reused

pt, err := mgrs.Decode(s, false /* SW corner of cell */)
z, err := mgrs.LongitudeZone(52.658, 5.892)
```

Runnable demo:

```bash
go run ./example -lat 52.658 -lon 5.892
go run ./example -decode 19TDJ3858897366
go run ./example -bench 500000 -lat 52.658 -lon 5.892
```

## Performance and allocations

**Zero alloc** (in Go benchmark terms) means the hot path does **not allocate new heap memory** on each call: no fresh slice backing store for the MGRS bytes, no string header for the reference. The work reuses memory you already own (a fixed `[]byte` or a slice you reset and grow once). That cuts garbage-collection work and tends to steady latency in tight loops.

| API | Typical use | Heap allocs per encode (when measured with `-benchmem`) |
|-----|----------------|--------------------------------------------------------|
| `EncodeTo` | Caller-owned buffer, same buffer reused | **0** |
| `AppendEncode` | Append into a slice with enough capacity, reset `len` each time | **0** |
| `EncodeBytes` | Returns a new `[]byte` each call | **1** (the output buffer) |
| `Encode` | Returns a `string` | **1** (buffer + string; still very small) |

Decode is allocation-light on typical inputs; measure with `go test ./mgrs -bench=BenchmarkDecode_knownCatalogMgrs -benchmem`.

Rough ballpark on a fast desktop CPU (see `go test ./mgrs -bench=. -benchmem`): encode is on the order of **~100 ns** for default precision; **`EncodeTo`** / reused **`AppendEncode`** are a bit faster than **`EncodeBytes`** because they skip that per-call allocation.

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
