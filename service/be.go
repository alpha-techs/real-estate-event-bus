package service

import (
	"errors"
	"os/exec"
	"strings"
)

func UpdateBe(version string) error {
	DockerRepo := "icylydia/real-estate-be-wt"
	ContainerName := "real-estate-be"

	err := exec.Command("docker", "pull", DockerRepo+":"+version).Run()
	if err != nil {
		return err
	}

	cid, err := exec.Command("docker", "ps", "-a", "--filter", "name="+ContainerName, "--format", "{{.Names}}").Output()
	if err != nil {
		return err
	}

	if len(cid) > 0 {
		err := exec.Command("docker", "rm", "-f", ContainerName).Run()
		if err != nil {
			return err
		}
	}

	err = exec.Command("docker", "run", "-d", "--name", ContainerName, "-p", "8000:9000", DockerRepo+":"+version).Run()
	if err != nil {
		return err
	}
	return nil
}

func GetBeVersion() (string, error) {
	cmd := "docker ps -a --filter name=real-estate-be --format {{.Image}} | awk -F ':' '{print $2}'"

	version, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		println(err.Error())
		return "", err
	}

	if len(version) == 0 {
		return "", errors.New("unable to get the version of the backend service")
	}

	trimmedVersion := strings.TrimSuffix(string(version), "\n")

	return trimmedVersion, nil
}
