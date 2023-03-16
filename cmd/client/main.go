package main

import "github.com/e-faizov/GophKeeper/internal/client"

func main() {
	//unauth zone
	client.RegOrAuth()

	//auth zone
	client.FirstAction()
}
