package models

type Service struct {
	Name    string
	Version int
	Schema  map[string]struct {
		JsonType string
	}
	EndpointUrl    string
	SupportedHooks []string
}
