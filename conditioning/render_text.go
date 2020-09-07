package conditioning

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
)

const (
	// The color.
	BLACK = true
	WHITE = false

	// The outline.
	OUTLINE    = true
	NO_OUTLINE = false
)

// RenderAffirmation writes text to the screen.
func RenderAffirmation(config Config, cr *cairo.Context, displayText DisplayText) {
	renderAffirmation(config, cr, displayText.X, displayText.Y, displayText.PangoMarkup, displayText.FontDescription, displayText.Black, displayText.Outline)
}

// renderAffirmation writes text to the screen.
func renderAffirmation(config Config, cr *cairo.Context, x, y int, pangoMarkup string, fontDescription *pango.FontDescription, black, outline bool) {

	// If outline, draw it first.
	if outline {

		// Create a pango layout.
		layout := pango.CairoCreateLayout(cr)

		// Position at the beginning of the text.
		cr.MoveTo(float64(x), float64(y))

		// Black or white? Do opposite of text.
		var outlineScale float64
		if black {
			cr.SetSourceRGB(255, 255, 255)
			outlineScale = config.BlackOutlineScale
		} else {
			cr.SetSourceRGB(0, 0, 0)
			outlineScale = config.WhiteOutlineScale
		}

		// Set the font description.
		layout.SetFontDescription(fontDescription)

		// Set the markup in the mask.
		layout.SetMarkup(pangoMarkup, -1)

		// Half of this stroke will be the outline.
		strokeWidth := (float64(fontDescription.GetSize()) / pango.PANGO_SCALE) * outlineScale
		cr.SetLineWidth(strokeWidth)

		// Create the mask and outline the text.
		pango.CairoLayoutPath(cr, layout)
		cr.Stroke()
	}

	// Create a pango layout.
	layout := pango.CairoCreateLayout(cr)

	// Position at the beginning of the text.
	cr.MoveTo(float64(x), float64(y))

	// Black or white?
	if black {
		cr.SetSourceRGB(0, 0, 0)
	} else {
		cr.SetSourceRGB(255, 255, 255)
	}

	// Set the font description.
	layout.SetFontDescription(fontDescription)

	// Set the markup in the mask.
	layout.SetMarkup(pangoMarkup, -1)

	// Create the mask and fill the text.
	pango.CairoShowLayout(cr, layout)
	cr.Fill()
}
