//go:build docker
// +build docker

package orchestrator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jspc/pipelines-orchestrator"
)

var validContainerProcessConfig = orchestrator.ProcessConfig{
	Name: "tests",
	Type: "container",
	ExecutionContext: map[string]string{
		"env":   `HELLO="WORLD"`,
		"image": "quay.io/podman/hello:latest",
	},
}

func TestNewContainerProcess(t *testing.T) {
	for _, test := range []struct {
		name        string
		pc          orchestrator.ProcessConfig
		expectError error
	}{
		{"empty config", orchestrator.ProcessConfig{}, orchestrator.ContainerImageMissingErr{}},
		{"full and valid config", validContainerProcessConfig, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := orchestrator.NewContainerProcess(test.pc)
			if err == nil && test.expectError != nil {
				t.Errorf("expected error, received none")
			} else if err != nil && test.expectError == nil {
				t.Errorf("unexpected error %#v", err)
			}

			if err != nil && test.expectError != nil {
				err.Error() // does nothing but increase codecoverage /shrug

				if !errors.Is(err, test.expectError) {
					t.Errorf("expected error of type %T, received %T", test.expectError, err)
				}
			}
		})
	}
}

func TestContainerProcess_Run(t *testing.T) {
	c, err := orchestrator.NewContainerProcess(validContainerProcessConfig)
	if err != nil {
		t.Fatal(err)
	}

	ps, err := c.Run(context.Background(), orchestrator.Event{})
	if err != nil {
		t.Fatal(err)
	}

	expect := orchestrator.ProcessSuccess
	if expect != ps.Status {
		t.Errorf("status should be %v, received %v", expect, ps.Status)
	}
}