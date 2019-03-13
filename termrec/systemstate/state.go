package systemstate

type SystemState struct {
	termWidth  uint16
	termHeight uint16
}

var state = SystemState{}

func Current() SystemState {
	return state
}

func TermWidth() uint16 {
	return state.termWidth
}

func TermHeight() uint16 {
	return state.termHeight
}

func UpdateTermHeight(h uint16) {
	if h > state.termHeight {
		state.termHeight = h
	}
}

func UpdateTermWidth(w uint16) {
	if w > state.termWidth {
		state.termWidth = w
	}
}
