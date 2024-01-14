package apachepulsar

import "fmt"

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

func (c *configValidator) Required(key string) string {
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

func (c *configValidator) Optional(key string) string {
	if c.err != nil {
		return ""
	}

	value, ok := c.cfg[key]
	if !ok {
		return ""
	}

	return value
}

func (c configValidator) Error() error {
	return c.err
}
