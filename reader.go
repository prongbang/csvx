package csvx

import "os"

func ReadByte(filename string) []byte {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
	return bytes
}
