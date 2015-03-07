package main

import (
	"strconv"
	"time"
)

type DateReporter struct {
	format string
}

func (d *DateReporter) Report() (string, error) {
	return time.Now().Format(d.format), nil
}

func (d *DateReporter) GetRefreshInterval() time.Duration {
	return time.Second
}

func NewDateReporter(format string) *DateReporter {
	return &DateReporter{
		format: format,
	}
}

type DummyReporter struct {
	count int
}

func (d *DummyReporter) Report() (string, error) {
	d.count++
	return strconv.Itoa(d.count), nil
}

func (d *DummyReporter) GetRefreshInterval() time.Duration {
	return time.Second * 10
}
