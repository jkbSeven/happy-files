package crypto

import (
	"fmt"
	"os"
	"testing"
)

var chacha20_init = [4]uint32{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574}
var key = [8]uint32{0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c, 0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c}
var counter uint32 = 0x00000001
var nonce = [3]uint32{0x09000000, 0x4a000000, 0x00000000}
var initial_block [16]uint32

const plaintext = "Ladies and Gentlemen of the class of '99: If I could offer you only one tip for the future, sunscreen would be it."
var ciphertext = []byte{0x6e, 0x2e, 0x35, 0x9a, 0x25, 0x68, 0xf9, 0x80,
                        0x41, 0xba, 0x07, 0x28, 0xdd, 0x0d, 0x69, 0x81,
                        0xe9, 0x7e, 0x7a, 0xec, 0x1d, 0x43, 0x60, 0xc2,
                        0x0a, 0x27, 0xaf, 0xcc, 0xfd, 0x9f, 0xae, 0x0b,
                        0xf9, 0x1b, 0x65, 0xc5, 0x52, 0x47, 0x33, 0xab,
                        0x8f, 0x59, 0x3d, 0xab, 0xcd, 0x62, 0xb3, 0x57,
                        0x16, 0x39, 0xd6, 0x24, 0xe6, 0x51, 0x52, 0xab,
                        0x8f, 0x53, 0x0c, 0x35, 0x9f, 0x08, 0x61, 0xd8,
                        0x07, 0xca, 0x0d, 0xbf, 0x50, 0x0d, 0x6a, 0x61,
                        0x56, 0xa3, 0x8e, 0x08, 0x8a, 0x22, 0xb6, 0x5e,
                        0x52, 0xbc, 0x51, 0x4d, 0x16, 0xcc, 0xf8, 0x06,
                        0x81, 0x8c, 0xe9, 0x1a, 0xb7, 0x79, 0x37, 0x36,
                        0x5a, 0xf9, 0x0b, 0xbf, 0x74, 0xa3, 0x5b, 0xe6,
                        0xb4, 0x0b, 0x8e, 0xed, 0xf2, 0x78, 0x5e, 0x42,
                        0x87, 0x4d,
                    }

func initBlock() {
    for i, v := range chacha20_init {
        initial_block[i] = v
    }
    for i, v := range key {
        initial_block[i + 4] = v
    }
    initial_block[12] = counter
    for i, v := range nonce {
        initial_block[i + 13] = v
    }
}

func TestMain(m *testing.M) {
    initBlock()
    code := m.Run()
    os.Exit(code)
}

func TestQuarterRound (t *testing.T) {
    var a, b, c, d uint32 = 0x11111111, 0x01020304, 0x9b8d6f43, 0x01234567
    var a_expected, b_expected, c_expected, d_expected uint32 = 0xea2a92f4, 0xcb1cf8ce, 0x4581472e, 0x5881c4bb

    execRound(&a, &b, &c, &d)

    if a != a_expected {
        t.Fatalf("a: got %x, want %x", a, a_expected)
    } else if b != b_expected {
        t.Fatalf("b: got %x, want %x", b, b_expected)
    } else if c != c_expected {
        t.Fatalf("c: got %x, want %x", c, c_expected)
    } else if d != d_expected {
        t.Fatalf("d: got %x, want %x", d, d_expected)
    }
}

func TestStreamGeneration (t *testing.T) {
    expected_block := [16]uint32{0xe4e7f110, 0x15593bd1, 0x1fdd0f50, 0xc47120a3, 0xc7f4d1c7, 0x0368c033, 0x9aaa2204, 0x4e6cd4c3,
                                0x466482d2, 0x09aa9f07, 0x05d7c214, 0xa2028bd9, 0xd19c12b5, 0xb94e16de, 0xe883d0cb, 0x4e3c50a2}

    block := generateStream(initial_block)

    for i, val := range block {
        if val != expected_block[i] {
            t.Fatalf("Position %d: got %x instead of %x", i, val, expected_block[i])
        }
    }
}

func TestSerialization (t *testing.T) {
    expected_block := [16]uint32{0x10f1e7e4, 0xd13b5915, 0x500fdd1f, 0xa32071c4, 0xc7d1f4c7, 0x33c06803, 0x0422aa9a, 0xc3d46c4e,
                                0xd2826446, 0x079faa09, 0x14c2d705, 0xd98b02a2, 0xb5129cd1, 0xde164eb9, 0xcbd083e8, 0xa2503c4e} 
    
    block := generateStream(initial_block)
    serialize(&block)

    for i, val := range block {
        if val != expected_block[i] {
            t.Fatalf("Position %d: got %x instead of %x", i, val, expected_block[i])
        }
    }
}

func TestEncrypt (t *testing.T) {
    nonce[0] = 0
    encrypted := Encrypt([]byte(plaintext), key, nonce)
    for i, v := range encrypted {
        if v != ciphertext[i] {
            fmt.Printf("Expected: %v\n\nGot: %v\n", ciphertext, encrypted)
            t.Fatalf("Encryption is not working properly")
        }
    }
    decrypted := Decrypt(encrypted, key, nonce)
    if string(decrypted) != plaintext {
        fmt.Println("Original:", []byte(plaintext))
        fmt.Println("Decrypted:", decrypted)
        t.Fatalf("Initial plaintext differs from decrypted plaintext")
    }
    nonce[0] = 0x09000000
}
