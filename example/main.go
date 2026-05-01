package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"git.quad4.io/Go-Libs/MGRS-Go/mgrs"
)

func main() {
	decodeFlag := flag.String("decode", "", "decode this MGRS string (southwest corner of cell)")
	centerFlag := flag.Bool("center", false, "with -decode: return centre of truncated cell instead of SW corner")
	benchIterations := flag.Uint64("bench", 0, "after other actions run this many EncodeTo iterations and print throughput (0 = skip)")
	latStr := flag.String("lat", "52.658", "latitude decimal degrees north (negative = south)")
	lonStr := flag.String("lon", "5.892", "longitude decimal degrees east")
	pairs := flag.Int("pairs", mgrs.DefaultDigitPairs, fmt.Sprintf("MGRS digit pairs (max %d)", mgrs.MaxDigitPairs))

	flag.Parse()

	if *pairs < 0 || *pairs > mgrs.MaxDigitPairs {
		die(errors.New("-pairs out of supported range"))
	}

	if flag.NArg() > 0 {
		die(errors.New("unexpected arguments"))
	}

	if dec := *decodeFlag; dec != "" {
		pt, err := mgrs.Decode(dec, *centerFlag)
		if err != nil {
			die(err)
		}
		fmt.Printf("%s %+v lat=%f lon=%f\n", dec, pt, pt.Lat, pt.Lon)
	}

	lat, err := strconv.ParseFloat(*latStr, 64)
	if err != nil {
		die(fmt.Errorf("-lat: %w", err))
	}
	lon, err := strconv.ParseFloat(*lonStr, 64)
	if err != nil {
		die(fmt.Errorf("-lon: %w", err))
	}

	ref, err := mgrs.Encode(lat, lon, *pairs)
	if err != nil {
		die(err)
	}
	fmt.Printf("encode %+f %+f pairs=%d -> %s\n", lat, lon, *pairs, ref)

	raw, err := mgrs.EncodeBytes(lat, lon, *pairs)
	if err != nil {
		die(err)
	}
	pt, err := mgrs.Decode(string(raw), false)
	if err != nil {
		die(err)
	}
	fmt.Printf("round-trip %+v delta lat=%+.6g delta lon=%+.6g\n", pt, pt.Lat-lat, pt.Lon-lon)

	if *benchIterations == 0 {
		return
	}
	var benchBuf [mgrs.MaxEncodedLen]byte
	const warmup = 1_000
	for i := 0; i < warmup; i++ {
		if _, err := mgrs.EncodeTo(benchBuf[:], lat, lon, *pairs); err != nil {
			die(err)
		}
	}
	start := time.Now()
	var n uint64
	for n = 0; n < *benchIterations; n++ {
		if _, err := mgrs.EncodeTo(benchBuf[:], lat, lon, *pairs); err != nil {
			die(err)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("%d EncodeTo calls in %s (avg %s/op)\n", n, elapsed, time.Duration(elapsed.Nanoseconds()/int64(n)))
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
