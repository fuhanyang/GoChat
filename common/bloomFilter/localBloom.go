package bloomFilter

import (
	"github.com/demdxx/gocast"
	"github.com/spaolacci/murmur3"
	"math"
)

type BloomFilter interface {
	Exists(origin string) bool
	Set(origin string)
}

type Encryptor struct {
}

func NewEncryptor() *Encryptor {
	return &Encryptor{}
}
func (e *Encryptor) Encrypt(origin string) int32 {
	hasher := murmur3.New32()
	_, _ = hasher.Write([]byte(origin))
	return int32(hasher.Sum32() % math.MaxInt32)
}

type LocalBloomFilter struct {
	m         int64
	k, n      int32
	bitmap    []int32
	encryptor *Encryptor
}

func NewLocalBloomFilter(m int64, k int32) *LocalBloomFilter {
	return &LocalBloomFilter{
		m:         m,
		k:         k,
		bitmap:    make([]int32, m/32+1),
		encryptor: NewEncryptor(),
	}
}

func (l *LocalBloomFilter) Exists(origin string) bool {
	for _, offset := range l.getKEncrypted(origin) {
		index := (int64(offset) % l.m) >> 5
		bitOffset := index & 31
		if l.bitmap[index]&(1<<bitOffset) == 0 {
			return false
		}

	}
	return true
}
func (l *LocalBloomFilter) getKEncrypted(val string) []int32 {
	encrypteds := make([]int32, 0, l.k)
	for i := 0; int32(i) < l.k; i++ {
		encrypted := l.encryptor.Encrypt(val)
		encrypteds = append(encrypteds, encrypted)
		if int32(i) == l.k-1 {
			break
		}
		val = gocast.ToString(encrypted)
	}
	return encrypteds
}

func (l *LocalBloomFilter) Set(origin string) {
	l.n++
	for _, offset := range l.getKEncrypted(origin) {
		index := (int64(offset) % l.m) >> 5
		bitOffset := index & 31

		l.bitmap[index] |= 1 << bitOffset
	}
}
