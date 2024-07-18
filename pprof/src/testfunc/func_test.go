package pprof_test

import (
	"math/rand"
	"testing"
	"time"
)

var test_vec []int = []int{1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}

func TestCacheWaterCorrectness(t *testing.T) {
	CacheWaterCorrectnessTest()
	test_vec_expect := 6
	if res := CacheWater_chaidaxuan(test_vec); res != test_vec_expect {
		t.Errorf("Cache water result is not as expect, result is %v, expect is %v", res, test_vec_expect)
		t.Fail()
	}
}

func BenchmarkCacheWaterWithRandData(b *testing.B) {
	b.StopTimer()
	testTimes := 10000
	vecSize := 10000
	r := rand.New(rand.NewSource(time.Now().Unix()))
	vecs := make([][]int, 0, testTimes)
	for i := 0; i < testTimes; i++ {
		vec := make([]int, 0, vecSize)
		for j := 0; j < vecSize; j++ {
			vec = append(vec, r.Intn(10))
		}
		vecs = append(vecs, vec)
	}
	b.StartTimer()
	for i := range vecs {
		CacheWater_chaidaxuan(vecs[i])
	}
}

func BenchmarkCacheWater(b *testing.B) {
	b.StopTimer()
	// testTimes := 10000
	// vecSize := 10000
	// r := rand.New(rand.NewSource(time.Now().Unix()))
	// vecs := make([][]int, 0, testTimes)
	// for i := 0; i < testTimes; i++ {
	// 	vec := make([]int, 0, vecSize)
	// 	for j := 0; j < vecSize; j++ {
	// 		vec = append(vec, r.Intn(10))
	// 	}
	// 	vecs = append(vecs, vec)
	// }
	// 需要提供用例
	b.StartTimer()
	CacheWater_chaidaxuan(test_vec)
}
