// +build mage

// under review. Am seeing if it is worth using magefiles or not

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

const (
	packageName = "aws-region-stats"
	GoVersion   = "1.10.2"
)

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var (
	goexe  = "go"
	docker = sh.RunCmd("docker")
)

func getDep() error {
	return sh.Run(goexe, "get", "-u", "github.com/golang/dep/cmd/dep")
}

// Install Go Dep and sync the tools vendored dependencies
func Vendor() error {
	mg.Deps(getDep)
	return sh.Run("dep", "ensure")
}

func dockerBuildCmds(goVersion string) ([]string, error) {
	localGoPath := os.Getenv("GOPATH")
	if localGoPath == "" {
		return []string{}, errors.New("GOPATH not set")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return []string{}, errors.Wrap(err, "Problem getting current Directory")
	}

	buildString := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%s:/go", localGoPath),
		"-v", fmt.Sprintf("%s:/go/src/app", currentDir),
		"-w", "/go/src/app",
		fmt.Sprintf("golang:%s", goVersion),
	}

	return buildString, nil
}

// Build uses docker to build binaries for different OSes.
func Build() error {
	fmt.Println("Building OSX App...")
	buildCmds, err := dockerBuildCmds(GoVersion)
	if err != nil {
		return errors.Wrap(err, "Problem getting build commands")
	}

	if err := docker(buildCmds...); err != nil {
		return errors.Wrap(err, "Problem running Docker")
	}

	return nil
}
