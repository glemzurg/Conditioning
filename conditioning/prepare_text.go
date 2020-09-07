package conditioning

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
	"strconv"
)

// DisplayText is everything neded to display text on the screen.
type DisplayText struct {
	// Coordinate.
	X int // Coordinate on the screen.
	Y int // Coordiante on the screen.
	// Text.
	PangoMarkup     string                 // The text to display with optional formatting.
	FontDescription *pango.FontDescription // The font (and font size to use).
	// The color and outline.
	Black   bool // If true display black, otherwise display white.
	Outline bool // If true, put an inversed outline underneath the text.
}

// PrepareText prepares a text for display on the screen.
func PrepareText(config Config, cr *cairo.Context, message string, textProperties TextProperties) (displayText DisplayText) {

	// Create markup.
	displayText.PangoMarkup = PangoMarkup(message)

	// Font information.
	fontFace := config.FontFace
	fontSize := config.FontSize
	if textProperties.FontSize != 0 {
		fontSize = textProperties.FontSize
	}
	fontString := fontFace + " " + strconv.Itoa(int(fontSize))

	displayText.FontDescription = pango.FontDescriptionFromString(fontString)

	// Color and outline.
	displayText.Black = textProperties.Black
	displayText.Outline = OUTLINE // Always outline.

	// Compute the absolute position of the text.

	// Create a pango layout.
	layout := pango.CairoCreateLayout(cr)

	// Set the font description.
	layout.SetFontDescription(displayText.FontDescription)

	// Set the markup in the mask.
	layout.SetMarkup(displayText.PangoMarkup, -1)

	// What are the dimensions of the layout.
	pangoWidth, pangoHeight := layout.GetSize()
	textWidth := pangoWidth / pango.PANGO_SCALE
	textHeight := pangoHeight / pango.PANGO_SCALE

	// What is the center point of the screen?
	centerX := int(config.ScreenWidth) / 2
	centerY := int(config.ScreenHeight) / 2

	// If we move back by half width and height we get centered text.
	centeredX := centerX - textWidth/2
	centeredY := centerY - textHeight/2

	// Set the position of this text using the offsets.
	displayText.X = centeredX + textProperties.OffsetX
	displayText.Y = centeredY + textProperties.OffsetY

	return displayText
}
