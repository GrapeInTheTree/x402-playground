package tui

// NavigateMsg requests navigation to a specific page.
type NavigateMsg struct {
	Page Page
}

// BackMsg requests navigation back to the previous page.
type BackMsg struct{}

// ErrorMsg carries an error to display.
type ErrorMsg struct {
	Err error
}

// StatusMsg updates the status bar text.
type StatusMsg struct {
	Text string
}
