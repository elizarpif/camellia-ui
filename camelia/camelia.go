package camelia

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"math/bits"
	"strconv"
)

const (
	MASK8   = 0xff
	MASK32  = 0xffffffff
	MASK64  = 0xffffffffffffffff
	MASK128 = 0xffffffffffffffffffffffffffffffff

	C1 = 0xA09E667F3BCC908B
	C2 = 0xB67AE8584CAA73B2
	C3 = 0xC6EF372FE94F82BE
	C4 = 0x54FF53A5F1D36F1C
	C5 = 0x10E527FADE682D1D
	C6 = 0xB05688C2B3E6C1FD

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

func (c cameliaCipher) Print() {
	for i, k := range c.k {
		fmt.Printf("k[%d] = %x\n", i+1, k)
	}
	for i, k := range c.ke {
		fmt.Printf("ke[%d]%x\n", i+1, k)
	}
	for i, k := range c.kw {
		fmt.Printf("kw[%d]%x\n", i+1, k)
	}
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

func crypt(c cameliaCipher, dst, src []byte) {
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

	d1 ^= c.kw[0]
	d2 ^= c.kw[1]

	d2 = d2 ^ f(d1, c.k[0])
	d1 = d1 ^ f(d2, c.k[1])
	d2 = d2 ^ f(d1, c.k[2])
	d1 = d1 ^ f(d2, c.k[3])
	d2 = d2 ^ f(d1, c.k[4])
	d1 = d1 ^ f(d2, c.k[5])

	d1 = fl(d1, c.ke[0])
	d2 = flinv(d2, c.ke[1])

	d2 = d2 ^ f(d1, c.k[6])
	d1 = d1 ^ f(d2, c.k[7])
	d2 = d2 ^ f(d1, c.k[8])
	d1 = d1 ^ f(d2, c.k[9])
	d2 = d2 ^ f(d1, c.k[10])
	d1 = d1 ^ f(d2, c.k[11])

	d1 = fl(d1, c.ke[2])
	d2 = flinv(d2, c.ke[3])

	d2 = d2 ^ f(d1, c.k[12])
	d1 = d1 ^ f(d2, c.k[13])
	d2 = d2 ^ f(d1, c.k[14])
	d1 = d1 ^ f(d2, c.k[15])
	d2 = d2 ^ f(d1, c.k[16])
	d1 = d1 ^ f(d2, c.k[17])

	if c.klen > 16 {
		// 24 or 32

		d1 = fl(d1, c.ke[4])
		d2 = flinv(d2, c.ke[5])

		d2 = d2 ^ f(d1, c.k[18])
		d1 = d1 ^ f(d2, c.k[19])
		d2 = d2 ^ f(d1, c.k[20])
		d1 = d1 ^ f(d2, c.k[21])
		d2 = d2 ^ f(d1, c.k[22])
		d1 = d1 ^ f(d2, c.k[23])
	}

	d2 = d2 ^ c.kw[2]
	d1 = d1 ^ c.kw[3]

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
func rotate128Key(k [2]uint64, rot uint) (hi uint64, lo uint64) {
	//if n > 64 {
	//	n %= 64
	//	// поворот порядка
	//	k[0], k[1] = k[1], k[0]
	//}
	////t := *r0 >> (32 - n)
	////*r0 = (*r0 << n) | (*r1 >> (32 - n))
	////*r1 = (*r1 << n) | (*r2 >> (32 - n))
	////*r2 = (*r2 << n) | (*r3 >> (32 - n))
	////*r3 = (*r3 << n) | t
	//
	//// сдвигаем
	//t := k[0] >> (64 - n)
	//r := (k[0] << n) | (k[1] >> (64 - n))
	//l := (k[1] << n) | t
	//
	//return r, l

	if rot > 64 {
		rot -= 64
		k[0], k[1] = k[1], k[0]
	}

	t := k[0] >> (64 - rot)
	hi = (k[0] << rot) | (k[1] >> (64 - rot))
	lo = (k[1] << rot) | t
	return hi, lo
}

// helpKeys128 - генерация вспомогательных ключей для ключа размером 128 бит
func (c *cameliaCipher) helpKeys128(ka, kl [2]uint64) {
	c.kw[0], c.kw[1] = rotate128Key(kl, 0)

	c.k[0], c.k[1] = rotate128Key(ka, 0)
	fmt.Printf("k[1] = %x, k[2] = %x\n", c.k[0], c.k[1])
	c.k[2], c.k[3] = rotate128Key(kl, 15)
	fmt.Printf("k[3] = %x, k[4] = %x\n", c.k[2], c.k[3])
	c.k[4], c.k[5] = rotate128Key(ka, 15)

	c.ke[0], c.ke[1] = rotate128Key(ka, 30)

	c.k[6], c.k[7] = rotate128Key(kl, 45)
	c.k[8], _ = rotate128Key(ka, 45)
	_, c.k[9] = rotate128Key(kl, 60)
	c.k[10], c.k[11] = rotate128Key(ka, 60)

	c.ke[2], c.ke[3] = rotate128Key(kl, 77)

	c.k[12], c.k[13] = rotate128Key(kl, 94)
	c.k[14], c.k[15] = rotate128Key(ka, 94)
	c.k[16], c.k[17] = rotate128Key(kl, 111)

	c.kw[2], c.kw[3] = rotate128Key(ka, 111)
}

// helpKeys256 - генерация вспомогательных ключей для ключа размером 194 or 256
func (c *cameliaCipher) helpKeys256(ka, kb, kl, kr [2]uint64) {
	c.kw[0], c.kw[1] = rotate128Key(kl, 0)

	c.k[0], c.k[1] = rotate128Key(kb, 0)
	c.k[2], c.k[3] = rotate128Key(kr, 15)
	c.k[4], c.k[5] = rotate128Key(ka, 15)

	c.ke[0], c.ke[1] = rotate128Key(kr, 30)

	c.k[6], c.k[7] = rotate128Key(kb, 30)
	c.k[8], c.k[9] = rotate128Key(kl, 45)
	c.k[10], c.k[11] = rotate128Key(ka, 45)

	c.ke[2], c.ke[3] = rotate128Key(kl, 60)

	c.k[12], c.k[13] = rotate128Key(kr, 60)
	c.k[14], c.k[15] = rotate128Key(kb, 60)
	c.k[16], c.k[17] = rotate128Key(kl, 77)

	c.ke[4], c.ke[5] = rotate128Key(ka, 77)

	c.k[18], c.k[19] = rotate128Key(kr, 94)
	c.k[20], c.k[21] = rotate128Key(ka, 94)
	c.k[22], c.k[23] = rotate128Key(kl, 111)

	c.kw[2], c.kw[3] = rotate128Key(kb, 111)
}

// kAkB - вычисления вспомогательных 128-битных чисел KA, KB
func kAkB(kl, kr [2]uint64) (ka [2]uint64, kb [2]uint64) {
	var d1, d2 uint64

	//  KB = (D1 << 64) | D2;

	//   D1 = (KL ^ KR) >> 64;
	//  D2 = (KL ^ KR) & MASK64;
	d1 = kl[0] ^ kr[0]
	d2 = kl[1] ^ kr[1]

	//  D2 = D2 ^ F(D1, C1);
	//  D1 = D1 ^ F(D2, C2);
	d2 = d2 ^ f(d1, C1)
	d1 = d1 ^ f(d2, C2)

	//  D1 = D1 ^ (KL >> 64);
	//  D2 = D2 ^ (KL & MASK64);
	//  D2 = D2 ^ F(D1, C3);
	d1 = d1 ^ (kl[0])
	d2 = d2 ^ (kl[1])
	d2 = d2 ^ f(d1, C3)

	//  D1 = D1 ^ F(D2, C4);
	//  KA = (D1 << 64) | D2;
	//  D1 = (KA ^ KR) >> 64;
	//  D2 = (KA ^ KR) & MASK64;
	//  D2 = D2 ^ F(D1, C5);
	//  D1 = D1 ^ F(D2, C6);
	d1 = d1 ^ f(d2, C4)
	ka[0] = d1
	ka[1] = d2
	d1 = ka[0] ^ kr[0]
	d2 = ka[1] ^ kr[1]
	d2 = d2 ^ f(d1, C5)
	d1 = d1 ^ f(d2, C6)

	kb[0] = d1
	kb[1] = d2

	fmt.Printf("K = %x %x\n", kl[0], kl[1])
	fmt.Printf("KB = %x %x\n", kr[0], kr[1])

	fmt.Printf("KA = %x %x\n", ka[0], ka[1])
	fmt.Printf("KB = %x %x\n", kb[0], kb[1])
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

	kl[0] = binary.BigEndian.Uint64(key[0:8])
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

	c.Print()
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
	t2 := sbox2(uint8(x >> 48))
	t3 := sbox3(uint8(x >> 40))
	t4 := sbox4(uint8(x >> 32))
	t5 := sbox2(uint8(x >> 24))
	t6 := sbox3(uint8(x >> 16))
	t7 := sbox4(uint8(x >> 8))
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
	return uint64(y1)<<56 | uint64(y2)<<48 | uint64(y3)<<40 | uint64(y4)<<32 | uint64(y5)<<24 | uint64(y6)<<16 | uint64(y7)<<8 | uint64(y8)
}
