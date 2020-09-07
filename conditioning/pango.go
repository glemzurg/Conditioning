package conditioning

import (
	"strings"
)

// PangoMarkup
func PangoMarkup(text string) (markup string) {
	markup = text

	// Replace all / with <i> and </i>
	open := true
	for strings.Index(markup, "/") != -1 {
		if open {
			markup = strings.Replace(markup, "/", "<i>", 1)
		} else {
			markup = strings.Replace(markup, "/", "<~i>", 1)
		}
		open = !open
	}
	markup = strings.ReplaceAll(markup, "~", "/")

	// Replace all * with <b> and </b>
	open = true
	for strings.Index(markup, "*") != -1 {
		if open {
			markup = strings.Replace(markup, "*", "<b>", 1)
		} else {
			markup = strings.Replace(markup, "*", "</b>", 1)
		}
		open = !open
	}

	return "<span>" + markup + "</span>"
}
