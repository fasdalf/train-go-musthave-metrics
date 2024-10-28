// Package printbuild - Print Build
package printbuild

import (
	"testing"
)

func TestPrint(t *testing.T) {
	bd := &Data{
		BuildVersion: "buildVersion",
		BuildDate:    "buildDate",
		BuildCommit:  "buildCommit",
	}
	bd.Print()
}
