package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	docker_root     = "/var/lib/docker"
	containers_root = filepath.Join(docker_root, "containers")
	volumes_root    = filepath.Join(docker_root, "volumes")
)

type Config struct {
	ID              string           `json:"ID"`
	Pid             int              `json:"Pid"`
	Name            string           `json:"Name"`
	ResolveConfPath string           `json:"ResolvConfPath"`
	HostnamePath    string           `json:"HostnamePath"`
	HostsPath       string           `json:"HostsPath"`
	LogPath         string           `json:"LogPath"`
	MountPoints     map[string]Mount `json:"MountPoints"`
}

type Mount struct {
	Name        string `json:"Name"`
	Destination string `json:"Destination"`
}

func GetAll(c *cli.Context) {
	dirs, err := ioutil.ReadDir(containers_root)
	if err != nil {
		fmt.Println("Readdir Error")
	}

	for _, dir := range dirs {
		fmt.Println("\nContainer " + dir.Name())
		if err := GetConVolumes(dir); err != nil {
			fmt.Println("Container " + dir.Name() + " has no volume details.")
			continue
		}
	}
}

func GetConVolumes(dir os.FileInfo) error {
	container_config_path := filepath.Join(containers_root, dir.Name(), "config.json")

	fi, _ := os.Open(container_config_path)
	defer fi.Close()

	fd, _ := ioutil.ReadAll(fi)
	// Get all data from config.json
	configData := string(fd)
	// Tranfer configData to []byte
	configDataBytes := []byte(configData)

	// Start to convert []bytes into map
	var config Config
	_ = json.Unmarshal(configDataBytes, &config)

	fmt.Println("Container ID:" + config.ID + "\nContainer Name:" + config.Name)

	//fmt.Println(config)

	var err error

	if err = checkSize("resolve.conf", config.ResolveConfPath); err != nil {
		fmt.Println("")
	}

	if err = checkSize("hostname", config.HostnamePath); err != nil {
		fmt.Println("")
	}

	if err = checkSize("hosts", config.HostsPath); err != nil {
		fmt.Println("")
	}

	if err = checkSize("json-log", config.LogPath); err != nil {
		fmt.Println("")
	}

	if err := checkDataVolume(config.MountPoints); err != nil {
		fmt.Println("")
	}
	return nil
}

func checkSize(filename string, path string) error {
	fmt.Print(filename + " Size: ")

	f, err := os.Lstat(path)

	if err != nil {
		return err
	}
	fmt.Println(f.Size())
	return nil
}

func checkDataVolume(mounts map[string]Mount) error {
	for _, value := range mounts {
		name := value.Name
		//fmt.Println("Destination:" + destination)
		volume_path := filepath.Join(volumes_root, name, "_data")
		//fmt.Println("Source:" + volume_path)

		if name != "" {
			// it means it is a data volume
			size := 0
			filepath.Walk(volume_path, func(_ string, file os.FileInfo, _ error) error {
				size += int(file.Size())
				return nil
			})
			fmt.Println("Volume Space Size:", size)
		}
	}
	return nil
}
