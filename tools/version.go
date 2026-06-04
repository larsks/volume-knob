package main

import "fmt"

const (
	version string = "1.2.0"
	vcs_ref string = "c1a7b639e7"
	vcs_date string = "2026-06-03 22:59:03 -0400"
)

func cmdVersion() {
	fmt.Printf("vkcfg (volume_knob config) version %s (%s) on %s\n", version, vcs_ref, vcs_date)
}
