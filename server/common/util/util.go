package util

///////////////////////////////////////////////////////
// uTil start
///////////////////////////////////////////////////////
// cutOver49把uint64对象的前14位清零（也就是保留0~49位）
func CutOver49(fullNum uint64) uint64{
	clearTag := uint64(18445618173802708992)	/* = 1111111111111100000000000000000000000000000000000000000000000000（14个1，50个0） */
	clearedNum := fullNum & (^clearTag)
	return clearedNum
}

// countUint64求解uint64对象中1的个数
func CountUint64(x uint64) int {
	var c int
	for i:=0; i<64; i++ {
		c += int(x&1)
		x = x>>1
	}
	return c
}

// FindLowBit求解uint64对象中第一个为1的索引
func FindLowBit(x uint64) int {
	for i:=0; i<64; i++ {
		if int(x&1) > 0 {
			return i
		}
		x = x>>1
	}
	return -1	// x中不存在为1的位！
}
///////////////////////////////////////////////////////
// uTil end
///////////////////////////////////////////////////////