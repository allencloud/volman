package main

import (
	"./utils"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func GetOne(c *cli.Context) {
	dirs, err := ioutil.ReadDir(containers_root)
	if err != nil {
		fmt.Println("Readdir Error")
	}
	sum := 0
	var rightone os.FileInfo
	input := c.Args()[0]
	fmt.Println(input)
	for _, dir := range dirs {
		if strings.Contains(dir.Name(), input) == true {
			rightone = dir
			sum++
			if sum >= 2 {
				fmt.Println("More than 2 containers' ID has your input prefix.")
				fmt.Println("Please check your input.")
				os.Exit(1)
			}
		}
	}

	if err := GetConVolumes(rightone); err != nil {
		fmt.Println("Container " + rightone.Name() + " has no volume details.")
	}
}

func GetAll(c *cli.Context) {
	dirs, err := ioutil.ReadDir(containers_root)
	if err != nil {
		fmt.Println("Readdir Error")
	}

	for _, dir := range dirs {
		if err := GetConVolumes(dir); err != nil {
			fmt.Println("Container " + dir.Name() + " has no volume details.")
			continue
		}
		fmt.Println("\n")
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

	fmt.Println("Container ID: " + config.ID + "\nContainer Name: " + config.Name)

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
	fmt.Println(utils.Convert(f.Size()))
	return nil
}

func checkDataVolume(mounts map[string]Mount) error {
	index := 1
	for _, value := range mounts {
		name := value.Name
		if name != "" {
			//fmt.Println("Destination:" + destination)
			volume_path := filepath.Join(volumes_root, name, "_data")
			fmt.Println(index, ". SourceDir: "+volume_path)
			index++

			// it means it is a data volume
			size := int64(0)
			filepath.Walk(volume_path, func(_ string, file os.FileInfo, _ error) error {
				size += file.Size()
				return nil
			})
			result := utils.Convert(size)
			fmt.Println("Volume Space Size:", result)
		}
	}

	return nil
}
