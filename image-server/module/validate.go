package module

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ValdateRequestPath(basePath string, request string) error {

	if request == "" {
		return errors.New("value is empty")
	}
	if len(request) > 100 {
		return errors.New("value is too long")
	}
	invalidChars := regexp.MustCompile(`[<>:"\\|?*\x00-\x1F]`)
	if invalidChars.MatchString(request) {
		return errors.New("value contain invalid character")
	}
	log.Println(request)
	fullPath := filepath.Join(basePath, request)
	log.Println(fullPath)
	cleanPath := filepath.Clean(fullPath)
	log.Println(cleanPath)

	basePathClean := filepath.Clean(basePath) + string(os.PathSeparator)
	log.Println(basePath)
	if !strings.HasPrefix(cleanPath, basePathClean) {
		return errors.New("value contains escape")
	}

	return nil
}
