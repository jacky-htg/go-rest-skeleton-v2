package database

import (
	"bytes"
	"fmt"
	"os/exec"
)

// StartContainer runs a mysql container to execute commands.
func StartPostgresContainer() {
	cmd := exec.Command("docker", "run", "-d", "--name", "postgres_temp", "--publish", "54320:5432", "--env", "POSTGRES_PASSWORD=1234", "postgres:16-alpine")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("could not start docker for postgres", err)
	}

}

// StopContainer stops and removes the specified container.
func StopPostgresContainer() {
	if err := exec.Command("docker", "container", "rm", "-f", "postgres_temp").Run(); err != nil {
		fmt.Println("could not stop postgres container", err)
	}
}

func StartRedisContainer() {
	cmd := exec.Command("docker", "run", "-d", "--name", "redis_temp", "--publish", "63790:6379", "redis:latest")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("could not start docker for redis", err)
	}

}

// StopContainer stops and removes the specified container.
func StopRedisContainer() {
	if err := exec.Command("docker", "container", "rm", "-f", "redis_temp").Run(); err != nil {
		fmt.Println("could not stop redis container", err)
	}
}
