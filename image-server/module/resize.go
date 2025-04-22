package module

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"os"

	"github.com/disintegration/imaging"
	"github.com/jdeng/goheif"
)

func Resize(srcPath string, dstPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(srcPath))

	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".heic":
		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		img, err = goheif.Decode(bytes.NewReader(data))
	default:
		return fmt.Errorf("unsupported formst : %s", ext)
	}
	if err != nil {
		return err
	}

	resized := imaging.Fit(img, 800, 800, imaging.Lanczos)
	outFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, resized, &jpeg.Options{Quality: 80})
}
