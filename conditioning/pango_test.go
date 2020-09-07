package conditioning

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type PangoSuite struct{}

var _ = Suite(&PangoSuite{})

// Add the tests.

func (s *PangoSuite) Test_PandoMarkup(c *C) {
	tests := []struct {
		text   string
		markup string
	}{
		{``, `<span></span>`},
		{`something here`, `<span>something here</span>`},
		{`something /here/`, `<span>something <i>here</i></span>`},
		{`something *here*`, `<span>something <b>here</b></span>`},
		{`something */here/*`, `<span>something <b><i>here</i></b></span>`},
		{`something /*here*/`, `<span>something <i><b>here</b></i></span>`},
		{`something *here* *again*`, `<span>something <b>here</b> <b>again</b></span>`},
	}
	for i, test := range tests {
		comment := Commentf("Case %v: %v", i, test)
		c.Check(PangoMarkup(test.text), Equals, test.markup, comment)
	}
}
