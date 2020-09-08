package conditioning

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type AffirmationSuite struct{}

var _ = Suite(&AffirmationSuite{})

// Add the tests.

func (s *AffirmationSuite) Test_Validate(c *C) {

	// A sample file.
	text := `

		// A fun title.

		// A comment.
		here is an affirmation!
		another one...

		// Another comment.
		yay!

		// Images.
		This is cool [something.jpeg]
		This is cool [something.jpeg:3.23]
		This is cool [something.jpeg:12,-34]
		This is cool [something.jpeg:12,-34:3.23]
		This is cool [   something.jpeg:3.23:12,-34   ]

		// Text.
		This is cool [b]
		This is cool [w:32]
		This is cool [b:12,-34]
		This is cool [w:12,-34:32]
		This is cool [   b:32:12,-34   ]

		// Image and text.
		This is cool [   something.jpeg:3.23:12,-34    b:32:12,-34   ]
		This is cool [    b:32:12,-34   something.jpeg:3.23:12,-34   ]
		`

	affirmations, title := parseAffirmations(text)
	c.Check(title, Equals, "A fun title.")
	c.Check(affirmations, DeepEquals, []Affirmation{
		Affirmation{
			Message: "here is an affirmation!",
		},
		Affirmation{
			Message: "another one...",
		},
		Affirmation{
			Message: "yay!",
		},

		// Images.
		Affirmation{
			Message: "This is cool",
			Image: AffirmationImage{
				Filename: "something.jpeg",
				OffsetX:  0,
				OffsetY:  0,
				Scale:    1.00,
			},
		},
		Affirmation{
			Message: "This is cool",
			Image: AffirmationImage{
				Filename: "something.jpeg",
				OffsetX:  0,
				OffsetY:  0,
				Scale:    3.23,
			},
		},
		Affirmation{
			Message: "This is cool",
			Image: AffirmationImage{
				Filename: "something.jpeg",
				OffsetX:  12,
				OffsetY:  -34,
				Scale:    1.00,
			},
		},
		Affirmation{
			Message: "This is cool",
			Image: AffirmationImage{
				Filename: "something.jpeg",
				OffsetX:  12,
				OffsetY:  -34,
				Scale:    3.23,
			},
		},
		Affirmation{
			Message: "This is cool",
			Image: AffirmationImage{
				Filename: "something.jpeg",
				OffsetX:  12,
				OffsetY:  -34,
				Scale:    3.23,
			},
		},

		// Text.
		Affirmation{
			Message: "This is cool",
			Text: TextProperties{
				Black:    BLACK,
				OffsetX:  0,
				OffsetY:  0,
				FontSize: 0,
			},
		},
		Affirmation{
			Message: "This is cool",
			Text: TextProperties{
				Black:    WHITE,
				OffsetX:  0,
				OffsetY:  0,
				FontSize: 32,
			},
		},
		Affirmation{
			Message: "This is cool",
			Text: TextProperties{
				Black:    BLACK,
				OffsetX:  12,
				OffsetY:  -34,
				FontSize: 0,
			},
		},
		Affirmation{
			Message: "This is cool",
			Text: TextProperties{
				Black:    WHITE,
				OffsetX:  12,
				OffsetY:  -34,
				FontSize: 32,
			},
		},
		Affirmation{
			Message: "This is cool",
			Text: TextProperties{
				Black:    BLACK,
				OffsetX:  12,
				OffsetY:  -34,
				FontSize: 32,
			},
		},

		// Image and text.
		Affirmation{
			Message: "This is cool",
			Image: AffirmationImage{
				Filename: "something.jpeg",
				OffsetX:  12,
				OffsetY:  -34,
				Scale:    3.23,
			},
			Text: TextProperties{
				Black:    BLACK,
				OffsetX:  12,
				OffsetY:  -34,
				FontSize: 32,
			},
		},
		Affirmation{
			Message: "This is cool",
			Image: AffirmationImage{
				Filename: "something.jpeg",
				OffsetX:  12,
				OffsetY:  -34,
				Scale:    3.23,
			},
			Text: TextProperties{
				Black:    BLACK,
				OffsetX:  12,
				OffsetY:  -34,
				FontSize: 32,
			},
		},
	})
}
