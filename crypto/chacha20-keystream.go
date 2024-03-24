package crypto

import (
	"encoding/binary"
	"math/bits"
)

const rounds = 20
var initParams = [4]uint32{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574}

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

func serialize(block *[16]uint32) {
    for i, val := range block {
        subblock_bytes := make([]byte, 4)
        binary.BigEndian.PutUint32(subblock_bytes, val)
        block[i] = binary.LittleEndian.Uint32(subblock_bytes)
    }
}

func generateStream(initial_block [16]uint32) [16]uint32 {
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

func uintToBytes(stream [16]uint32) []byte {
    size := len(stream) * 4
    out := make([]byte, size)
    start := 0
    for _, v := range stream {
        binary.BigEndian.PutUint32(out[start:start+4], v)
        start += 4
    }
    return out
}

func Encrypt(plaintext []byte, key [8]uint32, nonce [3]uint32) []byte {
    plainLen := len(plaintext)
    blockCount := plainLen / 64

    if plainLen % 64 != 0 {
        blockCount++
    }

    block := [16]uint32{
        initParams[0], initParams[1], initParams[2], initParams[3],
        key[0], key[1], key[2], key[3], 
        key[4], key[5], key[6], key[7], 
        1, nonce[0], nonce[1], nonce[2],
    }

    ciphertext := make([]byte, 0, plainLen)
    for i := range blockCount {
        stream := generateStream(block)
        serialize(&stream)
        streamBytes := uintToBytes(stream)

        if i == blockCount - 1 {
            for j, v := range plaintext[i * 64:] {
                ciphertext = append(ciphertext, v ^ streamBytes[j])
            }
            return ciphertext
        }

        for j, v := range streamBytes {
            ciphertext = append(ciphertext, v ^ plaintext[j + 64 * i])
        }
        block[12]++
    }

    return ciphertext
}

func Decrypt(ciphertext []byte, key [8]uint32, nonce [3]uint32) []byte {
    cipherLen := len(ciphertext)
    blockCount := cipherLen / 64

    if cipherLen % 64 != 0 {
        blockCount++
    }

    block := [16]uint32{
        initParams[0], initParams[1], initParams[2], initParams[3],
        key[0], key[1], key[2], key[3], 
        key[4], key[5], key[6], key[7], 
        1, nonce[0], nonce[1], nonce[2],
    }

    plaintext := make([]byte, 0, cipherLen)
    for i := range blockCount {
        stream := generateStream(block)
        serialize(&stream)
        streamBytes := uintToBytes(stream)

        if i == blockCount - 1 {
            for j, v := range ciphertext[i * 64:] {
                plaintext = append(plaintext, v ^ streamBytes[j])
            }
            return plaintext
        }

        for j, v := range streamBytes {
            plaintext = append(plaintext, v ^ ciphertext[j + 64 * i])
        }
        block[12]++
    }
    return plaintext
}
