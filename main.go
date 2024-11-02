package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/StackExchange/wmi"
)

type Win32_PnPEntity struct {
	Name        string
	DeviceID    string
	Description string
}

func main() {
	var devices []Win32_PnPEntity

	query := "WHERE Description LIKE '%camera%' OR Description LIKE '%webcam%' OR Description LIKE '%video%'"
	err := wmi.Query("SELECT Name, DeviceID, Description FROM Win32_PnPEntity "+query, &devices)
	if err != nil {
		fmt.Println("Error querying WMI:", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println("No cameras found.")
		return
	}

	fmt.Println("Connected cameras:")
	for i, device := range devices {
		fmt.Printf("[%d] Name: %s\n", i, device.Name)
	}

	fmt.Print("Select a camera: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	choice := strings.TrimSpace(input)

	selectedDeviceID := devices[0].DeviceID
	for i, device := range devices {
		if fmt.Sprintf("%d", i) == choice {
			selectedDeviceID = device.DeviceID
			break
		}
	}

	fmt.Print("Type 'enable' or 'disable': ")
	scanner.Scan()
	action := strings.TrimSpace(scanner.Text())

	var psCommand string
	if action == "disable" {
		psCommand = fmt.Sprintf("Get-PnpDevice -InstanceId '%s' | Disable-PnpDevice -Confirm:$false", selectedDeviceID)
	} else if action == "enable" {
		psCommand = fmt.Sprintf("Get-PnpDevice -InstanceId '%s' | Enable-PnpDevice -Confirm:$false", selectedDeviceID)
	} else {
		fmt.Println("Invalid action. Please type 'enable' or 'disable'.")
		return
	}

	cmd := exec.Command("powershell", "-Command", psCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing PowerShell command:", err)
	} else {
		fmt.Println("Command output:", string(output))
	}
}
