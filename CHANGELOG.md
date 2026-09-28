# Changelog

All notable changes are documented here. The format follows Keep a Changelog
and versions follow SemVer.

## [1.0.0] - 2026-09-28

First tagged release.

### Added

- Full-globe WGS84 MGRS encode and decode: UTM between 80S and 84N latitude,
  polar UPS outside that belt (bands A/B south, Y/Z north).
- Norway and Svalbard zone widening, grid-zone lettering, and upper-bound
  coordinate nudges matching GeographicLib / MSP GEOTRANS practice.
- Public API: Encode, EncodeGrid, Decode, DecodeParts, LatLonToGrid,
  LongitudeZone, FormatSpaced.
- Zero-allocation encode variants: EncodeTo, EncodeGridTo, AppendEncode,
  AppendEncodeGrid, EncodeBytes, EncodeGridBytes.
- Grid API for direct UTM/UPS metre input, including ZoneUPS for polar points.
- Example CLI under example/ with encode, decode, pretty, and bench modes.
- Checked-in golden corpus (mgrs/testdata/golden_mgrs.jsonl) and differential
  tests against PROJ and GeographicLib GeoConvert. Set MGRS_REQUIRE_DIFF=1 to
  fail instead of skip when the tools are missing.

### CI

- golangci-lint, staticcheck, gofumpt/goimports format check, race detector,
  bounded fuzz, CodeQL, OpenSSF Scorecard, and zizmor workflow audits.
- Generated shields endpoint badges via scripts/gen-badges.sh (make badges),
  published to the badges branch.

[1.0.0]: https://github.com/Quad4-Software/MGRS-Go/releases/tag/v1.0.0
