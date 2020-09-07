package main

import (
	"flag"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"glemzurg/conditioning"
)

const (
	// Key strokes used.
	_KEY_SPACE      = 32
	_KEY_L          = 108
	_KEY_R          = 114
	_KEY_LEFT  uint = 65361
	// _KEY_UP    uint = 65362
	_KEY_RIGHT uint = 65363
	// _KEY_DOWN  uint = 65364

	// The default window title.
	_DEFAULT_TITLE = "Affirmations"
)

func main() {
	var err error

	var configFilename, affirmationFilename string
	flag.StringVar(&configFilename, "config", "", "configuration")
	flag.StringVar(&affirmationFilename, "affirm", "", "affirmations")
	flag.Parse()

	log.Println(`config: `, configFilename)
	log.Println(`affirmations: `, affirmationFilename)

	// The image path is the root of the affirmations.
	imagePath := filepath.Dir(affirmationFilename) + "/images/"
	log.Println(`images: `, imagePath)

	// Get the config in a useable form.
	config, err := conditioning.LoadConfig(configFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", config)

	// Random seed.
	rand.Seed(time.Now().UnixNano())

	// Prime the system.
	system, err := conditioning.NewSystem(config, affirmationFilename, imagePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n\nCommands are LEFT+RIGHT ARROWS (go back and forth in slide show), SPACE (start/stop slide show), R (toggle order of slideshow to random), L (reload from file)\n\n")

	// Load the affiramtions.
	if err = system.Load(); err != nil {
		log.Fatal(err)
	}

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	title := _DEFAULT_TITLE
	if config.Title != "" {
		title = config.Title
	}
	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	drawingArea, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Unable to create drawing area:", err)
	}

	// Add the image to the window.
	win.Add(drawingArea)

	// Add a key press event.
	win.Connect("key-press-event", func(win *gtk.Window, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}

		switch keyEvent.KeyVal() {

		case _KEY_SPACE:
			if err = system.StartStopSlideShow(win); err != nil {
				log.Printf(`key-press-event StartStopSlideShow(): %+v`, err)
			}

		case _KEY_L:
			if err = system.Load(); err != nil {
				log.Printf(`key-press-event Load(): %+v`, err)
			}
			win.QueueDraw()

		case _KEY_R:
			system.RandomOnOff()
			win.QueueDraw()

		case _KEY_LEFT:
			if err = system.Left(); err != nil {
				log.Printf(`key-press-event Left(): %+v`, err)
			}
			win.QueueDraw()

		case _KEY_RIGHT:
			if err = system.Right(); err != nil {
				log.Printf(`key-press-event Right(): %+v`, err)
			}
			win.QueueDraw()
		}
	})

	// Add a drawing event that exposes the cairo context of the drawing area.
	drawingArea.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {

		// Get the window size.
		winWidth, winHeight := win.GetSize()

		// Get the affirmation.
		affirmationIndex, displayText, displayImage, displayBoth, affirmationFound := system.DisplayTextImage()

		// Get a cached slide if there is one.
		cachedPixbuf, cacheFound := system.GetCachedSlide(affirmationIndex, winWidth, winHeight)
		if cacheFound && displayBoth {

			// Paint the cached slide.
			gtk.GdkCairoSetSourcePixBuf(cr, cachedPixbuf, 0, 0)
			cr.Paint()

		} else {

			// Paint the screen black.
			cr.SetSourceRGB(0, 0, 0)
			cr.Rectangle(0, 0, 10000, 10000) // Lots of black.
			cr.Fill()

			// Pick the shortest ratio.
			widthRatio := float64(winWidth) / float64(config.ScreenWidth)
			heightRatio := float64(winHeight) / float64(config.ScreenHeight)
			ratio := 1.0
			switch {
			case widthRatio < heightRatio:
				ratio = widthRatio
			case heightRatio < widthRatio:
				ratio = heightRatio
			}

			// Pick the offset that matters.
			var xOffset, yOffset float64
			switch {
			case widthRatio < heightRatio:
				// We're centering on the y axis.
				yOffset = (float64(winHeight) - float64(config.ScreenHeight)*ratio) / 2.0
			case heightRatio < widthRatio:
				// We're centering on the x axis.
				xOffset = (float64(winWidth) - float64(config.ScreenWidth)*ratio) / 2.0
			}

			// Create a matrix that represents this transform.
			matrix := cairo.NewMatrix(ratio, 0.0, 0.0, ratio, xOffset, yOffset)
			cr.Transform(matrix)

			// If we have an affirmation render it.
			if affirmationFound {
				if displayImage != nil {
					conditioning.RenderImage(config, cr, *displayImage)
				}
				if displayBoth {
					conditioning.RenderAffirmation(config, cr, displayText)
				}
			}

			// Is the the whole slide?
			if displayBoth {
				// Attempt to cache the window image.
				winGdk, err := win.GetWindow()
				if err != nil {
					log.Printf("win.GetWindow() err: %+v", err)
				} else {
					pixbuf, err := winGdk.PixbufGetFromWindow(0, 0, winWidth, winHeight)
					if err != nil {
						log.Printf("winGdk.PixbufGetFromWindow() err: %+v", err)
					} else {
						// No error, we can cache this.
						system.CacheSlide(affirmationIndex, winWidth, winHeight, pixbuf)
					}
				}
			}
		}
	})

	win.QueueDraw()

	// Set the default window size.
	win.SetDefaultSize(int(config.ScreenWidth), int(config.ScreenHeight))

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
