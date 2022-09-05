package main

import (
	"strings"
	"testing"
)

func TestRegionMapSet(t *testing.T) {
	r := NewRegionMap()
	r.Set("nyc:New York City")
	r.Set("sfo:San Francisco")
	r.Set("tor:Toronto")

	if r["nyc"] != "New York City" {
		t.Errorf("region nyc got: %q, want: %q", r["nyc"], "New York City")
	}
	if r["sfo"] != "San Francisco" {
		t.Errorf("region sfo got: %q, want: %q", r["sfo"], "San Francisco")
	}
	if r["tor"] != "Toronto" {
		t.Errorf("region tor got: %q, want: %q", r["tor"], "Toronto")
	}
}

func TestRegionMapString(t *testing.T) {
	r := NewRegionMap()
	r.Set("nyc:New York City")
	r.Set("sfo:San Francisco")
	r.Set("tor:Toronto")

	s := r.String()

	wants := []string{`"nyc"="New York City"`, `"sfo"="San Francisco"`, `"tor"="Toronto"`}
	for _, want := range wants {
		if !strings.Contains(s, want) {
			t.Errorf("got: %q, does not contain: %q", s, want)
		}
	}
}
