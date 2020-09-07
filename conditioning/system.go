package conditioning

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type affirmationData struct {
	affirmation  Affirmation   // The affirmation loaded from a file.
	displayText  DisplayText   // The text prepared for rendering.
	displayImage *DisplayImage // The image prepared for rendering, if there is one.
}

// System is the wrapper for data.
type System struct {
	// Mutex bookkeeping.
	mux *sync.Mutex
	// Display.
	activeAffirmationIndex int  // The currently active affirmation.
	displayBoth            bool // If false, display image only.
	// Slide show.
	slideShowRandom   bool // If true randomly scramble the order of slides.
	slideShowTicker   *time.Ticker
	slideShowDoneChan chan bool
	slideShowIndexes  []int
	slideShowI        int
	// Cached slides.
	cachedWidth  int
	cachedHeight int
	cachedSlides map[int]*gdk.Pixbuf
	// Data.
	config              Config
	affirmationFilename string
	imagePath           string
	affirmations        []affirmationData
}

// NewSystem creates a wellformed system for displaying
func NewSystem(config Config, affirmationFilename, imagePath string) (system *System, err error) {
	if err = config.Validate(); err != nil {
		return nil, err
	}
	return &System{
		mux:                 &sync.Mutex{},
		config:              config,
		affirmationFilename: affirmationFilename,
		imagePath:           imagePath,
	}, nil
}

// RandomOnOff configures whether the slide show is ordered or random.
func (s *System) RandomOnOff() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.slideShowRandom = !s.slideShowRandom
}

// getDisplayText returns whether the text should be rendered.
func (s *System) getDisplayBoth() (displayBoth bool) {

	// If slide show is running, always display the text.
	if s.slideShowTicker != nil {
		return true // Always display text in a slide show.
	}

	return s.displayBoth // Manual navigation may or may not display text.
}

// Load loads all the affirmations and prepares them for display.
func (s *System) Load() (err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Load from the text file.
	affirmations, err := LoadAffirmations(s.affirmationFilename)
	if err != nil {
		return Error(err)
	}

	// Create a context of the proper dimensions for sizing everything.
	surface := cairo.CreateImageSurface(cairo.FORMAT_ARGB32, int(s.config.ScreenWidth), int(s.config.ScreenHeight))
	cr := cairo.Create(surface)

	// Prepare the affirmation data.
	var affirmationDatas []affirmationData
	for _, affirmation := range affirmations {

		// Prep the parts that must exist.
		data := affirmationData{
			affirmation: affirmation,
			displayText: PrepareText(s.config, cr, affirmation.Message, affirmation.Text),
		}

		// Is there an image?
		if affirmation.Image.Filename != "" {
			displayImage, err := PrepareImage(s.config, cr, s.imagePath, affirmation.Image)
			if err != nil {
				return Error(err)
			}
			data.displayImage = &displayImage
		}

		// Add the affirmation to the affirmations in the system.
		affirmationDatas = append(affirmationDatas, data)
	}
	s.affirmations = affirmationDatas

	// Keep the index in bounds.
	if s.activeAffirmationIndex > s.maxAffirmationIndex() {
		s.activeAffirmationIndex = s.maxAffirmationIndex()
	}

	// Create an index list for these affirmations we can shuffle.
	s.slideShowIndexes = s.calculateSlideShowIndixes()
	s.slideShowI = 0

	// Reset the rendered slide cache.
	s.clearSlideCacheIfNecessary(0, 0) // Passing 0, 0 should trigger a cache clear.

	return nil
}

// StartStopSlideShow starts a slide show or stops a running one.
func (s *System) StartStopSlideShow(win *gtk.Window) (err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// If there is no allocated slide show ticker, we need to start the slide show.
	if s.slideShowTicker == nil {

		s.slideShowTicker = time.NewTicker(time.Duration(s.config.SleepMilli) * time.Millisecond)
		s.slideShowDoneChan = make(chan bool)

		// Start go routine that operates the slide show.
		go func() {
			for { // Infinite loop.
				select {

				// A tick from teh ticker.
				case <-s.slideShowTicker.C:
					s.Random() // Pick a new slide.
					win.QueueDraw()

				// A stop command.
				case <-s.slideShowDoneChan:
					return // kill the goroutine.
				}
			}
		}()

	} else {
		// There is an allocated slide show ticker, we need to end the slide show.
		s.slideShowTicker.Stop()    // Stop ticker.
		s.slideShowDoneChan <- true // Escape the golang function responding to ticks.

		s.slideShowTicker = nil
		s.slideShowDoneChan = nil
	}

	return nil
}

// Random picks a random affirmation and makes it active.
func (s *System) Random() (err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Can only do random if there are a few things that we could pick.
	if s.maxAffirmationIndex() > 0 {

		// Have we gone through the slide show?
		if s.slideShowI >= len(s.slideShowIndexes)-1 {

			// Restart the slide show.
			s.slideShowI = 0
			s.slideShowIndexes = s.calculateSlideShowIndixes()

			// Now set the random slide show screen.
			s.activeAffirmationIndex = s.slideShowIndexes[s.slideShowI]

		} else {
			// Go to next slide of the slide show.
			s.slideShowI++
			s.activeAffirmationIndex = s.slideShowIndexes[s.slideShowI]
		}
	}

	return nil
}

// calculateSlideShowIndixes figures out the indixes for the slide show.
func (s *System) calculateSlideShowIndixes() (slideShowIndexes []int) {

	// Are we random?
	if s.slideShowRandom {

		// Get random slide show.
		slideShowIndexes = rand.Perm(len(s.affirmations))
		// If we are currently looking at the new first affirmation, move that affirmation to the end.
		if s.activeAffirmationIndex == slideShowIndexes[s.slideShowI] {
			slideShowIndexes = append(slideShowIndexes[1:], slideShowIndexes[:1]...)
		}

	} else {

		// Get ordered slide show.
		var orderedIndexes []int
		for i := 0; i < len(s.affirmations); i++ {
			orderedIndexes = append(orderedIndexes, i)
		}
		slideShowIndexes = orderedIndexes
	}

	return slideShowIndexes
}

// Left moves the pointer to displayed affirmation towards the top of the file.
func (s *System) Left() (err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// If we are currently displaying text and their is an image, then we should
	// no longer display the text, navigating the slide in two steps.
	if s.displayBoth && s.affirmations[s.activeAffirmationIndex].displayImage != nil {
		s.displayBoth = false // We've navigated to the beginning of a slide, show the image alone.
		return
	}

	if s.activeAffirmationIndex == 0 {
		s.activeAffirmationIndex = s.maxAffirmationIndex()
	} else {
		s.activeAffirmationIndex--
	}
	s.displayBoth = true // We just navigated backwards to a slide, show the text.
	return nil
}

// Right moves the pointer to displayed affirmation towards the bottom of the file.
func (s *System) Right() (err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// If we are currently not displaying text, just display it.
	if !s.displayBoth {
		s.displayBoth = true // Turn off on text so just the image exists.
		return               // We're done.
	}

	if s.activeAffirmationIndex == s.maxAffirmationIndex() {
		s.activeAffirmationIndex = 0
	} else {
		s.activeAffirmationIndex++
	}

	// Does this new slide have an image?
	if s.affirmations[s.activeAffirmationIndex].displayImage != nil {
		s.displayBoth = false // We've navigated to the beginning of a slide, show the image alone.
	}

	return nil
}

// maxAffirmationIndex gets the maximum affirmation index.
func (s *System) maxAffirmationIndex() (index int) {
	return len(s.affirmations) - 1
}

// GetDisplayTextImage gets the display text and image.
func (s *System) DisplayTextImage() (affirmationIndex int, displayText DisplayText, displayImage *DisplayImage, displayBoth, found bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if len(s.affirmations) == 0 {
		return 0, DisplayText{}, nil, false, false
	}

	return s.activeAffirmationIndex, s.affirmations[s.activeAffirmationIndex].displayText, s.affirmations[s.activeAffirmationIndex].displayImage, s.getDisplayBoth(), true
}

// CacheSlide caches a slide for quick rendering.
func (s *System) CacheSlide(affirmationIndex, winWidth, winHeight int, pixbuf *gdk.Pixbuf) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Clear the cache if our window size has changed.
	s.clearSlideCacheIfNecessary(winWidth, winHeight)

	// Cache the slide.
	s.cachedSlides[affirmationIndex] = pixbuf
}

// GetCachedSlide gets a cached slide if there is one for quick rendering.
func (s *System) GetCachedSlide(affirmationIndex, winWidth, winHeight int) (pixbuf *gdk.Pixbuf, found bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Clear the cache if our window size has changed.
	s.clearSlideCacheIfNecessary(winWidth, winHeight)

	// Get the cached value.
	pixbuf, found = s.cachedSlides[affirmationIndex]
	if !found {
		return nil, false
	}

	return pixbuf, true
}

// clearSlideCacheIfNecessary clears the cache of prerendered slides if width or height changed
func (s *System) clearSlideCacheIfNecessary(winWidth, winHeight int) {

	// Is this still the proper cache.
	if winWidth == s.cachedWidth && winHeight == s.cachedHeight {
		// The cache is still good.
		return
	}

	s.cachedWidth = winWidth
	s.cachedHeight = winHeight
	s.cachedSlides = map[int]*gdk.Pixbuf{}
}
