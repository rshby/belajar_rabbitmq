package helper

// method fibonaci
func Fibonaci(n int) int {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return Fibonaci(n-1) + Fibonaci(n-2)
	}
}
