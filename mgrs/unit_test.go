package mgrs

import (
	"errors"
	"testing"
)

func TestUnitEncode_digitPairsBoundary(t *testing.T) {
	_, err := Encode(10, 20, -1)
	if !errors.Is(err, ErrInvalidDigitPairs) {
		t.Fatalf("got %v want ErrInvalidDigitPairs", err)
	}
	_, err = Encode(10, 20, MaxDigitPairs+3)
	if !errors.Is(err, ErrInvalidDigitPairs) {
		t.Fatalf("got %v want ErrInvalidDigitPairs", err)
	}
}

func TestUnitEncode_latitudeBelt(t *testing.T) {
	_, err := Encode(-82, 40, DefaultDigitPairs)
	if !errors.Is(err, ErrLatOutOfUTMRange) {
		t.Fatalf("got %v", err)
	}
	_, err = Encode(85, 0, DefaultDigitPairs)
	if !errors.Is(err, ErrLatOutOfUTMRange) {
		t.Fatalf("got %v", err)
	}
}

func TestUnitDecode_whitespaceInvariant(t *testing.T) {
	s, err := Encode(-33.8136, 144.9631, DefaultDigitPairs)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := Decode(s, false)
	if err != nil {
		t.Fatal(err)
	}
	spaced, err := Decode("  \t"+s+"\n  ", false)
	if err != nil {
		t.Fatal(err)
	}
	if spaced != ref {
		t.Fatalf("whitespace divergence %+v vs %+v", ref, spaced)
	}
}

func TestUnitDecode_oddDigitTail(t *testing.T) {
	_, err := Decode("03UXD982817461", false)
	if err == nil {
		t.Fatal("expected error for odd digit count")
	}
}

func TestUnitDecode_zoneLeadingZeroRejected(t *testing.T) {
	_, err := Decode("00UXD98281746211", false)
	if !errors.Is(err, ErrInvalidMGRS) {
		t.Fatalf("got %v want ErrInvalidMGRS", err)
	}
}
