package types

import (
	"encoding/json"
	"errors"
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
	Leader        string
	Services      []Service
	EndPointIndex map[string]int
	Nodes         map[string][]IP_PORT_ROW
}

func (B *BS_State) InitState(bod []byte) {
	B.Leader = "None"
	json.Unmarshal(bod, &B.Services)
	B.Nodes = make(map[string][]IP_PORT_ROW)
	B.EndPointIndex = map[string]int{}
	for _, service := range B.Services {
		B.Nodes[service.ServiceName] = make([]IP_PORT_ROW, 0)
		for _, ep := range service.ServiceEndpoints {
			B.EndPointIndex[ep] = 0
		}
	}
}

func (B *BS_State) GetAvailabe(sname string, ename string) (IP_PORT_ROW, error) {
	// Check for the size of nodes also it should not exceed and increment in roud robin fashion
	var output IP_PORT_ROW
	if index, exists := B.EndPointIndex[ename]; exists {
		if len(B.Nodes[sname]) > 0 {
			output = B.Nodes[sname][index]
			B.EndPointIndex[ename] = ((B.EndPointIndex[ename] + 1) % len(B.Nodes[sname]))
			return output, nil
		}
	}
	return output, errors.New("not a valid thing")
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
	fmt.Println(bod.Sname)
	_, exists := B.Nodes[bod.Sname]
	if exists {
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
