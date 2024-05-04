package types

type SequencerConfigs struct {
	DaRPC           string
	DaKey           string
	DaType          string
	JunctionRpc     string
	JunctionKeyName string
	StationId       string
}
type StationInfo struct {
	StationType string `json:"stationType"`
}
