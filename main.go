//go:generate statik -dest=internal -p statikdata -src=assets/templates/ -include=*.gotmpl

package main

import (
	"github.com/shakahl/gqlassist/cmd"
)

func main() {
	cmd.Execute()
}
