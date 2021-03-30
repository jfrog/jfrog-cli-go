package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gookit/color"
	pipelinesservices "github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/jfrog/jfrog-client-go/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands/buildinfo"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands/generic"
	rtutils "github.com/jfrog/jfrog-cli-core/artifactory/utils"
	corecommoncommands "github.com/jfrog/jfrog-cli-core/common/commands"
	utilsconfig "github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-cli-core/utils/coreutils"
	"github.com/jfrog/jfrog-cli-core/utils/ioutils"
	buildinfocmd "github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray"
	xrayservices "github.com/jfrog/jfrog-client-go/xray/services"
	xrayutils "github.com/jfrog/jfrog-client-go/xray/services/utils"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	ConfigServerId          = "vcs-integration-platform"
	VcsConfigFile           = "jfrog-cli-vcs.conf"
	DefaultFirstBuildNumber = "0"
	DefaultWorkspace        = "./jfrog-vcs-workspace"
	pipelineUiPath          = "ui/pipelines/myPipelines/default/"
)

type GitProvider string

const (
	Github           = "GitHub"
	GithubEnterprise = "GitHub Enterprise"
	Bitbucket        = "Bitbucket"
	BitbucketServer  = "Bitbucket Server"
	Gitlab           = "GitLab"
)

type CiSetupCommand struct {
	defaultData *VcsData
	data        *VcsData
}

type VcsData struct {
	RepositoryName          string
	ProjectDomain           string
	LocalDirPath            string
	GitBranch               string
	BuildCommand            string
	BuildName               string
	ArtifactoryVirtualRepos map[Technology]string
	// A collection of technologies that was found with a list of theirs indications
	DetectedTechnologies map[Technology]bool
	VcsCredentials       VcsServerDetails
	GitProvider          GitProvider
}
type VcsServerDetails struct {
	Url         string `json:"url,omitempty"`
	User        string `json:"user,omitempty"`
	Password    string `json:"-"`
	AccessToken string `json:"-"`
}

func (cc *CiSetupCommand) SetData(data *VcsData) *CiSetupCommand {
	cc.data = data
	return cc
}
func (cc *CiSetupCommand) SetDefaultData(data *VcsData) *CiSetupCommand {
	cc.defaultData = data
	return cc
}

func RunCiSetupCmd() error {
	logBeginningInstructions()
	cc := &CiSetupCommand{}
	err := cc.prepareConfigurationData()
	if err != nil {
		return err
	}
	err = cc.Run()
	if err != nil {
		return err
	}
	return saveVcsConf(cc.data)

}

func logBeginningInstructions() {
	instructions := []string{
		"",
		colorTitle("About this command"),
		"This command sets up a basic CI pipeline which uses the JFrog Platform.",
		"It currently supports maven, gradle and npm, but additional package managers will be added in the future.",
		"The generated CI pipeline is based on JFrog Pipelines, but additional CI providers will be added in the future.",
		"",
		colorTitle("Important"),
		" 1. When asked to provide credentials for your JFrog Platform and Git provider, please make sure the credentials have admin privileges.",
		" 2. You can exit the command by hitting 'control + C' at any time. The values you provided before exiting are saved (with the exception of passwords and tokens) and will be set as defaults the next tine you run the command.",
		"",
	}
	log.Info(strings.Join(instructions, "\n"))
}

func colorTitle(title string) string {
	if terminal.IsTerminal(int(os.Stderr.Fd())) {
		return color.Green.Render(title)
	}
	return title
}

func (cc *CiSetupCommand) prepareConfigurationData() error {
	// If data is nil, initialize a new one
	if cc.data == nil {
		cc.data = new(VcsData)
	}

	// Get previous vcs data if exists
	defaultData, err := readVcsConf()
	cc.defaultData = defaultData
	return err
}

func readVcsConf() (conf *VcsData, err error) {
	conf = &VcsData{}
	path, err := coreutils.GetJfrogHomeDir()
	if err != nil {
		return
	}
	configFile, err := fileutils.ReadFile(filepath.Join(path, VcsConfigFile))
	if err != nil {
		return
	}
	err = json.Unmarshal(configFile, conf)
	return
}

func saveVcsConf(conf *VcsData) error {
	path, err := coreutils.GetJfrogHomeDir()
	if err != nil {
		return err
	}
	bytesContent, err := json.Marshal(conf)
	if err != nil {
		return errorutils.CheckError(err)
	}
	var content bytes.Buffer
	err = json.Indent(&content, bytesContent, "", "  ")
	if err != nil {
		return errorutils.CheckError(err)
	}
	err = ioutil.WriteFile(filepath.Join(path, VcsConfigFile), []byte(content.String()), 0600)
	if err != nil {
		return errorutils.CheckError(err)
	}
	return nil
}

func (cc *CiSetupCommand) Run() error {
	// Run JFrog config command
	err := runConfigCmd()
	if err != nil {
		return err
	}

	// Basic VCS questionnaire (URLs, Credentials, etc'...)
	err = cc.gitPhase()
	if err != nil || saveVcsConf(cc.data) != nil {
		return err
	}

	// Interactively create Artifactory repository based on the detected technologies and on going user input
	err = cc.artifactoryConfigPhase()
	if err != nil || saveVcsConf(cc.data) != nil {
		return err
	}
	// Publish empty build info.
	err = cc.publishFirstBuild()
	if err != nil || saveVcsConf(cc.data) != nil {
		return err
	}
	// Configure Xray to scan the new build.
	err = cc.xrayConfigPhase()
	if err != nil || saveVcsConf(cc.data) != nil {
		return err
	}
	return cc.runPipelinesPhase()
}

func getPipelinesToken() (string, error) {
	var err error
	var byteToken []byte
	for len(byteToken) == 0 {
		print("Please provide a JFrog Pipelines admin token (To generate the token, " +
			"log into the JFrog Platform UI --> Administration --> Identity and Access --> Access Tokens --> Generate Admin Token): ")
		byteToken, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", errorutils.CheckError(err)
		}
		// New-line required after the access token input:
		fmt.Println()
	}
	return string(byteToken), nil
}

func runConfigCmd() (err error) {
	for {
		configCmd := corecommoncommands.NewConfigCommand().SetInteractive(true).SetServerId(ConfigServerId).SetEncPassword(true)
		err = configCmd.Config()
		if err != nil {
			log.Error(err)
			continue
		}
		// Validate JFrog credentials by excute ping command
		serviceDetails, err := utilsconfig.GetSpecificConfig(ConfigServerId, false, false)
		if err != nil {
			return err
		}
		err = generic.NewPingCommand().SetServerDetails(serviceDetails).Run()
		if err == nil {
			return nil
		}
		log.Error(err)
	}
}

func (cc *CiSetupCommand) runPipelinesPhase() error {
	var pipelinesYamlBytes []byte
	var pipelineName string
	var err error
	// Ask for token and config pipelines. Run again if authentication problem.
	for {
		// Ask for pipelines token.
		pipelinesToken, err := getPipelinesToken()
		if err != nil {
			return err
		}
		// Run Pipelines setup
		pipelinesYamlBytes, pipelineName, err = configAndGeneratePipelines(cc.data, pipelinesToken)
		// If no error, continue with flow. Elseif unauthorized error, ask for token again.
		if err == nil {
			break
		}
		if _, ok := err.(*pipelinesservices.IntegrationUnauthorizedError); !ok {
			return err
		}
		log.Debug(err.Error())
		log.Info("There seems to be an authorization problem with the pipelines token you entered. Please try again.")
	}

	err = cc.saveYamlToFile(pipelinesYamlBytes)
	if err != nil {
		return err
	}
	err = cc.stagePipelinesYaml(pipelinesYamlPath)
	if err != nil {
		return err
	}
	return cc.logCompletionInstruction(pipelineName)
}

func (cc *CiSetupCommand) saveYamlToFile(yaml []byte) error {
	path := filepath.Join(cc.data.LocalDirPath, pipelinesYamlPath)
	log.Info("Generating pipelines.yml at: '" + path + "'...")
	return ioutil.WriteFile(path, yaml, 0644)
}

func (cc *CiSetupCommand) logCompletionInstruction(pipelineName string) error {
	serviceDetails, err := utilsconfig.GetSpecificConfig(ConfigServerId, false, false)
	if err != nil {
		return err
	}

	instructions := []string{
		"To complete the setup, run the following commands:", "",
		"cd " + cc.data.LocalDirPath,
		"git commit -m \"Add pipelines.yml\"",
		"git push", "",
		"Although your pipeline is configured, it hasn't run yet.",
		"It will run and become visible in the following URL, after the next git commit:",
		getPipelineUiPath(serviceDetails.Url, pipelineName), "",
	}
	log.Info(strings.Join(instructions, "\n"))
	return nil
}

func getPipelineUiPath(pipelinesUrl, pipelineName string) string {
	return utils.AddTrailingSlashIfNeeded(pipelinesUrl) + pipelineUiPath + pipelineName
}

func (cc *CiSetupCommand) publishFirstBuild() (err error) {
	println("Everytime the new pipeline builds the code, it generates a build entity (also known as build-info) and stores it in Artifactory.")
	ioutils.ScanFromConsole("Please choose a name for the build", &cc.data.BuildName, "${vcs.repo.name}-${branch}")
	cc.data.BuildName = strings.Replace(cc.data.BuildName, "${vcs.repo.name}", cc.data.RepositoryName, -1)
	cc.data.BuildName = strings.Replace(cc.data.BuildName, "${branch}", cc.data.GitBranch, -1)
	// Run BAG Command (in order to publish the first, empty, build info)
	buildAddGitConfigurationCmd := buildinfo.NewBuildAddGitCommand().SetDotGitPath(cc.data.LocalDirPath).SetServerId(ConfigServerId) //.SetConfigFilePath(c.String("config"))
	buildConfiguration := rtutils.BuildConfiguration{BuildName: cc.data.BuildName, BuildNumber: DefaultFirstBuildNumber}
	buildAddGitConfigurationCmd = buildAddGitConfigurationCmd.SetBuildConfiguration(&buildConfiguration)
	log.Info("Generating an initial build-info...")
	err = commands.Exec(buildAddGitConfigurationCmd)
	if err != nil {
		return err
	}
	// Run BP Command.
	serviceDetails, err := utilsconfig.GetSpecificConfig(ConfigServerId, false, false)
	if err != nil {
		return err
	}
	buildInfoConfiguration := buildinfocmd.Configuration{DryRun: false}
	buildPublishCmd := buildinfo.NewBuildPublishCommand().SetServerDetails(serviceDetails).SetBuildConfiguration(&buildConfiguration).SetConfig(&buildInfoConfiguration)
	err = commands.Exec(buildPublishCmd)
	if err != nil {
		return err

	}
	return
}

func (cc *CiSetupCommand) xrayConfigPhase() (err error) {
	serviceDetails, err := utilsconfig.GetSpecificConfig(ConfigServerId, false, false)
	if err != nil {
		return err
	}
	xrayDetails, err := serviceDetails.CreateXrayAuthConfig()
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(xrayDetails).
		Build()
	if err != nil {
		return err
	}
	xrayManager, err := xray.New(&xrayDetails, serviceConfig)
	if err != nil {
		return err
	}
	// AddBuildsToIndexing.
	buildsToIndex := []string{cc.data.BuildName}
	err = xrayManager.AddBuildsToIndexing(buildsToIndex)
	// Create new default policy.
	policyParams := xrayutils.NewPolicyParams()
	policyParams.Name = "vcs-integration-security-policy"
	policyParams.Type = xrayutils.Security
	policyParams.Description = "Basic Security policy."
	policyParams.Rules = []xrayutils.PolicyRule{
		{
			Name:     "min-severity-rule",
			Criteria: *xrayutils.CreateSeverityPolicyCriteria(xrayutils.Critical),
			Priority: 1,
		},
	}
	err = xrayManager.CreatePolicy(policyParams)
	if err != nil {
		// In case the error is from type PolicyAlreadyExistsError, we should continue with the regular flow.
		if _, ok := err.(*xrayservices.PolicyAlreadyExistsError); !ok {
			return err
		} else {
			log.Debug(err.(*xrayservices.PolicyAlreadyExistsError).InnerError)
			err = nil
		}
	}
	// Create new default watcher.
	watchParams := xrayutils.NewWatchParams()
	watchParams.Name = "vcs-integration-watch-all"
	watchParams.Description = "VCS Configured Build Watch"
	watchParams.Active = true

	// Need to be verified before merging
	watchParams.Builds.Type = xrayutils.WatchBuildAll
	watchParams.Policies = []xrayutils.AssignedPolicy{
		{
			Name: policyParams.Name,
			Type: "security",
		},
	}

	err = xrayManager.CreateWatch(watchParams)
	if err != nil {
		// In case the error is from type WatchAlreadyExistsError, we should continue with the regular flow.
		if _, ok := err.(*xrayservices.WatchAlreadyExistsError); !ok {
			return err
		} else {
			log.Debug(err.(*xrayservices.WatchAlreadyExistsError).InnerError)
			err = nil
		}
	}
	return
}

func (cc *CiSetupCommand) artifactoryConfigPhase() (err error) {

	cc.data.ArtifactoryVirtualRepos = make(map[Technology]string)
	// First create repositories for each technology in Artifactory according to user input
	for tech, detected := range cc.data.DetectedTechnologies {
		if detected && coreutils.AskYesNo(fmt.Sprintf("It looks like the source code is built using %s. Would you like to resolve the %s dependencies from Artifactory?", tech, tech), true) {
			err = cc.interactivelyCreatRepos(tech)
			if err != nil {
				return
			}
		}
	}
	// Ask for working build command
	prompt := "Please provide a single-line build command. You may use the && operator. Currently scripts (such as bash scripts) are not supported"
	ioutils.ScanFromConsole(prompt, &cc.data.BuildCommand, cc.defaultData.BuildCommand)
	return nil
}

func (cc *CiSetupCommand) interactivelyCreatRepos(technologyType Technology) (err error) {
	serviceDetails, err := utilsconfig.GetSpecificConfig(ConfigServerId, false, false)
	if err != nil {
		return err
	}
	// Get all relevant remotes to choose from
	remoteRepos, err := GetAllRepos(serviceDetails, Remote, string(technologyType))
	if err != nil {
		return err
	}

	// Ask if the user would like us to create a new remote or to choose from the exist repositories list
	remoteRepo, err := promptARepoSelection(remoteRepos, "Select remote repository")
	if err != nil {
		return nil
	}
	// The user choose to create a new remote repo
	if remoteRepo == NewRepository {
		for {
			var repoName, repoUrl string
			ioutils.ScanFromConsole("Repository Name", &repoName, GetRemoteDefaultName(technologyType))
			ioutils.ScanFromConsole("Repository URL", &repoUrl, GetRemoteDefaultUrl(technologyType))
			err = CreateRemoteRepo(serviceDetails, technologyType, repoName, repoUrl)
			if err != nil {
				log.Error(err)
			} else {
				remoteRepo = repoName
				for {
					// Create a new virtual repository as well
					ioutils.ScanFromConsole(fmt.Sprintf("Choose a name for a new virtual repository which will include %q remote repo", remoteRepo),
						&repoName, GetVirtualDefaultName(technologyType))
					err = CreateVirtualRepo(serviceDetails, technologyType, repoName, remoteRepo)
					if err != nil {
						log.Error(err)
					} else {
						// we created both remote and virtual repositories successfully
						cc.data.ArtifactoryVirtualRepos[technologyType] = repoName
						return
					}
				}
			}
		}
	}
	// Else, the user choose an existing remote repo
	virtualRepos, err := GetAllRepos(serviceDetails, Virtual, string(technologyType))
	if err != nil {
		return err
	}
	// Ask if the user would like us to create a new virtual or to choose from the exist repositories list
	virtualRepo, err := promptARepoSelection(virtualRepos, fmt.Sprintf("Select a virtual repository, which includes %s or choose to create a new repo:", remoteRepo))
	if virtualRepo == NewRepository {
		// Create virtual repository
		for {
			var repoName string
			ioutils.ScanFromConsole("Repository Name", &repoName, GetVirtualDefaultName(technologyType))
			err = CreateVirtualRepo(serviceDetails, technologyType, repoName, remoteRepo)
			if err != nil {
				log.Error(err)
			} else {
				virtualRepo = repoName
				break
			}
		}
	} else {
		// Validate that the chosen virtual repo contains the chosen remote repo
		chosenVirtualRepo, err := GetVirtualRepo(serviceDetails, virtualRepo)
		if err != nil {
			return err
		}
		if !contains(chosenVirtualRepo.Repositories, remoteRepo) {
			log.Error(fmt.Sprintf("The chosen virtual repo %q does not contain the chosen remote repo %q", virtualRepo, remoteRepo))
			return cc.interactivelyCreatRepos(technologyType)
		}
	}
	// Saves the new created repo name (key) in the results data structure.
	cc.data.ArtifactoryVirtualRepos[technologyType] = virtualRepo
	return
}

func promptARepoSelection(repoDetails *[]services.RepositoryDetails, promptMsg string) (selectedRepoName string, err error) {

	selectableItems := []ioutils.PromptItem{{Option: NewRepository, TargetValue: &selectedRepoName}}
	for _, repo := range *repoDetails {
		selectableItems = append(selectableItems, ioutils.PromptItem{Option: repo.Key, TargetValue: &selectedRepoName, DefaultValue: repo.Url})
	}
	println(promptMsg)
	err = ioutils.SelectString(selectableItems, "", func(item ioutils.PromptItem) {
		*item.TargetValue = item.Option
	})
	return
}

func promptGitProviderSelection() (selected string, err error) {
	gitProviders := []GitProvider{
		Github,
		GithubEnterprise,
		Bitbucket,
		BitbucketServer,
		Gitlab,
	}

	var selectableItems []ioutils.PromptItem
	for _, provider := range gitProviders {
		selectableItems = append(selectableItems, ioutils.PromptItem{Option: string(provider), TargetValue: &selected})
	}
	println("Choose your project Git provider:")
	err = ioutils.SelectString(selectableItems, "", func(item ioutils.PromptItem) {
		*item.TargetValue = item.Option
	})
	return
}

func (cc *CiSetupCommand) prepareVcsData() (err error) {
	cc.data.LocalDirPath = DefaultWorkspace
	for {
		err = fileutils.CreateDirIfNotExist(cc.data.LocalDirPath)
		if err != nil {
			return err
		}
		dirEmpty, err := fileutils.IsDirEmpty(cc.data.LocalDirPath)
		if err != nil {
			return err
		}
		if dirEmpty {
			break
		} else {
			log.Error("The '" + cc.data.LocalDirPath + "' directory isn't empty.")
			ioutils.ScanFromConsole("Choose a name for a directory to be used as the command's workspace", &cc.data.LocalDirPath, "")
		}

	}
	err = cc.cloneProject()
	if err != nil {
		return
	}
	err = cc.detectTechnologies()
	return
}

func (cc *CiSetupCommand) cloneProject() (err error) {
	// Create the desired path if necessary
	err = os.MkdirAll(cc.data.LocalDirPath, os.ModePerm)
	if err != nil {
		return err
	}
	cloneOption := &git.CloneOptions{
		URL:           cc.data.VcsCredentials.Url,
		Auth:          createCredentials(&cc.data.VcsCredentials),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", cc.data.GitBranch)),
		// Enable git submodules clone if there any.
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	cc.extractRepositoryName()
	// Clone the given repository to the given directory from the given branch
	log.Info(fmt.Sprintf("Cloning project %q from: %q into: %q", cc.data.RepositoryName, cc.data.VcsCredentials.Url, cc.data.LocalDirPath))
	_, err = git.PlainClone(cc.data.LocalDirPath, false, cloneOption)
	if err != nil {
		return err
	}
	return
}

func (cc *CiSetupCommand) stagePipelinesYaml(path string) error {
	log.Info("Staging pipelines.yml for git commit...")
	repo, err := git.PlainOpen(cc.data.LocalDirPath)
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	_, err = worktree.Add(path)
	return err
}

func (cc *CiSetupCommand) extractRepositoryName() {
	vcsUrl := cc.data.VcsCredentials.Url
	if vcsUrl == "" {
		return
	}
	// Trim trailing "/" if one exists
	vcsUrl = strings.TrimSuffix(vcsUrl, "/")
	cc.data.VcsCredentials.Url = vcsUrl
	splitUrl := strings.Split(vcsUrl, "/")
	repositoryName := splitUrl[len(splitUrl)-1]
	cc.data.ProjectDomain = splitUrl[len(splitUrl)-2]
	cc.data.RepositoryName = strings.TrimSuffix(repositoryName, ".git")
}

func (cc *CiSetupCommand) detectTechnologies() (err error) {
	indicators := GetTechIndicators()
	filesList, err := fileutils.ListFilesRecursiveWalkIntoDirSymlink(cc.data.LocalDirPath, false)
	if err != nil {
		return err
	}
	cc.data.DetectedTechnologies = make(map[Technology]bool)
	for _, file := range filesList {
		for _, indicator := range indicators {
			if indicator.Indicates(file) {
				cc.data.DetectedTechnologies[indicator.GetTechnology()] = true
				// Same file can't indicate on more than one technology.
				break
			}
		}
	}
	return
}

func createCredentials(serviceDetails *VcsServerDetails) (auth transport.AuthMethod) {
	var password, username string
	if serviceDetails.AccessToken != "" {
		password = serviceDetails.AccessToken
		// Authentication fails if the username string is empty. This can be anything except an empty string...
		username = "user"
	} else {
		password = serviceDetails.Password
		username = serviceDetails.User
	}
	return &http.BasicAuth{Username: username, Password: password}
}

func (cc *CiSetupCommand) gitPhase() (err error) {
	for {
		gitProvider, err := promptGitProviderSelection()
		if err != nil {
			log.Error(err)
			continue
		}
		cc.data.GitProvider = GitProvider(gitProvider)
		ioutils.ScanFromConsole("Git project URL", &cc.data.VcsCredentials.Url, cc.defaultData.VcsCredentials.Url)
		ioutils.ScanFromConsole("Git username", &cc.data.VcsCredentials.User, cc.defaultData.VcsCredentials.User)
		print("Git access token: ")
		byteToken, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Error(err)
			continue
		}
		// New-line required after the access token input:
		fmt.Println()
		cc.data.VcsCredentials.AccessToken = string(byteToken)
		ioutils.ScanFromConsole("Git branch", &cc.data.GitBranch, cc.defaultData.GitBranch)
		err = cc.prepareVcsData()
		if err != nil {
			log.Error(err)
		} else {
			return nil
		}
	}

}
