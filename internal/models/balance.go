package models

type Balance struct {
	Addr string `json:"addr"`
	Wei  string `json:"wei"`
	Eth  string `json:"eth"`
}
