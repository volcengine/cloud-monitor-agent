package display

import (
	"fmt"
)

func cursorDownRow(num string) {
	fmt.Printf("\033[" + num + "B")
}

func cursorUpRow(num string) {
	fmt.Printf("\033[" + num + "A")
}

func cursorLeftHead() {
	fmt.Printf("\r")
}

func cursorHide() {
	fmt.Printf("\033[?25l")
}

func cursorShow() {
	fmt.Printf("\033[?25h")
}

func cursorBackPriLocation(num string) {
	cursorLeftHead()
	cursorHide()
	cursorUpRow(num)
}

func cursorReset(num string) {
	cursorShow()
	cursorDownRow(num)
}
