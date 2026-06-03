package main

import "fmt"

const (
	version string = "1.1.0"
	vcs_ref string = "630ae44d10"
	vcs_date string = "2026-06-02 23:00:17 -0400"
)

func cmdVersion() {
	fmt.Printf("vkcfg (volume_knob config) version %s (%s) on %s\n", version, vcs_ref, vcs_date)
}
