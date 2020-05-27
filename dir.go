package stage

import (
	"os"
	"strconv"
)

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func numberOfDigit(i int) int {
	return len(strconv.Itoa(i))
}
