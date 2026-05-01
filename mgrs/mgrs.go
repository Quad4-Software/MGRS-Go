package mgrs

// Point expresses a geographic location on WGS84 in decimal degrees.
type Point struct {
	Lat float64
	Lon float64
}

const DefaultDigitPairs = 5
const MaxEncodedLen = 5 + MaxDigitPairs*2

func LongitudeZone(latDeg, lonDeg float64) (int, error) {
	switch {
	case latDeg < -80 || latDeg >= 84:
		return 0, ErrLatOutOfUTMRange
	default:
		return longitudeZoneStandard(latDeg, lonDeg), nil
	}
}

func marshalMGRS(latDeg, lonDeg float64, digitPairs int) ([]byte, error) {
	switch {
	case latDeg < -80 || latDeg >= 84:
		return nil, ErrLatOutOfUTMRange
	case digitPairs < 0 || digitPairs > MaxDigitPairs:
		return nil, ErrInvalidDigitPairs
	}
	northhemi := latDeg >= 0
	zone := longitudeZoneStandard(latDeg, lonDeg)
	if zone < 1 || zone > 60 {
		return nil, ErrLatOutOfUTMRange
	}
	east, north := forwardTM(latDeg, lonDeg, zone, !northhemi)
	return encodeUTM(zone, latDeg, northhemi, east, north, digitPairs)
}

func marshalMGRSInto(dst []byte, latDeg, lonDeg float64, digitPairs int) ([]byte, error) {
	switch {
	case latDeg < -80 || latDeg >= 84:
		return nil, ErrLatOutOfUTMRange
	case digitPairs < 0 || digitPairs > MaxDigitPairs:
		return nil, ErrInvalidDigitPairs
	}
	northhemi := latDeg >= 0
	zone := longitudeZoneStandard(latDeg, lonDeg)
	if zone < 1 || zone > 60 {
		return nil, ErrLatOutOfUTMRange
	}
	east, north := forwardTM(latDeg, lonDeg, zone, !northhemi)
	return encodeUTMInto(dst, zone, latDeg, northhemi, east, north, digitPairs)
}

// Encode formats (latDeg, lonDeg) as an MGRS reference string without spaces.
func Encode(latDeg, lonDeg float64, digitPairs int) (string, error) {
	buf, err := marshalMGRS(latDeg, lonDeg, digitPairs)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// EncodeBytes is like Encode but returns the ASCII bytes directly, avoiding an
// extra string allocation when callers need []byte downstream.
func EncodeBytes(latDeg, lonDeg float64, digitPairs int) ([]byte, error) {
	return marshalMGRS(latDeg, lonDeg, digitPairs)
}

// EncodeTo writes an encoded MGRS reference into dst and returns the number of
// bytes written. dst must have length >= 5+digitPairs*2. This allows a
// zero-allocation encode path when the caller reuses a fixed buffer.
func EncodeTo(dst []byte, latDeg, lonDeg float64, digitPairs int) (int, error) {
	switch {
	case digitPairs < 0 || digitPairs > MaxDigitPairs:
		return 0, ErrInvalidDigitPairs
	}
	need := 5 + digitPairs*2
	if len(dst) < need {
		return 0, ErrShortBuffer
	}
	out, err := marshalMGRSInto(dst[:0], latDeg, lonDeg, digitPairs)
	if err != nil {
		return 0, err
	}
	return len(out), nil
}

// AppendEncode appends a newly encoded reference to dst and returns the
// enlarged slice (convenient for buffering many references with one allocator).
func AppendEncode(dst []byte, latDeg, lonDeg float64, digitPairs int) ([]byte, error) {
	return marshalMGRSInto(dst, latDeg, lonDeg, digitPairs)
}

func Decode(reference string, centerCell bool) (Point, error) {
	data := compactUpperAsciiFromString(reference)
	switch len(data) {
	case 0:
		return Point{}, ErrInvalidMGRS
	default:
		return decodeNormalizedUpper(data, centerCell)
	}
}
