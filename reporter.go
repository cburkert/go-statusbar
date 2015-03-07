// go-statusbar project main.go
package main

import (
	"log"
	"math/big"
	"os/exec"
	"strings"
	"time"
)

type Reporter interface {
	Report() (string, error)
	GetRefreshInterval() time.Duration
}

type StatusBar struct {
	reporterInfos []reporterInfo
	separator     string
	tickDuration  time.Duration
}

type reporterInfo struct {
	reporter   Reporter
	lastReport string
	tickSkips  int
}

func (s *StatusBar) AddReporter(reporter Reporter) {
	refreshInterval := reporter.GetRefreshInterval()
	if s.tickDuration == 0 {
		s.tickDuration = refreshInterval
	} else {
		curTick := big.NewInt(s.tickDuration.Nanoseconds())
		newIntverval := big.NewInt(refreshInterval.Nanoseconds())
		newTick := big.NewInt(0)
		newTick.GCD(nil, nil, curTick, newIntverval)
		s.tickDuration = time.Duration(newTick.Int64())
	}
	repInfo := reporterInfo{
		reporter:   reporter,
		lastReport: "",
		tickSkips:  0,
	}
	s.reporterInfos = append(s.reporterInfos, repInfo)
}

func (s *StatusBar) build() (string, error) {
	reports := make([]string, len(s.reporterInfos))
	for i := range s.reporterInfos {
		repInfo := &s.reporterInfos[i]
		if repInfo.tickSkips > 0 {
			repInfo.tickSkips--
			reports[i] = repInfo.lastReport
		} else {
			repInfo.tickSkips = int(repInfo.reporter.GetRefreshInterval() / s.tickDuration)
			report, err := repInfo.reporter.Report()
			if err != nil {
				return "", err
			}
			repInfo.lastReport = report
			reports[i] = report
		}
	}
	return strings.Join(reports, s.separator), nil
}

func (s *StatusBar) buildOldReport() string {
	reports := make([]string, len(s.reporterInfos))
	for i := range s.reporterInfos {
		repInfo := &s.reporterInfos[i]
		reports[i] = repInfo.lastReport
	}
	return strings.Join(reports, s.separator)
}

func (s *StatusBar) Update() error {
	oldReport := s.buildOldReport()
	report, err := s.build()
	if err != nil {
		return err
	}
	if report != oldReport {
		exec.Command("xsetroot", "-name", report).Run()
	}
	return nil
}

func (s *StatusBar) Run() {
	for {
		err := s.Update()
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(s.tickDuration)
	}
}

func NewStatusBar(separator string) *StatusBar {
	return &StatusBar{
		separator: separator,
	}
}
