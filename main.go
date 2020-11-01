package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

type Service struct {
	Name string
	Port string
}

func (s *Service) ToString() []string {
	return []string{
		"port-forward",
		"svc/" + s.Name,
		s.Port+":"+s.Port,
	}
}

func main() {
	services := []Service{
		{
			Name: "mongo",
			Port: "27017",
		},
		{
			Name: "postgres",
			Port: "5432",
		},
		{
			Name: "rabbitmq",
			Port: "5672",
		},
		{
			Name: "redis",
			Port: "6379",
		},
		{
			Name: "account-service",
			Port: "10003",
		},
		{
			Name: "gcs-uploader",
			Port: "10004",
		},
		{
			Name: "komiic-feedback",
			Port: "10005",
		},
		{
			Name: "komiic-service",
			Port: "10002",
		},
	}
	defer func() {
		for _, s := range services {
			fmt.Printf("\nClosing: %s\n", s.Port)
			cmd := exec.Command("fuser", "-k", fmt.Sprintf("%s/tcp", s.Port))
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	done := make(chan string)
	for _, s := range services {
		v := s
		go func() {
			cmd, err := runCmd(v)
			if err != nil {
				done <- fmt.Sprintf("[%s] Error %+v", v.Name, err)
			}
			if err := cmd.Wait(); err != nil {
				done <- fmt.Sprintf("[%s] Error %+v", v.Name, err)
			}
			done <- fmt.Sprintf("[%s] Closed", v.Name)
		}()
		time.Sleep(1 * time.Second)
	}

	log.Printf("\nDone: %s\n", <-done)

}

func runCmd(s Service) (*exec.Cmd, error) {
	cmd, out, err := getCmdReadCloser(s)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(*out)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("\n[%s] %s", s.Name, line)
		}
	}()

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func getCmdReadCloser(s Service) (*exec.Cmd, *io.ReadCloser, error) {
	c1 := exec.Command("kubectl", s.ToString()...)
	out, err := c1.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	return c1, &out, nil
}
