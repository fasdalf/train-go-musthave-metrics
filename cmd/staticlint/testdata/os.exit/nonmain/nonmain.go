package nonmain

import (
	"fmt"
	"os"
)

func other() {
	fmt.Println("os.Exit()")
	os.Exit(0)
}
