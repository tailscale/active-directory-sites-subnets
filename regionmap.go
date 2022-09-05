package main

import (
	"fmt"
	"strings"
)

type RegionMap map[string]string

func NewRegionMap() RegionMap {
	return make(RegionMap)
}

func (r RegionMap) Set(value string) error {
	fields := strings.Split(value, ":")
	region := fields[0]
	site := fields[1]

	r[region] = site

	return nil
}

func (r RegionMap) String() string {
	var sb strings.Builder
	for key, value := range r {
		fmt.Fprintf(&sb, "%q=%q\n", key, value)
	}
	return sb.String()
}
