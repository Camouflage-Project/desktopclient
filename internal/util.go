package internal

import (
	"path/filepath"
)

func getFilenameFromProcessName(processName string) string {
	return filepath.Base(processName)
}
