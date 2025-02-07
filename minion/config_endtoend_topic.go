package minion

import (
	"fmt"
	"time"
)

type EndToEndTopicConfig struct {
	Enabled               bool          `koanf:"enabled"`
	Name                  string        `koanf:"name"`
	ReplicationFactor     int           `koanf:"replicationFactor"`
	PartitionsPerBroker   int           `koanf:"partitionsPerBroker"`
	ReconcilationInterval time.Duration `koanf:"reconcilationInterval"`
}

func (c *EndToEndTopicConfig) SetDefaults() {
	c.Enabled = true
	c.Name = "kminion-end-to-end"
	c.ReconcilationInterval = 10 * time.Minute
}

func (c *EndToEndTopicConfig) Validate() error {

	if c.ReplicationFactor < 1 {
		return fmt.Errorf("failed to parse replicationFactor, it should be more than 1, retrieved value %v", c.ReplicationFactor)
	}

	if c.PartitionsPerBroker < 1 {
		return fmt.Errorf("failed to parse partitionsPerBroker, it should be more than 1, retrieved value %v", c.ReplicationFactor)
	}

	// If the timeduration is 0s or 0ms or its variation of zero, it will be parsed as 0
	if c.ReconcilationInterval == 0 {
		return fmt.Errorf("failed to validate topic.ReconcilationInterval config, the duration can't be zero")
	}

	return nil
}
