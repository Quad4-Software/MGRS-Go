// Package mgrs converts between WGS84 geographic coordinates and the Military
// Grid Reference System (MGRS) used with UTM. It follows the letter-grid rules
// described in NGA references on UTM/MGRS (e.g. TM8358.1, Universal Grids and
// Grid Reference Systems).
//
// Coverage is latitudes [-80°, 84°) using UTM (the standard MGRS UTM belt).
// Latitudes outside that range normally use UPS polar grid references, which
// are not handled here.
//
// Accuracy and verification:
// Transverse Mercator uses WGS84 with constant-folded fictitious eastings and
// Sincos-based derivatives where possible; Norway/Svalbard zone logic matches the
// usual GeographicLib / MSP GEOTRANS conventions. Encode path uses compact
// fixed-capacity ASCII output plus a digit-scale lookup instead of naive pow loops.
//
// EncodeBytes and AppendEncode avoid the Unicode string allocation of Encode when
// you only need ASCII bytes or are batching coordinates. EncodeTo can write into
// a caller-owned fixed buffer for a zero-allocation encode path.
//
// Tests compare metre outputs against the PROJ `proj` binary when installed
// (mgrs.proj_test.go) and include fuzzing plus property checks.
package mgrs
