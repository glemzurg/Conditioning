package conditioning

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

// RenderImage draws an image to the screen.
func RenderImage(config Config, cr *cairo.Context, displayImage DisplayImage) {

	// Paint the graphic.
	gtk.GdkCairoSetSourcePixBuf(cr, displayImage.Pixbuf, displayImage.X, displayImage.Y)
	cr.Paint()
}
