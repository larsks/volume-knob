package main

import "fmt"

const (
	version string = "1.2.1"
	vcs_ref string = "93ed8fea8c"
	vcs_date string = "2026-06-04 08:41:55 -0400"
)

func cmdVersion() {
	fmt.Printf("vkcfg (volume_knob config) version %s (%s) on %s\n", version, vcs_ref, vcs_date)
}
