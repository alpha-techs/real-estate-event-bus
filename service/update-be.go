package service

import "os/exec"

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
