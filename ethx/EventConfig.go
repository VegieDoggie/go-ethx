package ethx

import (
	"errors"
	"log"
)

type EventConfig struct {
	IntervalBlocks, OverrideBlocks, DelayBlocks uint64
}

func (c *Clientx) mustClientxConfig(config []*ClientxConfig) *ClientxConfig {
	var cfg *ClientxConfig
	if len(config) > 0 {
		cfg = config[0]
		cfg.Event.panicIfNotValid()
	} else {
		cfg = c.config
	}
	return cfg
}

func (e *EventConfig) panicIfNotValid() {
	if e.IntervalBlocks == 0 {
		panic(errors.New("EventConfig::require IntervalBlocks > 0, eg: 800"))
	}
	if e.IntervalBlocks+e.OverrideBlocks > 2000 {
		panic(errors.New("EventConfig::require IntervalBlocks + OverrideBlocks <= 2000, eg: 800 + 800"))
	}
	if e.DelayBlocks == 0 {
		log.Printf("[WARN] EventConfig::If you are tracking logs, recommended DelayBlocks >= 3 (or risky). see: https://github.com/ethereum/go-ethereum/blob/master/core/types/log.go#L53")
	}
}
