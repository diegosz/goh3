package h3

import "strconv"

// uint64ToHex encodes value as an hex string without 0x prefix.
func uint64ToHex(value uint64) string {
	enc := make([]byte, 0, 10)
	return string(strconv.AppendUint(enc, value, 16))
}

// hexToUint64 decodes an hex string without 0x prefix to a uint64.
func hexToUint64(value string) (uint64, error) {
	return strconv.ParseUint(value, 16, 64)
}
