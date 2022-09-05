package main

import (
	"context"
	ts "tailscale.com/client/tailscale"
)

type Peer struct {
	ips  []string
	derp string
}

func getPeers(ctx context.Context, lc ts.LocalClient) (peers []Peer, err error) {
	status, err := lc.Status(ctx)
	if err != nil {
		return peers, err
	}

	peers = make([]Peer, len(status.Peer))
	peernum := 0
	for _, peer := range status.Peer {
		ips := make([]string, len(peer.TailscaleIPs))
		for idx, ip := range peer.TailscaleIPs {
			suffix := "/32"
			if ip.Is6() {
				suffix = "/128"
			}
			ips[idx] = ip.String() + suffix
		}
		peers[peernum] = Peer{ips, peer.Relay}
		peernum += 1
	}

	return peers, nil
}
