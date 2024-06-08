package types

import (
	"github.com/airchains-network/junction/x/junction/types"
)

type StationInfo struct {
	StationType string `json:"stationType"`
	//DaType      string `json:"daType"`
}

type GenesisDataType struct {
	StationId          string
	Creator            string
	CreationTime       string
	TxHash             string
	Tracks             []string
	TracksVotingPowers []uint64
	VerificationKey    interface{}
	ExtraArg           types.StationArg
	StationInfo        StationInfo
}
