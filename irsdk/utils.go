package irsdk

import (
	"encoding/binary"
	"math"
	"strings"
)

func byte4ToInt(in []byte) int {
	return int(binary.LittleEndian.Uint32(in))
}

func byte4ToFloat(in []byte) float32 {
	bits := binary.LittleEndian.Uint32(in)
	return math.Float32frombits(bits)
}

func byte8ToFloat(in []byte) float64 {
	bits := binary.LittleEndian.Uint64(in)
	return math.Float64frombits(bits)
}

//nolint:unused,mnd,gocritic // ok
func byte4toBitField(in []byte) []bool {
	result := make([]bool, 32)
	v := int(binary.LittleEndian.Uint32(in))

	for i := 0; i < 32; i++ {
		result[i] = v&1 == 1
		v = v >> 1
	}

	return result
}

func bytesToString(in []byte) string {
	return strings.TrimRight(string(in), "\x00")
}
