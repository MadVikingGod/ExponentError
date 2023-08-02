package main

import (
	"math"
	"math/bits"
)

func MapToIndexScale0(value float64) int32 {
	const ExponentWidth = 11
	const SignificandWidth = 52
	const SignificandMask = 1<<SignificandWidth - 1
	const ExponentBias = 1<<(ExponentWidth-1) - 1
	const ExponentMask = ((1 << ExponentWidth) - 1) << SignificandWidth
	rawBits := math.Float64bits(value)

	// rawExponent is an 11-bit biased representation of the base-2
	// exponent:
	// - value 0 indicates a subnormal representation or a zero value
	// - value 2047 indicates an Inf or NaN value
	// - value [1, 2046] are offset by ExponentBias (1023)
	rawExponent := (int64(rawBits) & ExponentMask) >> SignificandWidth

	// rawFragment represents (significand-1) for normal numbers,
	// where significand is in the range [1, 2).
	rawFragment := rawBits & SignificandMask

	// Check for subnormal values:
	if rawExponent == 0 {
		// Handle subnormal values: rawFragment cannot be zero
		// unless value is zero.  Subnormal values have up to 52 bits
		// set, so for example greatest subnormal power of two, 0x1p-1023, has
		// rawFragment = 0x8000000000000.  Expressed in 64 bits, the value
		// (rawFragment-1) = 0x0007ffffffffffff has 13 leading zeros.
		rawExponent -= int64(bits.LeadingZeros64(rawFragment-1) - 12)

		// In the example with 0x1p-1023, the preceding expression subtracts
		// (13-12)=1, leaving the rawExponent equal to -1.  The next statement
		// below subtracts `ExponentBias` (1023), leaving `ieeeExponent` equal
		// to -1024, which is the correct upper-inclusive bucket index for
		// the value 0x1p-1023.
	}
	ieeeExponent := int32(rawExponent - ExponentBias)
	// Note that rawFragment and rawExponent cannot both be zero,
	// or else the value is exactly zero, in which case the the ZeroCount
	// bucket is used.
	if rawFragment == 0 {
		// Special case for normal power-of-two values: subtract one.
		return ieeeExponent - 1
	}
	return ieeeExponent
}

// MapToIndexScale0 computes a bucket index at scale 0.
func MapToIndexFrexp(value float64) int {
	// Note: Frexp() rounds submnormal values to the smallest normal
	// value and returns an exponent corresponding to fractions in the
	// range [0.5, 1), whereas an exponent for the range [1, 2), so
	// subtract 1 from the exponent immediately.
	frac, exp := math.Frexp(value)
	exp--

	if frac == 0.5 {
		// Special case for powers of two: they fall into the bucket
		// numbered one less.
		exp--
	}
	return exp
}
