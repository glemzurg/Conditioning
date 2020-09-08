package conditioning

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	// The default slide show title.
	_DEFAULT_TITLE = "Affirmations"
)

// Affirmation is a single affirmation.
type Affirmation struct {
	Message string           // The text of the affirmation.
	Image   AffirmationImage // The image associated with the affirmation.
	Text    TextProperties   // Details about how to display the text.
}

// AffrimationImage is the image details of the affirmation.
type AffirmationImage struct {
	Filename string  // The filename for the image.
	OffsetX  int     // Offset from center.
	OffsetY  int     // Offset from center.
	Scale    float64 // How much to grow/shrink the image (1.0 is original size).
}

// AffrimationText is the text/font details of the affirmation.
type TextProperties struct {
	Black    bool // The color of the text.
	OffsetX  int  // Offset from center.
	OffsetY  int  // Offset from center.
	FontSize uint // The font size of the text.
}

// parseAffirmations parses the affirmation text.
func parseAffirmations(unparsed string) (affirmations []Affirmation, title string) {

	// Split the text on newlines.
	lines := strings.Split(unparsed, "\n")
	parsedNonBlankLine := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {

		case strings.HasPrefix(line, "//"):
			// This is a comment.
			if !parsedNonBlankLine {
				// The comment before all else, this is the title.
				title = strings.TrimSpace(strings.TrimLeft(line, "/"))
			}

			// We have parsed a line (even a comment)
			parsedNonBlankLine = true

		case line == "":
			// This is a blank line.

		default:
			// Not a comment? Not a blank line? This is an affirmation.

			// Split out the affirmation mesage from the display details.
			lineParts := strings.Split(line, "[")
			message := strings.TrimSpace(lineParts[0])

			// Get the image details.
			var image AffirmationImage
			var text TextProperties
			if len(lineParts) > 1 {

				// Split the display.
				displayParts := strings.Split(strings.Trim(lineParts[1], " ]"), " ")
				for i := 0; i < len(displayParts); i++ {
					part := displayParts[i]

					// Parse an image.
					parsedImage, parsed := parseImage(part)
					if parsed {
						image = parsedImage
					}

					// Parse a text display.
					parsedText, parsed := parseText(part)
					if parsed {
						text = parsedText
					}
				}
			}

			affirmations = append(affirmations, Affirmation{
				Message: message,
				Image:   image,
				Text:    text,
			})

			// We have parsed a line.
			parsedNonBlankLine = true
		}
	}

	return affirmations, title
}

// parseImage parses the part of an affirmation that
func parseImage(text string) (image AffirmationImage, parsed bool) {
	textParts := strings.Split(text, ":")

	// Is there a file extension?
	filename := textParts[0]
	if strings.Index(filename, ".") == -1 {
		return AffirmationImage{}, false // Not a filename.
	}

	// Examine each other part of the display.
	var scale float64 = 1.0
	var offsetX, offsetY int
	for i := 1; i < len(textParts); i++ {
		part := textParts[i]
		switch {

		// Is this coordinates?
		case strings.Index(part, ",") >= 0:
			coordinateParts := strings.Split(part, ",")
			x, xErr := strconv.Atoi(coordinateParts[0])
			y, yErr := strconv.Atoi(coordinateParts[1])
			if xErr == nil && yErr == nil {
				offsetX = x
				offsetY = y
			} else {
				if xErr != nil {
					log.Printf("%+v\n", xErr)
				}
				if yErr != nil {
					log.Printf("%+v\n", yErr)
				}
			}

			// Attempt to parse a scale.
		default:

			value, err := strconv.ParseFloat(part, 64)
			if err == nil {
				scale = value
			} else {
				if err != nil {
					log.Printf("%+v\n", err)
				}
			}
		}
	}

	return AffirmationImage{
		Filename: filename,
		OffsetX:  offsetX,
		OffsetY:  offsetY,
		Scale:    scale,
	}, true
}

// parseText parses the part of an affirmation that
func parseText(text string) (textDetails TextProperties, parsed bool) {
	textParts := strings.Split(text, ":")

	// At the beginning, we need to know if this is black or white text.
	color := textParts[0]
	var black bool
	switch color {
	case "b":
		black = true
	case "w":
		black = false
	default:
		return TextProperties{}, false // Not a color.
	}

	// Examine each other part of the display.
	var fontSize, offsetX, offsetY int
	for i := 1; i < len(textParts); i++ {
		part := textParts[i]
		switch {

		// Is this coordinates?
		case strings.Index(part, ",") >= 0:
			coordinateParts := strings.Split(part, ",")
			x, xErr := strconv.Atoi(coordinateParts[0])
			y, yErr := strconv.Atoi(coordinateParts[1])
			if xErr == nil && yErr == nil {
				offsetX = x
				offsetY = y
			} else {
				if xErr != nil {
					log.Printf("%+v\n", xErr)
				}
				if yErr != nil {
					log.Printf("%+v\n", yErr)
				}
			}

			// Attempt to parse a font size.
		default:

			value, err := strconv.Atoi(part)
			if err == nil {
				fontSize = value
			} else {
				if err != nil {
					log.Printf("%+v\n", err)
				}
			}
		}
	}

	return TextProperties{
		Black:    black,
		OffsetX:  offsetX,
		OffsetY:  offsetY,
		FontSize: uint(fontSize),
	}, true
}

// LoadAffirmations loads the affirmations from the affirmations file.
func LoadAffirmations(affirmationFilename string) (affirmations []Affirmation, title string, err error) {

	// Open the file.
	file, err := os.Open(affirmationFilename)
	if err != nil {
		return nil, "", Error(err)
	}
	defer file.Close()

	// Load the bytes from the file.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", Error(err)
	}

	// Extract the affirmations.
	affirmations, title = parseAffirmations(string(bytes))
	if title == "" {
		title = _DEFAULT_TITLE
	}

	return affirmations, title, nil
}
