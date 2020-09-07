package conditioning

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type ConfigSuite struct{}

var _ = Suite(&ConfigSuite{})

// Add the tests.

func (s *ConfigSuite) Test_Validate(c *C) {
	tests := []struct {
		config Config
		errstr string
	}{

		// All value set.
		{
			config: Config{
				SleepMilli:        1,
				ScreenWidth:       1,
				ScreenHeight:      1,
				FontFace:          "Georgia",
				FontSize:          1,
				BlackOutlineScale: 0.5,
				WhiteOutlineScale: 0.5,
			},
			errstr: ``,
		},

		// Check missing values.
		{
			config: Config{
				SleepMilli:        0,
				ScreenWidth:       1,
				ScreenHeight:      1,
				FontFace:          "Georgia",
				FontSize:          1,
				BlackOutlineScale: 0.5,
				WhiteOutlineScale: 0.5,
			},
			errstr: `invalid SleepMilli: 0`,
		},
		{
			config: Config{
				SleepMilli:        1,
				ScreenWidth:       0,
				ScreenHeight:      1,
				FontFace:          "Georgia",
				FontSize:          1,
				BlackOutlineScale: 0.5,
				WhiteOutlineScale: 0.5,
			},
			errstr: `invalid ScreenWidth: 0`,
		},
		{
			config: Config{
				SleepMilli:        1,
				ScreenWidth:       1,
				ScreenHeight:      0,
				FontFace:          "Georgia",
				FontSize:          1,
				BlackOutlineScale: 0.5,
				WhiteOutlineScale: 0.5,
			},
			errstr: `invalid ScreenHeight: 0`,
		},
		{
			config: Config{
				SleepMilli:        1,
				ScreenWidth:       1,
				ScreenHeight:      1,
				FontFace:          "",
				FontSize:          1,
				BlackOutlineScale: 0.5,
				WhiteOutlineScale: 0.5,
			},
			errstr: `invalid FontFace: ''`,
		},
		{
			config: Config{
				SleepMilli:        1,
				ScreenWidth:       1,
				ScreenHeight:      1,
				FontFace:          "Georgia",
				FontSize:          0,
				BlackOutlineScale: 0.5,
				WhiteOutlineScale: 0.5,
			},
			errstr: `invalid FontSize: 0`,
		},
		{
			config: Config{
				SleepMilli:        1,
				ScreenWidth:       1,
				ScreenHeight:      1,
				FontFace:          "Georgia",
				FontSize:          1,
				BlackOutlineScale: 0,
				WhiteOutlineScale: 0.5,
			},
			errstr: `invalid BlackOutlineScale: 0`,
		},
		{
			config: Config{
				SleepMilli:        1,
				ScreenWidth:       1,
				ScreenHeight:      1,
				FontFace:          "Georgia",
				FontSize:          1,
				BlackOutlineScale: 0.5,
				WhiteOutlineScale: 0,
			},
			errstr: `invalid WhiteOutlineScale: 0`,
		},
	}
	for i, test := range tests {
		comment := Commentf("Case %v: %v", i, test)

		err := test.config.Validate()
		if test.errstr == "" {
			c.Check(err, IsNil, comment)
		} else {
			c.Check(err, ErrorEquals, test.errstr, comment)
		}
	}
}
