package testdata

import "os"

func main() {
	os.Exit(1) // want "usage of os.Exit in the main function is not allowed"
}
