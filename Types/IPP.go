package types

import "fmt"

type IP_PORT struct {
	Ip    string `json:"IP"`
	Port  int    `json:"PORT"`
	Sname string `json:"SERVICE_NAME"`
}

type IP_PORT_ROW struct {
	Ip   string
	Port int
}

func (I *IP_PORT_ROW) ToString() string {
	output := fmt.Sprintf("http://%s:%d", I.Ip, I.Port)
	return output
}
func (I *IP_PORT_ROW) ToBytes() []byte {
	output := fmt.Appendf(nil, "IP:%s,PORT:%d", I.Ip, I.Port)
	return output
}

func (I *IP_PORT) ToString() string {
	return fmt.Sprintf("IP:%s,PORT:%d", I.Ip, I.Port)
}
