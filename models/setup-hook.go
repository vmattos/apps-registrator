package models

type SetupHook struct {
	Account   string   `json:"account"`
	Workspace string   `json:"workspace"`
	Hash      string   `json:"hash"`
	Removals  []string `json:"removals"`
	Additions []string `json:"additions"`
}
