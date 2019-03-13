package dialog

import (
	"fmt"
	"strings"
	"time"

	"github.com/theparanoids/ashirt/termrec/fancy"
)

// ShowLoadingAnimation presents the given text plus a looping dot animation.
// Should be called as a goroutine, otherwise this is likely to be an infinite loop
//
// To stop, set the stopCheck parameter to true
func ShowLoadingAnimation(text string, stopCheck *bool) {
	count := 0
	max := 4
	dots := func() string {
		return strings.Repeat(".", count)
	}
	for {
		if !*stopCheck {
			fmt.Print(fancy.ClearLine(text+dots(), 0))
			count = (count + 1) % max
		}
		time.Sleep(500 * time.Millisecond)
	}
}
