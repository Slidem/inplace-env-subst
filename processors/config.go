package processors

type Config interface {

	ShouldIgnoreEnv(val string) bool
}
