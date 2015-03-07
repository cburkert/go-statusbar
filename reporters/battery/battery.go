package battery

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"
)

type PowerReporter struct {
	classDir string
}

func (p *PowerReporter) readClass(pathSuffix string) (string, error) {
	raw, err := ioutil.ReadFile(path.Join(p.classDir, pathSuffix))
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (p *PowerReporter) getBatteryRates() ([]string, error) {
	dirEnts, err := ioutil.ReadDir(p.classDir)
	if err != nil {
		return nil, err
	}
	batteries := make([]string, 0, len(dirEnts))
	for _, dirEnt := range dirEnts {
		if strings.HasPrefix(dirEnt.Name(), "BAT") {
			batteries = append(batteries, dirEnt.Name())
		}
	}

	rates := make([]string, 0, len(batteries))
	for _, battery := range batteries {
		full, err := p.readClass(path.Join(battery, "charge_full_design"))
		if err != nil {
			return nil, err
		}
		fullInt, err := strconv.Atoi(strings.TrimSpace(full))
		if err != nil {
			return nil, err
		}
		now, err := p.readClass(path.Join(battery, "charge_now"))
		if err != nil {
			return nil, err
		}
		nowInt, err := strconv.Atoi(strings.TrimSpace(now))
		if err != nil {
			return nil, err
		}
		rate := float64(nowInt*100) / float64(fullInt)
		rates = append(rates, strconv.FormatFloat(rate, 'f', 0, 64))
	}

	return rates, nil
}

func (p *PowerReporter) OnAC() (bool, error) {
	raw, err := p.readClass(path.Join("AC", "online"))
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(raw) == "1", nil
}

func (p *PowerReporter) Report() (string, error) {
	online, err := p.OnAC()
	if err != nil {
		return "", err
	}
	var report string
	if online {
		report = "AC"
	} else {
		rates, err := p.getBatteryRates()
		if err != nil {
			return "", err
		}
		report = "BAT " + strings.Join(rates, " ")
	}
	return report, nil
}

func (p *PowerReporter) GetRefreshInterval() time.Duration {
	return time.Minute
}

func NewPowerReporter(classDir string) *PowerReporter {
	return &PowerReporter{
		classDir: classDir,
	}
}
