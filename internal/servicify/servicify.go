package servicify

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const defaultOutputPath string = "/opt/delicious-gps"

const warning string = `Installing as a service will create /opt/delicious-gps and create a systemd
service definition in /etc/systemd/system/delcicious-gps.service

Are you sure you want to install delicious-gps as a service? (y/n): `

const serviceDefinitionDir string = "/etc/systemd/system"
const serviceDefinitionPath string = "/etc/systemd/system/delcicious-gps.service"
const serviceDefinition string = `[Unit]
Description=Delicious GPS Datalogger
After=gpsd.service
Requires=gpsd.service
StartLimitIntervalSec=0
AssertPathExists=/opt/delicious-gps

[Service]
Type=simple
Restart=always
RestartSec=2
User=root
ExecStart=/usr/local/sbin/delicious-gps --database /opt/delicious-gps/data.db

[Install]
WantedBy=multi-user.target`

func confirm() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(warning)
	text, _ := reader.ReadString('\n')
	if text == "y\n" {
		return true
	}
	return false
}

func InstallAsService() error {
	if !confirm() {
		fmt.Print("Aborting")
		return nil
	}

	if _, err := os.Stat(defaultOutputPath); os.IsNotExist(err) {
		mkdirErr := os.Mkdir(defaultOutputPath, 0664)
		if mkdirErr != nil {
			return errors.New(fmt.Sprintf("Failed to create %s, aborting", defaultOutputPath))
		}
	}

	if _, err := os.Stat(serviceDefinitionDir); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s does not exist, create it first, aborting", serviceDefinitionDir))
	}

	f, err := os.Create(serviceDefinitionPath)
	if err != nil {
		return err
	}

	_, err = f.WriteString(serviceDefinition)
	if err != nil {
		return err
	}

	fmt.Println("Written\n", serviceDefinitionPath)
	fmt.Println(`Start delicious-gps using:
systemctl enable delcious-gps
systemctl start delicous-gps`)

	return nil
}
