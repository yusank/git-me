package utils

func Round(f float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int((f*pow)+0.5)) / pow
}

// f=123.3349995 return 123.34
func RoundSpec(f float64, decimals int) float64 {
	return Round(Round(f, decimals+1), decimals)
}
