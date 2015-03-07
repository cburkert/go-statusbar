package main

import (
	"github.com/cburkert/go-statusbar/reporters/battery"
	"github.com/cburkert/go-statusbar/reporters/volume"
)

func main() {
	statusBar := NewStatusBar(" | ")
	statusBar.AddReporter(volume.NewVolumeReporter())
	statusBar.AddReporter(battery.NewPowerReporter("/sys/class/power_supply/"))
	statusBar.AddReporter(&DateReporter{"Mon 02 Ý 15:04"})
	statusBar.Run()
}
