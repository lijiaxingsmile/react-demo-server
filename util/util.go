package util

import "os"

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Response(success bool, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}
}
