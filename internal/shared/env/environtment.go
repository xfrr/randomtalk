package env

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentTest        Environment = "test"
	ProductionEnvironment  Environment = "production"
)

func (e Environment) String() string {
	return string(e)
}
