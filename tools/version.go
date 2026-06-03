package main

import "fmt"

const (
	version  string = ""
	vcs_ref  string = ""
	vcs_date string = ""
)

func cmdVersion() {
	fmt.Printf("vkcfg (volume_knob config) version %s (%s) on %s\n", version, vcs_ref, vcs_date)
}
