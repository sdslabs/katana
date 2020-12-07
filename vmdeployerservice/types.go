package vmdeployerservice

type allocation struct {
	ID     string `json:"AllocId"`
	IP     string `json:"IP"`
	Serial string `json:"Serial"`
	Pid    string `json:"Pid"`
}

type status struct {
	Error interface{}
	Data  allocation
}
