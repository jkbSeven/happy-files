package crypto

import (
    "testing"
)

var counter uint32 = 0x00000001
var key = []uint32{0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c, 0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c}
var nonce = []uint32{0x09000000, 0x4a000000, 0x00000000}

func TestQuarterRound (t *testing.T) {
    var a, b, c, d uint32 = 0x11111111, 0x01020304, 0x9b8d6f43, 0x01234567
    var a_expected, b_expected, c_expected, d_expected uint32 = 0xea2a92f4, 0xcb1cf8ce, 0x4581472e, 0x5881c4bb

    exec_round(&a, &b, &c, &d)

    if a != a_expected {
        t.Fatalf("Got a = %x, want a = %x", a, a_expected)
    } else if b != b_expected {
        t.Fatalf("Got b = %x, want b = %x", b, b_expected)
    } else if c != c_expected {
        t.Fatalf("Got c = %x, want c = %x", c, c_expected)
    } else if d != d_expected {
        t.Fatalf("Got d = %x, want d = %x", d, d_expected)
    }
}

func TestStreamGeneration (t *testing.T) {
    expected_block := []uint32{0xe4e7f110, 0x15593bd1, 0x1fdd0f50, 0xc47120a3, 0xc7f4d1c7, 0x0368c033, 0x9aaa2204, 0x4e6cd4c3,
                                0x466482d2, 0x09aa9f07, 0x05d7c214, 0xa2028bd9, 0xd19c12b5, 0xb94e16de, 0xe883d0cb, 0x4e3c50a2}

    block := Generate_stream(key, nonce, counter)

    for i, val := range block {
        if val != expected_block[i] {
            t.Fatalf("Position %d: Got %x instead of %x", i, val, expected_block[i])
        }
    }
}

func TestSerialization (t *testing.T) {
    expected_block := []uint32{0x10f1e7e4, 0xd13b5915, 0x500fdd1f, 0xa32071c4, 0xc7d1f4c7, 0x33c06803, 0x0422aa9a, 0xc3d46c4e,
                                0xd2826446, 0x079faa09, 0x14c2d705, 0xd98b02a2, 0xb5129cd1, 0xde164eb9, 0xcbd083e8, 0xa2503c4e} 
    
    block := Generate_stream(key, nonce, counter)
    Serialize(block)

    for i, val := range block {
        if val != expected_block[i] {
            t.Fatalf("Position %d: Got %x instead of %x", i, val, expected_block[i])
        }
    }
}
