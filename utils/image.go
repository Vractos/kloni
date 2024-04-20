package utils

import (
	"image"

	"github.com/disintegration/imaging"
)

func ResizeImage(image image.Image, width, height int) image.Image {
	resized := imaging.Resize(image, width, height, imaging.Lanczos)

	return resized
}
