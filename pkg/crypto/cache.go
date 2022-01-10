package crypto

import "image"

type Cache interface {
	LoadImage(url string) (image.Image, error)
	SaveImage(url string, img image.Image)
}
