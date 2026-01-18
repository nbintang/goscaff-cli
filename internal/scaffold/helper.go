package scaffold

import (
	"fmt" 
)

func printStep(title string, commands ...string) {
	fmt.Println()
	cStepTitle.Printf("  %s\n", title)
	for _, c := range commands {
		cBullet.Print("    $ ")
		cCmd.Println(c)
	}
} 
