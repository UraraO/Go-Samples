package pprof_test

import "fmt"

// var test_vec2 []int = []int{1}

func CacheWater_chaidaxuan(vec []int) (res int) {
	lmax, rmax, res := 0, 0, 0
	for il, ir := 0, len(vec)-1; il < ir; {
		lmax = max(lmax, vec[il])
		rmax = max(rmax, vec[ir])
		if lmax < rmax {
			res += min(lmax, rmax) - vec[il]
			il++
		} else {
			res += min(lmax, rmax) - vec[ir]
			ir--
		}
	}
	return res
}

func CacheWaterCorrectnessTest() {
	test_vec_expect := 6
	test_vec := []int{1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
	res := CacheWater_chaidaxuan(test_vec)
	fmt.Println(res == test_vec_expect, res)
}
