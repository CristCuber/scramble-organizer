package main

import (
	"fmt"
	"scramble-organizer/organizer"
)

func main() {
	o := organizer.NewOrganizer()

	if err := o.OrganizScramble(); err != nil {
		fmt.Println(err)
	}
}
