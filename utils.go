package apachepulsar

import (
	"fmt"
	"strconv"
	"time"
)

type configValidator struct {
	cfg map[string]string
	err error
}

func newConfigValidator(cfg map[string]string) *configValidator {
	return &configValidator{
		cfg: cfg,
		err: nil,
	}
}

func (c *configValidator) RequiredStr(key string) string {
	if c.err != nil {
		return ""
	}

	value, ok := c.cfg[key]
	if !ok {
		c.err = fmt.Errorf("required parameter %s not found", key)
		return ""
	}
	if value == "" {
		c.err = fmt.Errorf("required parameter %s is empty", key)
		return ""
	}

	return value
}

func (c *configValidator) OptionalStr(key string) string {
	if c.err != nil {
		return ""
	}

	value, ok := c.cfg[key]
	if !ok {
		return ""
	}

	return value
}

func (c *configValidator) OptionalInt(key string) int {
	if c.err != nil {
		return 0
	}

	value, ok := c.cfg[key]
	if !ok {
		return 0
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		c.err = fmt.Errorf("invalid int for %s: %w", key, err)
		return 0
	}

	return v
}

func (c *configValidator) OptionalBool(key string) bool {
	if c.err != nil {
		return false
	}

	value, ok := c.cfg[key]
	if !ok {
		return false
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		c.err = fmt.Errorf("invalid bool for %s: %w", key, err)
		return false
	}

	return v
}

func (c *configValidator) OptionalDuration(key string) time.Duration {
	if c.err != nil {
		return 0
	}

	value, ok := c.cfg[key]
	if !ok {
		return 0
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		c.err = fmt.Errorf("invalid duration for %s: %w", key, err)
		return 0
	}

	return duration
}

func (c *configValidator) Error() error {
	return c.err
}

func parseSubscriptionType(v *configValidator, key string) SubscriptionType {
	if v.err != nil {
		return ""
	}

	value, ok := v.cfg[key]
	if !ok {
		return ""
	}

	switch value {
	case "exclusive":
		return Exclusive
	case "shared":
		return Shared
	case "failover":
		return Failover
	case "key_shared":
		return KeyShared
	default:
		v.err = fmt.Errorf("invalid subscription type %s", value)
		return ""
	}
}
