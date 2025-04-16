package main

import (
	"fmt"
	"os"
)

func ensureDBDirectoryExists(path string, perm os.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, perm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		err = os.Chmod(path, perm)
		if err != nil {
			return fmt.Errorf("failed to change directory permissions: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check directory existence: %v", err)
	}
	return nil
}
