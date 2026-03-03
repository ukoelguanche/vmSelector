package model

import (
	"encoding/json"
	"fmt"
	"log"

	"apodeiktikos.com/fbtest/httpClient"
	"apodeiktikos.com/fbtest/util"
)

func GetVMConfig(node string, vmid int) map[string]interface{} {
	url := fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmid)
	response := httpClient.DoRequest("GET", url, nil)
	defer response.Body.Close()

	var configBody struct {
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(response.Body).Decode(&configBody); err != nil {
		log.Fatal("Error decoding json")
	}

	return configBody.Data
}

func GetVMs() *VMList {
	response := httpClient.DoRequest("GET", "/cluster/resources?type=vm", nil)
	defer response.Body.Close()

	var pveRes VMList
	if err := json.NewDecoder(response.Body).Decode(&pveRes); err != nil {
		return nil
	}

	for i := range pveRes.Data {
		vm := &pveRes.Data[i]
		vm.Config = GetVMConfig(vm.Node, vm.VMID)
	}

	return &pveRes
}

func GetVMByName(hostname string) *VM {
	pveRes := GetVMs()
	if pveRes == nil {

	}
	for _, vm := range pveRes.Data {
		if vm.Name == hostname {
			return &vm
		}
	}
	return nil
}

func GetVMById(id int) *VM {
	pveRes := GetVMs()
	for _, vm := range pveRes.Data {
		if vm.VMID == id {
			return &vm
		}
	}
	return nil
}

func PowerOffVM(vm *VM) {
	PowerVM(vm, "stop")
}

func PowerOnVM(vm *VM) {
	PowerVM(vm, "start")
}

func PowerVM(vm *VM, action string) {
	url := fmt.Sprintf("/nodes/%s/qemu/%d/status/%s", vm.Node, vm.VMID, action)
	httpClient.DoRequest("POST", url, nil)
}

func SetVMDescription(vm *VM, description string) {
	url := fmt.Sprintf("/nodes/%s/qemu/%d/config", vm.Node, vm.VMID)
	options := httpClient.RequestOptions{
		Body: fmt.Sprintf("description=%s", description),
	}
	httpClient.DoRequest("POST", url, &options)
}

func GetRunningVMWithGPU() *VM {
	vms := GetVMs()

	for _, vm := range vms.Data {
		if vm.HasSpecificGPU(util.ContextStorage.GpuString) && vm.Status == "running" {

			return &vm
		}
	}

	return nil
}

func GetVMsWithGPU(gpuString string, centinelVM *VM) []VM {
	vms := GetVMs()
	if vms == nil {
		log.Fatal("Could not find any VMs")
	}

	var filtered []VM

	for _, vm := range vms.Data {
		if vm.HasSpecificGPU(gpuString) && vm.Name != centinelVM.Name {
			filtered = append(filtered, vm)
		}
	}

	filtered = append(filtered, *centinelVM)

	return filtered
}

func SwitchToVM(centinelVM *VM, targetVM VM) {
	if centinelVM.VMID == targetVM.VMID {
		SetVMDescription(centinelVM, "power_off")
	} else {
		SetVMDescription(centinelVM, fmt.Sprintf("target_vm_id %d", targetVM.VMID))
	}
	PowerOffVM(centinelVM)
}
