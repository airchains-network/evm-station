package junction

import (
	"time"
)

type AccountDetails struct {
	Address       string    `yaml:"address"`
	Mnemonic      string    `yaml:"mnemonic"`
	CreationTime  time.Time `yaml:"creationTime"`
	AddressPrefix string    `yaml:"addressPrefix"`
}
