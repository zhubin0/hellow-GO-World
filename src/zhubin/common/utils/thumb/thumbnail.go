package thumb

import (
	_ "image/jpeg"
	"io"

	_ "image/gif"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"

	"bytes"

	"image"

	"datamesh.com/MeshExpert/global"
	"github.com/nfnt/resize"
	"image/png"
)

// Generate thumbnail.
// NOTE for decoding and encoding in particular format, you need to import the decoder for side effect.
// NOTE call is responsible for closing the reader.
func GenThumbnail(image io.Reader, out *bytes.Buffer) error {
	th := thumbnail{original: image, thumb: out}
	return th.generate()
}

type thumbnail struct {
	original io.Reader
	thumb    *bytes.Buffer
}

func (t *thumbnail) generate() error {
	// decode jpeg into image.Image
	img, _, err := image.Decode(t.original)
	if err != nil {
		return err
	}
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(global.IMG_PROFILE_THUMBNAIL_WIDTH, global.IMG_PROFILE_THUMBNAIL_HEIGHT, img, resize.Lanczos3)
	// write new image
	return png.Encode(t.thumb, m)
}
