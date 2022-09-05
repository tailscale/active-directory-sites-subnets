package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	ts "tailscale.com/client/tailscale"
)

var (
	sites       = NewRegionMap()
	lc          = ts.LocalClient{}
	defaultSite string
)

func getDerpRegions(ctx context.Context, lc ts.LocalClient) (derps map[string]int) {
	dm, err := lc.CurrentDERPMap(ctx)
	if err != nil {
		log.Fatal(err)
	}

	derps = make(map[string]int, 0)
	for _, region := range dm.Regions {
		derps[region.RegionCode] = 0
	}

	return
}

func processFlags(ctx context.Context) {
	flag.Var(sites, "map", "derp:site where derp is the regionCode from https://login.tailscale.com/derpmap/default and site is the ActiveDirectory site name")
	mapDefault := flag.String("map-default", "", "default ActiveDirectory site if not overridden by a map argument")
	flag.Parse()

	derpRegions := getDerpRegions(ctx, lc)
	failed := false
	for region, _ := range sites {
		if _, ok := derpRegions[region]; !ok {
			fmt.Printf("invalid DERP region %q in map argument\n", region)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}

	defaultSite = *mapDefault
}

func siteForDerp(derp string) string {
	if site, ok := sites[derp]; ok {
		return site
	}
	return defaultSite
}

func main() {
	ctx := context.Background()
	processFlags(ctx)

	peers, err := getPeers(ctx, lc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(`Get-ADReplicationSubnet -Filter "Location -like '*Tailscale'" | Remove-ADReplicationSubnet`)
	for _, peer := range peers {
		site := siteForDerp(peer.derp)
		if site == "" {
			continue
		}
		for _, ip := range peer.ips {
			fmt.Printf("New-ADReplicationSubnet -Name '%s' -Location 'Tailscale' -Site %s\n", ip, site)
		}
	}
}
