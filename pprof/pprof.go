package main

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

func main() {

}
