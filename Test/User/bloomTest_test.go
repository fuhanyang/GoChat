package main

import (
	"common/bloomFilter"
	"common/snowflake"
	"strconv"
	"testing"
	"time"
	"user/settings"
)

func TestBloomFilter(t *testing.T) {
	var (
		bitmapLen int64 = 1008547758
		hashCount int32 = 5
	)
	f := bloomFilter.NewLocalBloomFilter(bitmapLen, hashCount)
	err := snowflake.Init(&settings.SnowflakeConfig{StartTime: "2006-01-02T15:04:05Z"})
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 100000000; i++ {
		time.Sleep(time.Microsecond * 10)
		f.Set(strconv.FormatUint(uint64(snowflake.GetID()), 10))
	}
}
