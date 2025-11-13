package models

import "fmt"

type Service struct {
	Name       string `yaml:"name"`
	Repository string `yaml:"repo"`
	RunCommand string `yaml:"run_command"`
}

func (s *Service) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	if s.Repository == "" {
		return fmt.Errorf("service '%s' must reference a repository", s.Name)
	}

	if s.RunCommand == "" {
		return fmt.Errorf("service '%s' must have a run_command", s.Name)
	}

	return nil
}
