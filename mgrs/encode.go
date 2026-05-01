package mgrs

import "math"

// MaxDigitPairs is the finest supported easting+northing pair count encoded in
// the numerical MGRS tail (mirrors GeographicLib permitting sub-metre granularity).
const MaxDigitPairs = 11

const latBandAngeps = 1e-13

func encodeUTM(zone int, latDeg float64, north bool, easting, northing float64, digitPairs int) ([]byte, error) {
	return encodeUTMInto(nil, zone, latDeg, north, easting, northing, digitPairs)
}

func encodeUTMInto(dst []byte, zone int, latDeg float64, north bool, easting, northing float64, digitPairs int) ([]byte, error) {
	if digitPairs < 0 || digitPairs > MaxDigitPairs {
		return nil, ErrInvalidDigitPairs
	}
	if zone < 1 || zone > 60 {
		return nil, ErrInvalidMGRS
	}

	iband := resolveBandIndex(latDeg, north)
	xh, yh := eastNorthTileIndices(easting, northing)
	colLetters := utmcols[(zone-1)%3]
	icol := xh - minUTMCols
	if icol < 0 || icol >= len(colLetters) {
		return nil, ErrInvalidGridSquareBand
	}

	rowPeriodic := posMod(yh, utmerowPeriod)
	iTruth := utmRow(iband, icol, rowPeriodic)
	exp := yTileExpectation(north, yh)
	if iTruth != exp || iTruth == maxUTMsRowTiles {
		return nil, ErrLatitudeCoordsClash
	}

	base := len(dst)
	n := 5 + digitPairs*2
	need := base + n
	if cap(dst) < need {
		grow := make([]byte, base, need)
		copy(grow, dst)
		dst = grow
	}
	dst = dst[:need]
	out := dst[base:need]
	out[0] = digits[zone/mgrsdigitBase]
	out[1] = digits[zone%mgrsdigitBase]
	out[2] = latBand[10+iband]
	out[3] = colLetters[icol]
	rowPos := posMod(yh+((zone-1)&1)*utmevenRowShift, utmerowPeriod)
	out[4] = utmrows[rowPos]

	if digitPairs == 0 {
		return dst[:base+5], nil
	}

	ix := int64(math.Floor(easting * coordMultiply))
	iy := int64(math.Floor(northing * coordMultiply))
	ix -= tileScaleMicro * int64(xh)
	iy -= tileScaleMicro * int64(yh)

	divIdx := MaxDigitPairs - digitPairs
	if divIdx < 0 || divIdx >= len(pow10Int64) {
		divIdx = 0
	}
	div := pow10Int64[divIdx]
	if div == 0 {
		div = 1
	}
	ix /= div
	iy /= div

	off := 5
	for i := digitPairs - 1; i >= 0; i-- {
		out[off+i] = digits[ix%mgrsdigitBase]
		ix /= mgrsdigitBase
		out[off+digitPairs+i] = digits[iy%mgrsdigitBase]
		iy /= mgrsdigitBase
	}

	return dst, nil
}

func eastNorthTileIndices(easting, northing float64) (xh int, yh int) {
	x := math.Floor(easting * coordMultiply)
	y := math.Floor(northing * coordMultiply)
	return int(int64(x) / tileScaleMicro), int(int64(y) / tileScaleMicro)
}

func yTileExpectation(north bool, yh int) int {
	if north {
		return yh
	}
	return yh - maxUTMsRowTiles
}

func resolveBandIndex(latDeg float64, north bool) int {
	switch {
	case latDeg > -latBandAngeps && latDeg < latBandAngeps:
		if north {
			return 0
		}
		return -1
	default:
		return latitudeBand(latDeg)
	}
}
