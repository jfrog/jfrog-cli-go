package golang

import (
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/artifactory/utils/golang/project"
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/utils/config"
	goutils "github.com/jfrog/jfrog-cli-go/jfrog-cli/artifactory/utils/golang"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/utils/errorutils"
	"os/exec"
	"strings"
)

func Publish(publishPackage bool, dependencies, targetRepo, version, buildName, buildNumber string, details *config.ArtifactoryDetails) (succeeded, failed int, err error) {
	err = validatePrerequisites()
	if err != nil {
		return
	}

	isCollectBuildInfo := len(buildName) > 0 && len(buildNumber) > 0
	if isCollectBuildInfo {
		err = utils.SaveBuildGeneralDetails(buildName, buildNumber)
		if err != nil {
			return
		}
	}

	goProject, err := project.Load(version)
	if err != nil {
		return
	}

	// Publish the package to Artifactory
	if publishPackage {
		err = goProject.PublishPackage(targetRepo, buildName, buildNumber, details)
		if err != nil {
			return
		}
	}

	// Publish the package dependencies to Artifactory
	depsList := strings.Split(dependencies, ",")
	if len(depsList) > 0 {
		succeeded, failed, err = goProject.PublishDependencies(targetRepo, details, depsList)
	}
	if publishPackage {
		succeeded++
	}
	if err != nil {
		return
	}

	// Publish the build-info to Artifactory
	if isCollectBuildInfo {
		err = utils.SaveBuildInfo(buildName, buildNumber, goProject.BuildInfo(true))
	}

	return
}

func ExecuteGo(noRegistry bool, goArg, targetRepo, buildName, buildNumber string, details *config.ArtifactoryDetails) error {
	isCollectBuildInfo := len(buildName) > 0 && len(buildNumber) > 0
	if isCollectBuildInfo {
		err := utils.SaveBuildGeneralDetails(buildName, buildNumber)
		if err != nil {
			return err
		}
	}

	if !noRegistry {
		goutils.SetGoProxyEnvVar(details, targetRepo)
	}
	err := goutils.RunGo(goArg)
	if err != nil {
		return err
	}

	if isCollectBuildInfo {
		// The version is not necessary because we are collecting the dependencies only.
		goProject, err := project.Load("-")
		if err != nil {
			return err
		}
		err = utils.SaveBuildInfo(buildName, buildNumber, goProject.BuildInfo(false))
	}

	return err
}

func validatePrerequisites() error {
	_, err := exec.LookPath("go")
	if err != nil {
		return errorutils.CheckError(err)
	}
	return nil
}
