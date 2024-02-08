package main

import(
    "math/bits"
    "encoding/binary"
)

type session struct {
    key []uint32
    counter []uint32
    nonce []uint32
}

func exec_round(a, b, c, d uint32) (uint32, uint32, uint32, uint32){
    a, _ = bits.Add32(a, b, 0)
    d = d ^ a
    d = bits.RotateLeft32(d, 16)

    c, _ = bits.Add32(c, d, 0)
    b = b ^ c
    b = bits.RotateLeft32(b, 12)

    a, _ = bits.Add32(a, b, 0)
    d = d ^ a
    d = bits.RotateLeft32(d, 8)

    c, _ = bits.Add32(c, d, 0)
    b = b ^ c
    b = bits.RotateLeft32(b, 7)
    return a, b, c, d
}

func serialize(block []uint32) {
    for i, val := range block {
        subblock_little := make([]byte, 4)
        binary.BigEndian.PutUint32(subblock_little, val)
        block[i] = binary.LittleEndian.Uint32(subblock_little)
    }
}

func generate_stream(ses *session) []uint32 {
    const rounds = 20
    chacha20_constant := []uint32{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574}

    var block []uint32
    block = append(block, chacha20_constant...)
    block = append(block, ses.key...)
    block = append(block, ses.counter...)
    block = append(block, ses.nonce...)

    initial_block := make([]uint32, 16)
    copy(initial_block, block)

    for i := 0; i < rounds / 2; i++ {
        block[0], block[4], block[8], block[12] = exec_round(block[0], block[4], block[8], block[12])
        block[1], block[5], block[9], block[13] = exec_round(block[1], block[5], block[9], block[13])
        block[2], block[6], block[10], block[14] = exec_round(block[2], block[6], block[10], block[14])
        block[3], block[7], block[11], block[15] = exec_round(block[3], block[7], block[11], block[15])

        block[0], block[5], block[10], block[15] = exec_round(block[0], block[5], block[10], block[15])
        block[1], block[6], block[11], block[12] = exec_round(block[1], block[6], block[11], block[12])
        block[2], block[7], block[8], block[13] = exec_round(block[2], block[7], block[8], block[13])
        block[3], block[4], block[9], block[14] = exec_round(block[3], block[4], block[9], block[14])
    }

    for i, val := range initial_block {
        block[i], _ = bits.Add32(block[i], val, 0)
    }

    serialize(block)
    ses.counter[0]++

    return block
}

func main() {
    var ses session
    ses.key = append(ses.key, 0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c, 0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c)
    ses.counter = append(ses.counter, 0x00000001)
    ses.nonce = append(ses.nonce, 0x09000000, 0x4a000000, 0x00000000)

    generate_stream(&ses)
}
