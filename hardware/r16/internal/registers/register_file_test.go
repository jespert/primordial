package registers_test

import (
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/registers"
	"github.com/jespert/primordial/internal/quality/expect"
)

func TestFile_Write_to_zero_is_hardcoded(t *testing.T) {
	var file registers.File
	file.Write(0, 1)
	expect.Equal(t, 0, file.Read(0))
}

func TestFile_Write_and_read_from_general_register(t *testing.T) {
	var file registers.File
	for i := 1; i < registers.NumRegisters; i++ {
		file.Write(i, int16(i))
		expect.Equal(t, int16(i), file.Read(i))
	}
}
