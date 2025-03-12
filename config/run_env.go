package config

type RunEnv int

const (
	NoEnv RunEnv = iota
	Local
	Dev
	Prod
	OnPrem
)

func (r RunEnv) String() string {
	return [...]string{"No Env", "local", "dev", "prod", "on-prem"}[r]
}

func StringToRunEnv(s string) RunEnv {
	switch s {
	case "local":
		return Local
	case "dev":
		return Dev
	case "prod":
		return Prod
	case "on-prem":
		return OnPrem
	}
	return NoEnv
}
