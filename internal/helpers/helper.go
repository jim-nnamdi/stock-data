package helpers

func copyslice(a []int, b []int) []int {
	out := make([]int, 0)
	copy(out, a)
	copy(out[len(a):], b)
	return out
}

func dot(sp []float64, sp2 []float64) float64 {
	total := 0.0
	for i, v := range sp {
		total += v * sp[i]
	}
	return total
}
