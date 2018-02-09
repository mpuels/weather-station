package pulse

import (
	"log"
	"testing"

	"github.com/jckuester/weather-station/protocol"
	"github.com/stretchr/testify/assert"
)

func TestProtocolMatches(t *testing.T) {
	p := &PulseInfo{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}
	pc := &protocol.Protocol{
		SeqLength: []int{76},
		Lengths:   []int{496, 2048, 4068, 8960},
	}

	assert.Equal(t, true, protocolMatches(p, pc))
}

func TestProtocolMatches_pulseSeqTooShort(t *testing.T) {
	p := &PulseInfo{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "010202010102020102",
	}
	pc := &protocol.Protocol{
		SeqLength: []int{76},
		Lengths:   []int{496, 2048, 4068, 8960},
	}

	assert.Equal(t, false, protocolMatches(p, pc))
}

func TestProtocolMatches_pulseLengthDeviationTooHigh(t *testing.T) {
	p := &PulseInfo{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}
	pc := &protocol.Protocol{
		SeqLength: []int{76},
		Lengths:   []int{496, 2048, 2000, 8960},
	}

	assert.Equal(t, false, protocolMatches(p, pc))
}

func TestProtocolMatches_numberOfPuleLengthDiffer(t *testing.T) {
	p := &PulseInfo{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}
	pc := &protocol.Protocol{
		SeqLength: []int{76},
		Lengths:   []int{496, 2048, 2000},
	}

	assert.Equal(t, false, protocolMatches(p, pc))
}

func TestMapPulse(t *testing.T) {
	pulseSeq := "020101020201020201010103"
	bits := "10011011000"

	var pulsesToBinaryMapping = map[string]string{
		"01": "0", // binary 0
		"02": "1", // binary 1
		"03": "",  // footer
	}

	mapped, err := Map(pulseSeq, pulsesToBinaryMapping)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, mapped, bits)
}

func TestPrepareCompressedPulse(t *testing.T) {
	input := "255 2904 1388 771 11346 0 0 0 0100020002020000020002020000020002000202000200020002000200000202000200020000020002000200020002020002000002000200000002000200020002020002000200020034"

	p, _ := PrepareCompressed(input)

	assert.Equal(t, []int{255, 771, 1388, 2904, 11346}, p.Lengths)
	assert.Equal(t, "0300020002020000020002020000020002000202000200020002000200000202000200020000020002000200020002020002000002000200000002000200020002020002000200020014", p.Seq)
}

func TestInvalidCharacters(t *testing.T) {
	pC := "544 4128 2100 100 140 320 808 188 01020202010202020202020202020202020101020102010101010J�G_YJ�Üxx�1��Nz�8��&[��"

	_, err := PrepareCompressed(pC)

	assert.Error(t, err)
}

func TestSortIndices(t *testing.T) {
	a := []int{200, 600, 500}

	sortedIndices := sortIndices(a)

	assert.Equal(t, sortedIndices, []int{0, 2, 1})
}
