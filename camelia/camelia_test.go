package camelia

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func Test_rotate128Key(t *testing.T) {
	type args struct {
		k [2]uint64
		n int
	}
	tests := []struct {
		name  string
		args  args
		want  uint64
		want1 uint64
	}{
		{
			name: "test rotate",
			args: args{
				k: [2]uint64{uint64(12), uint64(5)},
				n: 0,
			},
			want:  uint64(12),
			want1: uint64(5),
		},
		{
			name: "test rotate 2",
			args: args{
				k: [2]uint64{uint64(1), uint64(1)},
				n: 2,
			},
			want:  uint64(4),
			want1: uint64(4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := rotate128Key(tt.args.k, tt.args.n)
			if got != tt.want {
				t.Errorf("rotate128Key() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("rotate128Key() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}


type testVector struct {
	key, plaintext, ciphertext string
}

func fromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

// Test vectors from RFC3713 - https://www.ietf.org/rfc/rfc3713.txt
var vectors = []struct {
	key, plaintext, ciphertext string
}{
	{
		key:        "0123456789abcdeffedcba9876543210",
		plaintext:  "0123456789abcdeffedcba9876543210",
		ciphertext: "67673138549669730857065648eabe43",
	},
	{
		key:        "0123456789abcdeffedcba98765432100011223344556677",
		plaintext:  "0123456789abcdeffedcba9876543210",
		ciphertext: "b4993401b3e996f84ee5cee7d79b09b9",
	},
	{
		key:        "0123456789abcdeffedcba987654321000112233445566778899aabbccddeeff",
		plaintext:  "0123456789abcdeffedcba9876543210",
		ciphertext: "9acc237dff16d76c20ef7c919e3a7509",
	},
}

func TestVectors(t *testing.T) {
	for i, v := range vectors {
		key := fromHex(v.key)
		plaintext := fromHex(v.plaintext)
		ciphertext := fromHex(v.ciphertext)
		buf := make([]byte, BLOCKSIZE)

		c, err := NewCipher(key)
		if err != nil {
			t.Fatalf("Test vector %d: Failed to create Camellia instance: %s", i, err)
		}

		c.Encrypt(buf, plaintext)
		if !bytes.Equal(ciphertext, buf) {
			t.Fatalf("Test vector %d:\nEncryption failed\nFound:    %s\nExpected: %s", i, hex.EncodeToString(buf), hex.EncodeToString(ciphertext))
		}
		c.Decrypt(buf, buf)
		if !bytes.Equal(plaintext, buf) {
			t.Fatalf("Test vector %d:\nDecryption failed\nFound:    %s\nExpected: %s", i, hex.EncodeToString(buf), hex.EncodeToString(plaintext))
		}
	}
}
