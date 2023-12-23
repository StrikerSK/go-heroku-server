package utils

import "fmt"

// GetFileSize - Resolve ideal file size up to MegaBytes
func GetFileSize(fileSize int64) (outputSize string) {
	switch {
	case fileSize < 1024:
		outputSize = fmt.Sprintf("%d B", fileSize)
		break
	case fileSize < 1048576:
		fileSize = fileSize / 1024
		outputSize = fmt.Sprintf("%d kB", fileSize)
		break
	default:
		fileSize = fileSize / 1048576
		outputSize = fmt.Sprintf("%d MB", fileSize)
		break
	}
	return
}
