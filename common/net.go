package common

type WakeCommand struct {
	Subnet     string `json:"subnet"`
	MacAddress string `json:"mac"`
}

type PingCommand struct {
	Subnet    string `json:"subnet"`
	IpAddress string `json:"ip"`
}

type RunnerResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"msg"`
}
