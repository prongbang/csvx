package csvx

import "os"

// ReadByte reads the entire contents of the file with the specified filename and returns them as a slice of bytes.
// If an error occurs during the operation, a nil slice will be returned. This function can be used to read the contents
// of text files or other files that are encoded as byte streams. Note that this function may not be suitable for reading
// large files, as it reads the entire file into memory at once. For large files, consider using the os package or
// a buffered reader to read the file in smaller chunks.
func ReadByte(filename string) []byte {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
	return bytes
}
