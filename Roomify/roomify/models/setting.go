package models

type Setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SettingRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}