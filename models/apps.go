package models

type GateApp struct {
	Model
	Host string // example.com
	Name string
	User string
}

func (a *GateApp) ToMap() map[string]any {
	return map[string]any{
		"name": a.Name,
	}
}

func (a *GateApp) ToJson() any { return a.ToMap() }
