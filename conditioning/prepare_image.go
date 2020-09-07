package conditioning

import (
	"fmt"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
)

const (
	_PRESERVE_ASPECT_RATIO = true
)

// DisplayImage is everything neded to display image on the screen.
type DisplayImage struct {
	// Coordinate.
	X float64 // Coordinate on the screen.
	Y float64 // Coordinate on the screen.
	// Image.
	Filename string
	Pixbuf   *gdk.Pixbuf
}

// PrepareImage prepares a image for display on the screen.
func PrepareImage(config Config, cr *cairo.Context, imagePath string, affirmationImage AffirmationImage) (displayImage DisplayImage, err error) {

	// The full filename.
	displayImage.Filename = imagePath + affirmationImage.Filename

	// Load the data of an image file.
	pixbuf, err := gdk.PixbufNewFromFileAtScale(displayImage.Filename, int(config.ScreenWidth), int(config.ScreenHeight), _PRESERVE_ASPECT_RATIO)
	if err != nil {
		fmt.Printf("%T\n", err)
		return DisplayImage{}, Error(err)
	}

	// Are we changing the size?
	if affirmationImage.Scale != 1.0 {

		// What are the dimensions of the image.
		origWidth := pixbuf.GetWidth()
		origHeight := pixbuf.GetHeight()

		// Compute the new width and hight.
		destWidth := int(float64(origWidth) * affirmationImage.Scale)
		destHeight := int(float64(origHeight) * affirmationImage.Scale)

		// Resize.
		// The interp is defined here: http://openbooks.sourceforge.net/books/wga/graphics-gdk-pixbuf.html
		// Lowest to highest quality (fastest to slowest speed):
		// INTERP_NEAREST
		// INTERP_TILES
		// INTERP_BILINEAR
		// INTERP_HYPER
		scaledPixbuf, err := pixbuf.ScaleSimple(destWidth, destHeight, gdk.INTERP_HYPER)
		if err != nil {
			return DisplayImage{}, Error(err)
		}
		pixbuf = scaledPixbuf
	}

	// What are the dimensions of the image.
	imageWidth := pixbuf.GetWidth()
	imageHeight := pixbuf.GetHeight()

	// What is the center point of the screen?
	centerX := int(config.ScreenWidth) / 2
	centerY := int(config.ScreenHeight) / 2

	// If we move back by half width and height we get centered image.
	centeredX := centerX - imageWidth/2
	centeredY := centerY - imageHeight/2

	// Set the position of this text using the offsets.
	displayImage.X = float64(centeredX + affirmationImage.OffsetX)
	displayImage.Y = float64(centeredY + affirmationImage.OffsetY)

	// Attach the image data itself.
	displayImage.Pixbuf = pixbuf

	return displayImage, nil
}
