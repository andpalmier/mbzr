package cmd

import "fmt"

// printRootHelp displays the help message for the root command
func printRootHelp() {
	fmt.Println("mbzr is a command-line tool for interacting with MalwareBazaar API.\n  built by @andpalmier\n")
	fmt.Println("Usage:\n  mbzr <command> [flags]\n\nAvailable Commands:")
	fmt.Println("  query			Query MalwareBazaar for information about a sample")
	fmt.Println("  download		Download a sample by its sha256 hash")
	fmt.Println("  upload		Upload a file or directory to MalwareBazaar")
	fmt.Println("  comment		Add a comment to a malware sample")
	fmt.Println("  recent_detections	Get recent malware detections within a specified timeframe")
	fmt.Println("  cscb			Query the Code Signing Certificate Blocklist (CSCB)")
	fmt.Println("  update		Update metadata of a malware sample")
	fmt.Println("\nUse \"mbzr <command> -h or --help\" for more information about a command.")
	fmt.Println("\nExamples:")
	fmt.Println("  mbzr query -tag Emotet -limit 50")
	fmt.Println("  mbzr download -hash ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652")
	fmt.Println("  mbzr upload sample.exe")
	fmt.Println("  mbzr comment -sha256 ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652 -comment 'This is a test comment'")
	fmt.Println("  mbzr recent_detections -hours 12")
	fmt.Println("  mbzr cscb\n")
}
