package conditioning

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config is the configuration for the system.
type Config struct {

	// The logic.
	SleepMilli uint // How long to show each slide.

	// The screen.
	ScreenWidth  uint // The basic screen width.
	ScreenHeight uint // The basic screen height.

	// The font.
	FontFace          string  // The default font face.
	FontSize          uint    // The default font size.
	WhiteOutlineScale float64 // For white text. The 0.0-1.0 % of the font size for the outline (only half will show).
	BlackOutlineScale float64 // For black text. The 0.0-1.0 % of the font size for the outline (only half will show).
}

// Validate the config is well-formed.
func (c *Config) Validate() (err error) {
	if c.SleepMilli <= 0 {
		return Errorf(`invalid SleepMilli: %d`, c.SleepMilli)
	}
	if c.ScreenWidth <= 0 {
		return Errorf(`invalid ScreenWidth: %d`, c.ScreenWidth)
	}
	if c.ScreenHeight <= 0 {
		return Errorf(`invalid ScreenHeight: %d`, c.ScreenHeight)
	}
	if c.FontFace == "" {
		return Errorf(`invalid FontFace: '%s'`, c.FontFace)
	}
	if c.FontSize <= 0 {
		return Errorf(`invalid FontSize: %d`, c.FontSize)
	}
	if c.BlackOutlineScale <= 0 {
		return Errorf(`invalid BlackOutlineScale: %+v`, c.BlackOutlineScale)
	}
	if c.WhiteOutlineScale <= 0 {
		return Errorf(`invalid WhiteOutlineScale: %+v`, c.WhiteOutlineScale)
	}
	return nil
}

// LoadConfig loads the config file a filename.
func LoadConfig(configFilename string) (config Config, err error) {

	// Open the file.
	file, err := os.Open(configFilename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Load the bytes from the file.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	// Parse the
	if err = json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}

	// The populated config.
	return config, nil
}
