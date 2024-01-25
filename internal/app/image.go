package app

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
)

func ResizeImage(blob []byte, width, height uint) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
