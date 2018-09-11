package radix

func bit(i uint8, c byte) uint8 {
	if 1<<(i-1)&c > 0 {
		return 1
	}
	return 0
}
