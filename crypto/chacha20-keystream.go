package crypto

import (
    "math/bits"
    "encoding/binary"
)

const rounds = 20
const (
    chacha20_init0 = 0x61707865 
    chacha20_init1 = 0x3320646e
    chacha20_init2 = 0x79622d32
    chacha20_init3 = 0x6b206574
)

func exec_round(a, b, c, d *uint32) {
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

func Serialize(block []uint32) {
    for i, val := range block {
        subblock_little := make([]byte, 4)
        binary.BigEndian.PutUint32(subblock_little, val)
        block[i] = binary.LittleEndian.Uint32(subblock_little)
    }
}

func Generate_stream(key, nonce []uint32, counter uint32) []uint32 {
    var block []uint32
    block = append(block, chacha20_init0, chacha20_init1, chacha20_init2, chacha20_init3)
    block = append(block, key...)
    block = append(block, counter)
    block = append(block, nonce...)

    initial_block := make([]uint32, 16)
    copy(initial_block, block)

    for i := 0; i < rounds / 2; i++ {
        exec_round(&block[0], &block[4], &block[8], &block[12])
        exec_round(&block[1], &block[5], &block[9], &block[13])
        exec_round(&block[2], &block[6], &block[10], &block[14])
        exec_round(&block[3], &block[7], &block[11], &block[15])

        exec_round(&block[0], &block[5], &block[10], &block[15])
        exec_round(&block[1], &block[6], &block[11], &block[12])
        exec_round(&block[2], &block[7], &block[8], &block[13])
        exec_round(&block[3], &block[4], &block[9], &block[14])
    }

    for i, val := range initial_block {
        block[i], _ = bits.Add32(block[i], val, 0)
    }

    return block
}

