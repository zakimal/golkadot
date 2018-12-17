package u8util

import (
	"fmt"
	"math"
	"math/big"

	"github.com/c3systems/go-substrate/common/bn"
)

var alphabet = "0123456789abcdef"

// Concat concatenates multiple uint8 slices into a new slice
func Concat(slices ...[]uint8) []uint8 {
	var ret []uint8
	for _, x := range slices {
		ret = append(ret, x...)
	}

	return ret
}

// FixLength shifts a uint8 slice to a specific bitLength.  Returns a uint8 slice with the specified number of bits contained in the return value. (If bitLength is -1, length checking is not done). Values with more bits are trimmed to the specified length.
func FixLength(value []uint8, bitLength int, atStart bool) []uint8 {
	byteLength := int(math.Ceil(float64(bitLength) / float64(8)))

	if bitLength == -1 || len(value) == byteLength {
		return value
	}

	if len(value) > byteLength {
		return value[0:byteLength]
	}

	result := make([]uint8, byteLength)
	if atStart {
		return append(result, value...)
	}

	start := byteLength - len(value)
	for i := 0; i+start < byteLength; i++ {
		result[i+start] = value[i]
	}

	return result
}

// ToString creates a utf-8 string from a uint8 slice.
func ToString(value []uint8) string {
	return string(value)
}

// ToHex creates a hex string from a uint8 slice.
func ToHex(value []uint8, bitLength int, isPrefixed bool) string {
	byteLength := int(math.Ceil(float64(bitLength) / float64(8)))

	if byteLength > 0 && len(value) > byteLength {
		halfLength := int(math.Ceil(float64(byteLength) / float64(2)))

		return fmt.Sprintf("%s…%s", ToHex(value[0:halfLength], -1, isPrefixed), ToHex(value[len(value)-halfLength:], -1, false))
	}

	result := ""
	if isPrefixed {
		result = "0x"
	}

	for i := 0; i < len(value); i++ {
		v := value[i]
		result = result + string(alphabet[v>>4]) + string(alphabet[v&0xf])
	}

	return result
}

// ToBN creates a utf-8 string from a uint8 slice.
func ToBN(value []uint8, isLittleEndian bool) *big.Int {
	return bn.ToBN(value, isLittleEndian)
}