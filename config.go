package inplaceenvsubst

type Config struct {
	FailOnMissingVariables bool
	RunInParallel          bool
	ErrorListener          ErrorListener
}
