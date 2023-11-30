package ethx

import (
	"errors"
	"log"
)

var DefaultEventConfig = EventConfig{
	IntervalBlocks: 500,
	OverrideBlocks: 1000,
	DelayBlocks:    4,
}

type EventConfig struct {
	IntervalBlocks, OverrideBlocks, DelayBlocks uint64
}

func (c *Clientx) newEventConfig(config []EventConfig) EventConfig {
	var _config EventConfig
	if len(config) > 0 {
		_config = config[0]
		_config.panicIfNotValid()
	} else {
		_config = DefaultEventConfig
	}
	return _config
}

func (e *EventConfig) panicIfNotValid() {
	if e.IntervalBlocks == 0 {
		panic(errors.New("EventConfig::require IntervalBlocks > 0, eg: 800"))
	}
	if e.IntervalBlocks+e.OverrideBlocks > 2000 {
		panic(errors.New("EventConfig::require IntervalBlocks + OverrideBlocks <= 2000, eg: 800 + 800"))
	}
	if e.DelayBlocks == 0 {
		log.Printf("[WARN] EventConfig::If you are tracking the latest logs, DelayBlocks==0 is risky, recommended >= 3. see: https://github.com/ethereum/go-ethereum/blob/master/core/types/log.go#L53")
	}
}
