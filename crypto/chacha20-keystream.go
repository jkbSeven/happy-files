package crypto

import (
    "math/bits"
    "encoding/binary"
)

const rounds = 20

func execRound(a, b, c, d *uint32) {
    *a, _ = bits.Add32(*a, *b, 0)
    *d = *d ^ *a
    *d = bits.RotateLeft32(*d, 16)

    *c, _ = bits.Add32(*c, *d, 0)
    *b = *b ^ *c
    *b = bits.RotateLeft32(*b, 12)

    *a, _ = bits.Add32(*a, *b, 0)
    *d = *d ^ *a
    *d = bits.RotateLeft32(*d, 8)

    *c, _ = bits.Add32(*c, *d, 0)
    *b = *b ^ *c
    *b = bits.RotateLeft32(*b, 7)
}

func Serialize(block *[16]uint32) {
    for i, val := range block {
        subblock_bytes := make([]byte, 4)
        binary.BigEndian.PutUint32(subblock_bytes, val)
        block[i] = binary.LittleEndian.Uint32(subblock_bytes)
    }
}

func GenerateStream(initial_block [16]uint32) [16]uint32 {
    var block [16]uint32

    for i, v := range initial_block {
        block[i] = v
    }

    for i := 0; i < rounds / 2; i++ {
        execRound(&block[0], &block[4], &block[8], &block[12])
        execRound(&block[1], &block[5], &block[9], &block[13])
        execRound(&block[2], &block[6], &block[10], &block[14])
        execRound(&block[3], &block[7], &block[11], &block[15])

        execRound(&block[0], &block[5], &block[10], &block[15])
        execRound(&block[1], &block[6], &block[11], &block[12])
        execRound(&block[2], &block[7], &block[8], &block[13])
        execRound(&block[3], &block[4], &block[9], &block[14])
    }

    for i, val := range initial_block {
        block[i], _ = bits.Add32(block[i], val, 0)
    }

    return block
}

