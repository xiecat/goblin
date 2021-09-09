package options

import (
	"fmt"
)

const banner = `
────────╔╗──╔╗
────────║║──║║
╔══╗╔══╗║╚═╗║║─╔╗╔══╗
║╔╗║║╔╗║║╔╗║║║─╠╣║╔╗╗
║╚╝║║╚╝║║╚╝║║╚╗║║║║║║
╚═╗║╚══╝╚══╝╚═╝╚╝╚╝╚╝
╔═╝║  %s - %s
╚══╝
`

// Version is the current version of goblin
var (
	Version = "unknown"
	Commit  = "unknown"
	Branch  = "unknown"
	Release = "unknown"
)

// showBanner is used to show the banner to the user
func showBanner() {
	fmt.Printf("%s", fmt.Sprintf(banner, Version, Release))
	//fmt.Printf("\t version: %s\n\n", Version)
	fmt.Printf("\tFrom: %s\n\n", "https://github.com/xiecat/goblin")
	fmt.Println("Please use this tool within the scope of the license.")
	fmt.Println("goblin is not responsible for any risks arising from the use of the tool.")
	fmt.Println("Use agrees to this statement\n ")
}
