package mgrs

func stringIndexAscii(s string, needle byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == needle {
			return i
		}
	}
	return -1
}

func byteIndexAscii(chars []byte, c byte) int {
	for i, v := range chars {
		if v == c {
			return i
		}
	}
	return -1
}

func parseDigit(b byte) (int, bool) {
	switch {
	case b >= '0' && b <= '9':
		return int(b - '0'), true
	default:
		return 0, false
	}
}

func parseZonePrefix(buf []byte) (zone int, used int, err error) {
	if len(buf) == 0 {
		return 0, 0, ErrInvalidMGRS
	}
	z0, ok := parseDigit(buf[0])
	if !ok || z0 == 0 {
		return 0, 0, ErrInvalidMGRS
	}
	zone = z0
	used = 1
	if used < len(buf) {
		if z1, ok := parseDigit(buf[used]); ok {
			next := zone*10 + z1
			if next > 60 {
				return 0, 0, ErrInvalidMGRS
			}
			zone = next
			used = 2
		}
	}
	if zone < 1 || zone > 60 {
		return 0, 0, ErrInvalidMGRS
	}
	return zone, used, nil
}

func compactUpperAsciiFromString(ref string) []byte {
	out := make([]byte, 0, len(ref))
	for i := 0; i < len(ref); i++ {
		ch := ref[i]
		switch ch {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			if ch >= 'a' && ch <= 'z' {
				ch -= 'a' - 'A'
			}
			out = append(out, ch)
		}
	}
	return out
}

func decodeNormalizedUpper(data []byte, centerCell bool) (Point, error) {
	zone, used, err := parseZonePrefix(data)
	if err != nil {
		return Point{}, err
	}
	idx := used
	if idx >= len(data) {
		return Point{}, ErrInvalidMGRS
	}

	bandLUT := byteIndexAscii(latBand, data[idx])
	if bandLUT < 0 {
		return Point{}, ErrInvalidMGRS
	}
	idx++
	northp := bandLUT >= 10

	if idx+2 > len(data) {
		return Point{}, ErrInvalidMGRS
	}

	col := stringIndexAscii(utmcols[(zone-1)%3], data[idx])
	if col < 0 {
		return Point{}, ErrInvalidMGRS
	}
	idx++

	rowPeriodic := byteIndexAscii(utmrows, data[idx])
	if rowPeriodic < 0 {
		return Point{}, ErrInvalidMGRS
	}
	idx++

	iRow := rowPeriodic
	if (zone-1)&1 != 0 {
		iRow = posMod(iRow+(utmerowPeriod-utmevenRowShift), utmerowPeriod)
	}

	iBand := bandLUT - 10
	rowTiles := utmRow(iBand, col, iRow)
	if rowTiles == maxUTMsRowTiles {
		return Point{}, ErrInvalidGridSquareBand
	}
	if !northp {
		rowTiles += 100
	}

	e100k := col + minUTMCols

	tail := data[idx:]
	if len(tail)%2 != 0 {
		return Point{}, ErrInvalidMGRS
	}
	pairs := len(tail) / 2
	if pairs > MaxDigitPairs {
		return Point{}, ErrInvalidMGRS
	}

	for _, ch := range tail {
		if ch < '0' || ch > '9' {
			return Point{}, ErrInvalidMGRS
		}
	}

	scale := float64(1)

	xQty := float64(e100k)
	yQty := float64(rowTiles)

	eastDigits := tail[:pairs]
	northDigits := tail[pairs:]

	for i := 0; i < pairs; i++ {
		ed, okE := parseDigit(eastDigits[i])
		nd, okN := parseDigit(northDigits[i])
		if !okE || !okN {
			return Point{}, ErrInvalidMGRS
		}
		scale *= float64(mgrsdigitBase)
		xQty = float64(mgrsdigitBase)*xQty + float64(ed)
		yQty = float64(mgrsdigitBase)*yQty + float64(nd)
	}

	if centerCell {
		scale *= 2
		xQty = 2*xQty + 1
		yQty = 2*yQty + 1
	}

	easting := mgrsTileSize * xQty / scale
	northing := mgrsTileSize * yQty / scale

	latDeg, lonDeg := inverseTM(easting, northing, zone, !northp)

	return Point{Lat: latDeg, Lon: lonDeg}, nil
}
