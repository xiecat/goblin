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
	fmt.Printf("%s\n", fmt.Sprintf(banner, Version, Release))
	fmt.Printf("\t version: %s\n\n", Version)

	fmt.Printf("Use with caution. You are responsible for your actions\n")
	fmt.Printf("Developers assume no liability and are not responsible for any misuse or damage.\n")
}
