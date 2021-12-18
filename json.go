package main

import (
	"encoding/json"
)

type Status struct {
	Self  Machine
	Peers map[string]Machine
}

type Machine struct {
	RawMachine
	DisplayName Name
}

func (s *Status) UnmarshalJSON(data []byte) error {
	rawStatus := new(rawStatus)

	if err := json.Unmarshal(data, &rawStatus); err != nil {
		return err
	}

	peers := map[string]Machine{}

	for name, rawPeer := range rawStatus.Peers {
		peers[name] = rawPeer.ToMachine(rawStatus.MagicDNSSuffix)
	}

	self := rawStatus.Self.ToMachine(rawStatus.MagicDNSSuffix)

	*s = Status{
		Self:  self,
		Peers: peers,
	}

	return nil
}

type rawStatus struct {
	Self           RawMachine            `json:"Self"`
	Peers          map[string]RawMachine `json:"Peer"`
	MagicDNSSuffix string                `json:"MagicDNSSuffix"`
}

type RawMachine struct {
	DNSName      string   `json:"DNSName"`
	HostName     string   `json:"HostName"`
	TailscaleIPs []string `json:"TailscaleIPs"`
}

func (rm RawMachine) ToMachine(dnsSuffix string) Machine {
	return Machine{
		RawMachine:  rm,
		DisplayName: dnsOrQuoteHostname(dnsSuffix, rm),
	}
}
