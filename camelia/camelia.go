package camelia

import (
	"crypto/cipher"
	"encoding/binary"
	"math/bits"
	"strconv"
)

const (
	MASK8   = 0xff
	MASK32  = 0xffffffff
	MASK64  = 0xffffffffffffffff
	MASK128 = 0xffffffffffffffffffffffffffffffff

	C1 uint64 = 0xA09E667F3BCC908B
	C2 uint64 = 0xB67AE8584CAA73B2
	C3 uint64 = 0xC6EF372FE94F82BE
	C4 uint64 = 0x54FF53A5F1D36F1C
	C5 uint64 = 0x10E527FADE682D1D
	C6 uint64 = 0xB05688C2B3E6C1FD

	BLOCKSIZE = 16
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return "camelia: invalid key size " + strconv.Itoa(int(k))
}

type cameliaCipher struct {
	kw   [5]uint64
	k    [25]uint64
	ke   [7]uint64
	klen int

	ka, kb, kl, kr [2]uint64
}

func newCamelliaCipher(klen int) *cameliaCipher {
	return &cameliaCipher{klen: klen}
}

func (c cameliaCipher) BlockSize() int {
	return BLOCKSIZE
}

// src = 128-битное (16-байтовое сообщение)
func (c cameliaCipher) Encrypt(dst, src []byte) {
	crypt(c, dst, src)
}

func crypt(c cameliaCipher, dst, src []byte){
	//   D1 = M >> 64;            // Шифруемое сообщение делится на две 64-битные части
	//  D2 = M & MASK64;
	//  D1 = D1 ^ kw1;           // Предварительное забеливание
	//  D2 = D2 ^ kw2;
	//  D2 = D2 ^ F(D1, k1);
	//  D1 = D1 ^ F(D2, k2);
	//  D2 = D2 ^ F(D1, k3);
	//  D1 = D1 ^ F(D2, k4);
	//  D2 = D2 ^ F(D1, k5);
	//  D1 = D1 ^ F(D2, k6);
	//  D1 = FL   (D1, ke1);     // FL
	//  D2 = FLINV(D2, ke2);     // FLINV
	//  D2 = D2 ^ F(D1, k7);
	//  D1 = D1 ^ F(D2, k8);
	//  D2 = D2 ^ F(D1, k9);
	//  D1 = D1 ^ F(D2, k10);
	//  D2 = D2 ^ F(D1, k11);
	//  D1 = D1 ^ F(D2, k12);
	//  D1 = FL   (D1, ke3);     // FL
	//  D2 = FLINV(D2, ke4);     // FLINV
	//  D2 = D2 ^ F(D1, k13);
	//  D1 = D1 ^ F(D2, k14);
	//  D2 = D2 ^ F(D1, k15);
	//  D1 = D1 ^ F(D2, k16);
	//  D2 = D2 ^ F(D1, k17);
	//  D1 = D1 ^ F(D2, k18);
	//  D2 = D2 ^ kw3;           // Финальное забеливание
	//  D1 = D1 ^ kw4;
	//  C = (D2 << 64) | D1;
	d1 := binary.BigEndian.Uint64(src[0:])
	d2 := binary.BigEndian.Uint64(src[8:])

	d1 ^= c.kw[1]
	d2 ^= c.kw[2]

	d2 = d2 ^ f(d1, c.k[1])
	d1 = d1 ^ f(d2, c.k[2])
	d2 = d2 ^ f(d1, c.k[3])
	d1 = d1 ^ f(d2, c.k[4])
	d2 = d2 ^ f(d1, c.k[5])
	d1 = d1 ^ f(d2, c.k[6])

	d1 = fl(d1, c.ke[1])
	d2 = flinv(d2, c.ke[2])

	d2 = d2 ^ f(d1, c.k[7])
	d1 = d1 ^ f(d2, c.k[8])
	d2 = d2 ^ f(d1, c.k[9])
	d1 = d1 ^ f(d2, c.k[10])
	d2 = d2 ^ f(d1, c.k[11])
	d1 = d1 ^ f(d2, c.k[12])

	d1 = fl(d1, c.ke[3])
	d2 = flinv(d2, c.ke[4])

	d2 = d2 ^ f(d1, c.k[13])
	d1 = d1 ^ f(d2, c.k[14])
	d2 = d2 ^ f(d1, c.k[15])
	d1 = d1 ^ f(d2, c.k[16])
	d2 = d2 ^ f(d1, c.k[17])
	d1 = d1 ^ f(d2, c.k[18])

	if c.klen > 16 {
		// 24 or 32

		d1 = fl(d1, c.ke[5])
		d2 = flinv(d2, c.ke[6])

		d2 = d2 ^ f(d1, c.k[19])
		d1 = d1 ^ f(d2, c.k[20])
		d2 = d2 ^ f(d1, c.k[21])
		d1 = d1 ^ f(d2, c.k[22])
		d2 = d2 ^ f(d1, c.k[23])
		d1 = d1 ^ f(d2, c.k[24])
	}

	d2 = d2 ^ c.kw[3]
	d1 = d1 ^ c.kw[4]

	binary.BigEndian.PutUint64(dst[0:], d2)
	binary.BigEndian.PutUint64(dst[8:], d1)
}

func (c cameliaCipher) Decrypt(dst, src []byte) {
	if c.klen == 128/8 {
		c.kw[0], c.kw[2] = c.kw[2], c.kw[0]
		c.kw[1], c.kw[3] = c.kw[3], c.kw[1]
		c.k[0], c.k[17] = c.k[17], c.k[0]
		// 2 - 17
		c.k[1], c.k[16] = c.k[16], c.k[1]
		// 3 - 16
		c.k[2], c.k[15] = c.k[15], c.k[2]
		// 4 - 15
		c.k[3], c.k[14] = c.k[14], c.k[3]
		// 5 - 14
		c.k[4], c.k[13] = c.k[13], c.k[4]
		// 6 - 13
		c.k[5], c.k[12] = c.k[12], c.k[5]
		// 7 - 12
		c.k[6], c.k[11] = c.k[11], c.k[6]
		// 8 - 11
		c.k[7], c.k[10] = c.k[10], c.k[7]
		// 9 - 10
		c.k[8], c.k[9] = c.k[9], c.k[8]

		// ke1 <-> ke4
		// ke2 <-> ke3
		c.ke[0], c.ke[3] = c.ke[3], c.ke[0]
		c.ke[1], c.ke[2] = c.ke[2], c.ke[1]

		crypt(c, dst, src)
		return
	}

	// 192 or 256
	c.kw[0], c.kw[2] = c.kw[2], c.kw[0]
	c.kw[1], c.kw[3] = c.kw[3], c.kw[1]
	c.k[0], c.k[23] = c.k[23], c.k[0]
	c.k[1], c.k[22] = c.k[22], c.k[1]
	c.k[2], c.k[21] = c.k[21], c.k[2]
	c.k[3], c.k[20] = c.k[20], c.k[3]
	c.k[4], c.k[19] = c.k[19], c.k[4]
	c.k[5], c.k[18] = c.k[18], c.k[5]
	c.k[6], c.k[17] = c.k[17], c.k[6]
	c.k[7], c.k[16] = c.k[16], c.k[7]
	c.k[8], c.k[15] = c.k[15], c.k[8]
	c.k[9], c.k[14] = c.k[14], c.k[9]
	c.k[10], c.k[13] = c.k[13], c.k[10]
	c.k[11], c.k[12] = c.k[12], c.k[11]

	c.ke[0], c.ke[5] = c.ke[5], c.ke[0]
	c.ke[1], c.ke[4] = c.ke[4], c.ke[1]
	c.ke[2], c.ke[3] = c.ke[3], c.ke[2]

	crypt(c, dst, src)
}

// циклический сдвиг 128-битного ключа
func rotate128Key(k [2]uint64, n int) (uint64, uint64) {
	if n > 64 {
		n %= 64
		// поворот порядка
		k[0], k[1] = k[1], k[0]
	}
	//t := *r0 >> (32 - n)
	//*r0 = (*r0 << n) | (*r1 >> (32 - n))
	//*r1 = (*r1 << n) | (*r2 >> (32 - n))
	//*r2 = (*r2 << n) | (*r3 >> (32 - n))
	//*r3 = (*r3 << n) | t

	// сдвигаем
	t := k[0] >> (64 - n)
	r := (k[0] << n) | (k[1] >> (64 - n))
	l := (k[1] << n) | t

	return r, l
}

// helpKeys128 - генерация вспомогательных ключей для ключа размером 128 бит
func (c cameliaCipher) helpKeys128(ka, kl [2]uint64) {
	c.kw[1], c.kw[2] = rotate128Key(kl, 0)

	c.k[1], c.k[2] = rotate128Key(ka, 0)
	c.k[3], c.k[4] = rotate128Key(kl, 15)
	c.k[5], c.k[6] = rotate128Key(ka, 15)

	c.ke[1], c.ke[2] = rotate128Key(ka, 30)

	c.k[7], c.k[8] = rotate128Key(kl, 45)
	c.k[9], _ = rotate128Key(ka, 45)
	_, c.k[10] = rotate128Key(kl, 60)
	c.k[11], c.k[12] = rotate128Key(ka, 60)

	c.ke[3], c.ke[4] = rotate128Key(kl, 77)

	c.k[13], c.k[14] = rotate128Key(kl, 94)
	c.k[15], c.k[16] = rotate128Key(ka, 94)
	c.k[17], c.k[18] = rotate128Key(kl, 111)

	c.kw[3], c.kw[4] = rotate128Key(ka, 111)
}

// helpKeys256 - генерация вспомогательных ключей для ключа размером 194 or 256
func (c cameliaCipher) helpKeys256(ka, kb, kl, kr [2]uint64) {
	c.kw[1], c.kw[2] = rotate128Key(kl, 0)

	c.k[1], c.k[2] = rotate128Key(kb, 0)
	c.k[3], c.k[4] = rotate128Key(kr, 15)
	c.k[5], c.k[6] = rotate128Key(ka, 15)

	c.ke[1], c.ke[2] = rotate128Key(kr, 30)

	c.k[7], c.k[8] = rotate128Key(kb, 30)
	c.k[9], c.k[10] = rotate128Key(kl, 45)
	c.k[11], c.k[12] = rotate128Key(ka, 45)

	c.ke[3], c.ke[4] = rotate128Key(kl, 60)

	c.k[13], c.k[14] = rotate128Key(kr, 60)
	c.k[15], c.k[16] = rotate128Key(kb, 60)
	c.k[17], c.k[18] = rotate128Key(kl, 77)

	c.ke[5], c.ke[6] = rotate128Key(ka, 77)

	c.k[19], c.k[20] = rotate128Key(kr, 94)
	c.k[21], c.k[22] = rotate128Key(ka, 94)
	c.k[23], c.k[24] = rotate128Key(kl, 111)

	c.kw[3], c.kw[4] = rotate128Key(kb, 111)
}

// kAkB - вычисления вспомогательных 128-битных чисел KA, KB
func kAkB(kl, kr [2]uint64) (ka [2]uint64, kb [2]uint64) {
	var d1, d2 uint64

	d1 = kl[0] ^ kr[0]
	d2 = kl[1] ^ kr[1]

	d2 = d2 ^ f(d1, C1)
	d1 = d1 ^ f(d2, C2)

	d1 = d1 ^ (kl[0])
	d2 = d2 ^ (kl[1])
	d2 = d2 ^ f(d1, C3)
	ff := f(d2, C4)
	d1 = d1 ^ ff
	ka[0] = d1
	ka[1] = d2
	d1 = ka[0] ^ kr[0]
	d2 = ka[1] ^ kr[1]
	d2 = d2 ^ f(d1, C5)
	d1 = d1 ^ f(d2, C6)

	kb[0] = d1
	kb[1] = d2

	return ka, kb
}

// 128 бит = 16 байт
// 256 бит = 32 байт
// 192 бит = 24 байт
func NewCipher(key []byte) (cipher.Block, error) {
	k := len(key)

	switch k {
	default:
		return nil, KeySizeError(k)
	case 16, 24, 32:
		break
	}

	var kl [2]uint64
	var kr [2]uint64
	var ka [2]uint64
	var kb [2]uint64

	kl[0] = binary.BigEndian.Uint64(key[0:])
	kl[1] = binary.BigEndian.Uint64(key[8:])

	switch k {
	case 24:
		kr[0] = binary.BigEndian.Uint64(key[16:])
		kr[1] = ^kr[0]
	case 32:
		kr[0] = binary.BigEndian.Uint64(key[16:])
		kr[1] = binary.BigEndian.Uint64(key[24:])
	}

	ka, kb = kAkB(kl, kr)

	c := newCamelliaCipher(k)

	if k == 16 {
		c.helpKeys128(ka, kl)
	} else {
		// 24 or 32
		c.helpKeys256(ka, kb, kl, kr)
	}

	return c, nil
}

func fl(flIn, ke uint64) uint64 {
	//var x1, x2 as 32-bit unsigned integer;
	//var k1, k2 as 32-bit unsigned integer;
	//x1 = FL_IN >> 32;
	//x2 = FL_IN & MASK32;
	//k1 = KE >> 32;
	//k2 = KE & MASK32;
	//x2 = x2 ^ ((x1 & k1) <<< 1);
	//x1 = x1 ^ (x2 | k2);
	//FL_OUT = (x1 << 32) | x2;

	var x1, x2 uint32

	x1 = uint32(flIn >> 32)
	x2 = uint32(flIn & MASK32)
	k1 := uint32(ke >> 32)
	k2 := uint32(ke & MASK32)

	x2 = x2 ^ bits.RotateLeft32(x1&k1, 1)
	x1 = x1 ^ (x2 | k2)

	return (uint64(x1) << 32) | uint64(x2)
}

func flinv(flinvIn, ke uint64) uint64 {
	//var y1, y2 as 32-bit unsigned integer;
	//var k1, k2 as 32-bit unsigned integer;
	//y1 = FLINV_IN >> 32;
	//y2 = FLINV_IN & MASK32;
	//k1 = KE >> 32;
	//k2 = KE & MASK32;
	//y1 = y1 ^ (y2 | k2);
	//y2 = y2 ^ ((y1 & k1) <<< 1);
	//FLINV_OUT = (y1 << 32) | y2;
	var y1, y2, k1, k2 uint32

	y1 = uint32(flinvIn >> 32)
	y2 = uint32(flinvIn & MASK32)
	k1 = uint32(ke >> 32)
	k2 = uint32(ke & MASK32)

	y1 = y1 ^ (y2 | k2)
	y2 = y2 ^ (bits.RotateLeft32(y1&k1, 1))

	return (uint64(y1) << 32) | uint64(y2)
}

func f(fIn, ke uint64) uint64 {
	// 	   x  = F_IN ^ KE;
	x := fIn ^ ke
	/*
	   t1 =  x >> 56;
	   t2 = (x >> 48) & MASK8;
	   t3 = (x >> 40) & MASK8;
	   t4 = (x >> 32) & MASK8;
	   t5 = (x >> 24) & MASK8;
	   t6 = (x >> 16) & MASK8;
	   t7 = (x >>  8) & MASK8;
	   t8 =  x        & MASK8;
	*/

	t1 := sbox1(uint8(x >> 56))
	t2 := sbox2(uint8(x>>48))
	t3 := sbox3(uint8(x>>40))
	t4 := sbox4(uint8(x>>32))
	t5 := sbox2(uint8(x>>24))
	t6 := sbox3(uint8(x>>16))
	t7 := sbox4(uint8(x>>8))
	t8 := sbox1(uint8(x))

	/*
	   y1 = t1 ^ t3 ^ t4 ^ t6 ^ t7 ^ t8;
	   y2 = t1 ^ t2 ^ t4 ^ t5 ^ t7 ^ t8;
	   y3 = t1 ^ t2 ^ t3 ^ t5 ^ t6 ^ t8;
	   y4 = t2 ^ t3 ^ t4 ^ t5 ^ t6 ^ t7;
	   y5 = t1 ^ t2 ^ t6 ^ t7 ^ t8;
	   y6 = t2 ^ t3 ^ t5 ^ t7 ^ t8;
	   y7 = t3 ^ t4 ^ t5 ^ t6 ^ t8;
	   y8 = t1 ^ t4 ^ t5 ^ t6 ^ t7;
	*/
	y1 := t1 ^ t3 ^ t4 ^ t6 ^ t7 ^ t8
	y2 := t1 ^ t2 ^ t4 ^ t5 ^ t7 ^ t8
	y3 := t1 ^ t2 ^ t3 ^ t5 ^ t6 ^ t8
	y4 := t2 ^ t3 ^ t4 ^ t5 ^ t6 ^ t7
	y5 := t1 ^ t2 ^ t6 ^ t7 ^ t8
	y6 := t2 ^ t3 ^ t5 ^ t7 ^ t8
	y7 := t3 ^ t4 ^ t5 ^ t6 ^ t8
	y8 := t1 ^ t4 ^ t5 ^ t6 ^ t7

	// F_OUT = (y1 << 56) | (y2 << 48) | (y3 << 40) | (y4 << 32)| (y5 << 24) | (y6 << 16) | (y7 <<  8) | y8;
	return (y1 << 56) | (y2 << 48) | (y3 << 40) | (y4 << 32) | (y5 << 24) | (y6 << 16) | (y7 << 8) | y8
}

var sbox = [256]uint64{
	0x70, 0x82, 0x2c, 0xec, 0xb3, 0x27, 0xc0, 0xe5, 0xe4, 0x85, 0x57, 0x35, 0xea, 0x0c, 0xae, 0x41,
	0x23, 0xef, 0x6b, 0x93, 0x45, 0x19, 0xa5, 0x21, 0xed, 0x0e, 0x4f, 0x4e, 0x1d, 0x65, 0x92, 0xbd,
	0x86, 0xb8, 0xaf, 0x8f, 0x7c, 0xeb, 0x1f, 0xce, 0x3e, 0x30, 0xdc, 0x5f, 0x5e, 0xc5, 0x0b, 0x1a,
	0xa6, 0xe1, 0x39, 0xca, 0xd5, 0x47, 0x5d, 0x3d, 0xd9, 0x01, 0x5a, 0xd6, 0x51, 0x56, 0x6c, 0x4d,
	0x8b, 0x0d, 0x9a, 0x66, 0xfb, 0xcc, 0xb0, 0x2d, 0x74, 0x12, 0x2b, 0x20, 0xf0, 0xb1, 0x84, 0x99,
	0xdf, 0x4c, 0xcb, 0xc2, 0x34, 0x7e, 0x76, 0x05, 0x6d, 0xb7, 0xa9, 0x31, 0xd1, 0x17, 0x04, 0xd7,
	0x14, 0x58, 0x3a, 0x61, 0xde, 0x1b, 0x11, 0x1c, 0x32, 0x0f, 0x9c, 0x16, 0x53, 0x18, 0xf2, 0x22,
	0xfe, 0x44, 0xcf, 0xb2, 0xc3, 0xb5, 0x7a, 0x91, 0x24, 0x08, 0xe8, 0xa8, 0x60, 0xfc, 0x69, 0x50,
	0xaa, 0xd0, 0xa0, 0x7d, 0xa1, 0x89, 0x62, 0x97, 0x54, 0x5b, 0x1e, 0x95, 0xe0, 0xff, 0x64, 0xd2,
	0x10, 0xc4, 0x00, 0x48, 0xa3, 0xf7, 0x75, 0xdb, 0x8a, 0x03, 0xe6, 0xda, 0x09, 0x3f, 0xdd, 0x94,
	0x87, 0x5c, 0x83, 0x02, 0xcd, 0x4a, 0x90, 0x33, 0x73, 0x67, 0xf6, 0xf3, 0x9d, 0x7f, 0xbf, 0xe2,
	0x52, 0x9b, 0xd8, 0x26, 0xc8, 0x37, 0xc6, 0x3b, 0x81, 0x96, 0x6f, 0x4b, 0x13, 0xbe, 0x63, 0x2e,
	0xe9, 0x79, 0xa7, 0x8c, 0x9f, 0x6e, 0xbc, 0x8e, 0x29, 0xf5, 0xf9, 0xb6, 0x2f, 0xfd, 0xb4, 0x59,
	0x78, 0x98, 0x06, 0x6a, 0xe7, 0x46, 0x71, 0xba, 0xd4, 0x25, 0xab, 0x42, 0x88, 0xa2, 0x8d, 0xfa,
	0x72, 0x07, 0xb9, 0x55, 0xf8, 0xee, 0xac, 0x0a, 0x36, 0x49, 0x2a, 0x68, 0x3c, 0x38, 0xf1, 0xa4,
	0x40, 0x28, 0xd3, 0x7b, 0xbb, 0xc9, 0x43, 0xc1, 0x15, 0xe3, 0xad, 0xf4, 0x77, 0xc7, 0x80, 0x9e,
}

func sbox1(x uint8) uint64 {
	return sbox[x]
}

// SBOX2[x] = SBOX1[x] <<< 1;
func sbox2(x uint8) uint64 {
	return bits.RotateLeft64(sbox[x], 1)
}

// SBOX3[x] = SBOX1[x] <<< 7;
func sbox3(x uint8) uint64 {
	return bits.RotateLeft64(sbox[x], 7)
}

// SBOX4[x] = SBOX1[x <<< 1];
func sbox4(x uint8) uint64 {
	xx := bits.RotateLeft8(x, 1)
	return sbox[xx]
}
