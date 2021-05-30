package inplaceenvsubst

type Config struct {
	FailOnMissingVariables bool
	RunInParallel          bool
	ErrorListener          ErrorListener
	WhitelistEnvs          StringSet
	BlacklistEnvs          StringSet
}

func (c *Config) ShouldIgnoreEnv(env string) bool {
	if c.WhitelistEnvs.IsEmpty() && c.BlacklistEnvs.IsEmpty() || !c.WhitelistEnvs.IsEmpty() && !c.BlacklistEnvs.IsEmpty() {
		return false
	}
	if !c.WhitelistEnvs.IsEmpty() {
		_, whitelisted := c.WhitelistEnvs[env]
		return !whitelisted
	} else {
		_, blacklisted := c.BlacklistEnvs[env]
		return blacklisted
	}
}
