package core

type Observation struct {
	Success bool
	Output  string
}

type Step struct {
	Decision    Decision
	Observation Observation
}
