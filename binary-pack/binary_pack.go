package binary_pack

type BinaryPack struct {}

func (bp *BinaryPack) Pack(format string, args ...interface{}) []byte {
	return []byte{}
}

func (bp *BinaryPack) UnPack(format string, msg []byte) []interface{} {
	return make([]interface{}, 1)
}

func (bp *BinaryPack) CalcSize(format string) int {
	return 0
}
