package main

import (
	"bufio"
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
		s.Port,
	}
}

func main() {
	services := []Service{
		{
			Name: "gcs-uploader",
			Port: "10004:10004",
		},
		{
			Name: "komiic-service",
			Port: "10002:10002",
		},
		{
			Name: "mongo",
			Port: "27017:27017",
		},
		{
			Name: "postgres",
			Port: "5432:5432",
		},
		{
			Name: "rabbitmq",
			Port: "5672:5672",
		},
		{
			Name: "redis",
			Port: "6379:6379",
		},
	}
	done := make(chan string)
	for _, s := range services {
		v := s
		go func() {
			cmd, err := runCmd(v)
			if err != nil {
				log.Fatal(err)
			}
			if err := cmd.Wait(); err != nil {
				log.Fatal(err)
			}
			done <- fmt.Sprintf("[%s] Closed", v.Name)
		}()
		time.Sleep(3 * time.Second)
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
