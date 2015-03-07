package volume

// #cgo LDFLAGS: -lasound
// #include "getvol.h"
import "C"

import (
	"strconv"
	"time"
)

type VolumeReporter struct {
	// emtpy
}

func (v *VolumeReporter) getVolume() int {
	return int(C.get_volume_perc())
}

func (v *VolumeReporter) Report() (string, error) {
	volume := "Vol " + strconv.Itoa(v.getVolume())
	return volume, nil
}

func (v *VolumeReporter) GetRefreshInterval() time.Duration {
	return time.Second
}

func NewVolumeReporter() *VolumeReporter {
	return &VolumeReporter{}
}
