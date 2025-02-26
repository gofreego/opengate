package keysources

type KeySourceConfig struct {
}

func NewKeySourceConfig() *KeySourceConfig {
	return &KeySourceConfig{}
}

func (c *KeySourceConfig) GetSigningKey() (string, error) {
	return "", nil
}
