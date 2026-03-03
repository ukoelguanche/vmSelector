package model

import (
	"fmt"
	"strings"
)

type VM struct {
	VMID   int                    `json:"vmid"`
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Node   string                 `json:"node"`
	Status string                 `json:"status"`
	Config map[string]interface{} `json:"config,omitempty"`
}

type VMList struct {
	Data []VM `json:"data"`
}

func (vm *VM) HasSpecificGPU(gpuString string) bool {
	if vm.Config == nil {
		return false
	}

	for key, value := range vm.Config {
		if strings.HasPrefix(key, "hostpci") {
			valStr := fmt.Sprintf("%v", value)
			if strings.Contains(valStr, gpuString) {
				return true
			}
		}
	}
	return false
}

func (vm *VM) String() string {
	if vm == nil {
		return "<nil>"
	}
	return fmt.Sprintf("[%s|%d]", vm.Name, vm.VMID)
}

func (vmA *VM) Equals(vmB *VM) bool {
	if vmA == nil && vmB == nil {
		return true
	}
	return vmA != nil && vmB != nil && vmA.VMID == vmB.VMID
}

func (vm *VM) GetOS() string {
	description := vm.Config["description"].(string)
	lines := strings.Split(description, "\n")

	for _, line := range lines {
		if strings.Contains(line, "OS: ") {
			return strings.Replace(line, "OS: ", "", -1)
		}
	}

	return ""
}
