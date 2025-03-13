package helper

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func Hash(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ScanPattern(path, pattern string, recursive bool) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.Contains(path, pattern) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func ExtractPath(path string) string {
	s := strings.Split(path, "/")
	return strings.Join(s[:len(s)-1], "/")
}
