package types

import (
	"encoding/json"
	"fmt"
)

type Service struct {
	ServiceName      string   `json:"name"`
	ServiceEndpoints []string `json:"endpoints"`
	CanLead          bool     `json:"canlead"`
}

type IP_Record struct {
	Leader   *string
	IP_Table map[string]int
}

type BS_State struct {
	Leader   string
	Services []Service
	Nodes    map[string][]IP_PORT_ROW
}

func (B *BS_State) InitState(bod []byte) {
	B.Leader = "None"
	json.Unmarshal(bod, &B.Services)
	B.Nodes = make(map[string][]IP_PORT_ROW)
	for _, service := range B.Services {
		B.Nodes[service.ServiceName] = make([]IP_PORT_ROW, 0)
	}
}
func CheckCapability(services []Service, sname string) bool {
	output := false
	for _, i := range services {
		if i.ServiceName == sname {
			output = i.CanLead
			break
		}
	}
	return output
}
func (B *BS_State) CheckLeader(bod IP_PORT) bool {
	output := false
	if _, exists := B.Nodes[bod.Sname]; exists {
		if CheckCapability(B.Services, bod.Sname) {

			if B.Leader == "None" {
				output = true
				B.Nodes[bod.Sname] = append(B.Nodes[bod.Sname], IP_PORT_ROW{
					Ip:   bod.Ip,
					Port: bod.Port,
				})
				B.Leader = bod.ToString()
			} else {
				B.Nodes[bod.Sname] = append(B.Nodes[bod.Sname], IP_PORT_ROW{
					Ip:   bod.Ip,
					Port: bod.Port,
				})
				output = false
			}
		} else {
			B.Nodes[bod.Sname] = append(B.Nodes[bod.Sname], IP_PORT_ROW{
				Ip:   bod.Ip,
				Port: bod.Port,
			})
			output = false

		}
	}
	fmt.Printf("%v\n", B.Nodes)
	return output

}
func (B *BS_State) GetLeader() []byte {
	return []byte(B.Leader)
}
