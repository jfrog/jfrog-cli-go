package artifactory

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/yarn"
	"github.com/jfrog/jfrog-cli-core/v2/common/spec"
	mvndoc "github.com/jfrog/jfrog-cli/docs/artifactory/mvn"
	yarndocs "github.com/jfrog/jfrog-cli/docs/artifactory/yarn"
	"github.com/jfrog/jfrog-cli/docs/artifactory/yarnconfig"

	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/container"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/dotnet"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/permissiontarget"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/usersmanagement"
	commandsutils "github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/utils"
	containerutils "github.com/jfrog/jfrog-cli-core/v2/artifactory/utils/container"
	coreCommonCommands "github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	"github.com/jfrog/jfrog-cli-core/v2/utils/ioutils"
	"github.com/jfrog/jfrog-cli/docs/artifactory/accesstokencreate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/builddockercreate"
	dotnetdocs "github.com/jfrog/jfrog-cli/docs/artifactory/dotnet"
	"github.com/jfrog/jfrog-cli/docs/artifactory/dotnetconfig"
	"github.com/jfrog/jfrog-cli/docs/artifactory/groupaddusers"
	"github.com/jfrog/jfrog-cli/docs/artifactory/groupcreate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/groupdelete"
	"github.com/jfrog/jfrog-cli/docs/artifactory/permissiontargetcreate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/permissiontargetdelete"
	"github.com/jfrog/jfrog-cli/docs/artifactory/permissiontargettemplate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/permissiontargetupdate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/podmanpull"
	"github.com/jfrog/jfrog-cli/docs/artifactory/podmanpush"
	"github.com/jfrog/jfrog-cli/docs/artifactory/usercreate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/userscreate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/usersdelete"
	logUtils "github.com/jfrog/jfrog-cli/utils/log"
	"github.com/jfrog/jfrog-cli/utils/progressbar"
	ioUtils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jszwec/csvutil"

	"github.com/codegangsta/cli"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/buildinfo"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/curl"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/generic"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/golang"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/gradle"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/mvn"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/npm"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/pip"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/replication"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/repository"
	commandUtils "github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/utils"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	npmUtils "github.com/jfrog/jfrog-cli-core/v2/artifactory/utils/npm"
	"github.com/jfrog/jfrog-cli-core/v2/common/commands"
	corecommon "github.com/jfrog/jfrog-cli-core/v2/docs/common"
	coreConfig "github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli/docs/artifactory/buildadddependencies"
	"github.com/jfrog/jfrog-cli/docs/artifactory/buildaddgit"
	"github.com/jfrog/jfrog-cli/docs/artifactory/buildappend"
	"github.com/jfrog/jfrog-cli/docs/artifactory/buildclean"
	"github.com/jfrog/jfrog-cli/docs/artifactory/buildcollectenv"
	"github.com/jfrog/jfrog-cli/docs/artifactory/builddiscard"
	"github.com/jfrog/jfrog-cli/docs/artifactory/buildpromote"
	"github.com/jfrog/jfrog-cli/docs/artifactory/buildpublish"
	"github.com/jfrog/jfrog-cli/docs/artifactory/buildscan"
	copydocs "github.com/jfrog/jfrog-cli/docs/artifactory/copy"
	curldocs "github.com/jfrog/jfrog-cli/docs/artifactory/curl"
	"github.com/jfrog/jfrog-cli/docs/artifactory/delete"
	"github.com/jfrog/jfrog-cli/docs/artifactory/deleteprops"
	"github.com/jfrog/jfrog-cli/docs/artifactory/dockerpromote"
	"github.com/jfrog/jfrog-cli/docs/artifactory/dockerpull"
	"github.com/jfrog/jfrog-cli/docs/artifactory/dockerpush"
	"github.com/jfrog/jfrog-cli/docs/artifactory/download"
	"github.com/jfrog/jfrog-cli/docs/artifactory/gitlfsclean"
	"github.com/jfrog/jfrog-cli/docs/artifactory/gocommand"
	"github.com/jfrog/jfrog-cli/docs/artifactory/goconfig"
	"github.com/jfrog/jfrog-cli/docs/artifactory/gopublish"
	gradledoc "github.com/jfrog/jfrog-cli/docs/artifactory/gradle"
	"github.com/jfrog/jfrog-cli/docs/artifactory/gradleconfig"
	"github.com/jfrog/jfrog-cli/docs/artifactory/move"
	"github.com/jfrog/jfrog-cli/docs/artifactory/mvnconfig"
	"github.com/jfrog/jfrog-cli/docs/artifactory/npmci"
	"github.com/jfrog/jfrog-cli/docs/artifactory/npmconfig"
	"github.com/jfrog/jfrog-cli/docs/artifactory/npminstall"
	"github.com/jfrog/jfrog-cli/docs/artifactory/npmpublish"
	nugetdocs "github.com/jfrog/jfrog-cli/docs/artifactory/nuget"
	"github.com/jfrog/jfrog-cli/docs/artifactory/nugetconfig"
	nugettree "github.com/jfrog/jfrog-cli/docs/artifactory/nugetdepstree"
	"github.com/jfrog/jfrog-cli/docs/artifactory/ping"
	"github.com/jfrog/jfrog-cli/docs/artifactory/pipconfig"
	"github.com/jfrog/jfrog-cli/docs/artifactory/pipinstall"
	"github.com/jfrog/jfrog-cli/docs/artifactory/replicationcreate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/replicationdelete"
	"github.com/jfrog/jfrog-cli/docs/artifactory/replicationtemplate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/repocreate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/repodelete"
	"github.com/jfrog/jfrog-cli/docs/artifactory/repotemplate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/repoupdate"
	"github.com/jfrog/jfrog-cli/docs/artifactory/search"
	"github.com/jfrog/jfrog-cli/docs/artifactory/setprops"
	"github.com/jfrog/jfrog-cli/docs/artifactory/upload"
	"github.com/jfrog/jfrog-cli/docs/common"
	"github.com/jfrog/jfrog-cli/utils/cliutils"
	buildinfocmd "github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func GetCommands() []cli.Command {
	return cliutils.GetSortedCommands(cli.CommandsByName{
		{
			Name:         "upload",
			Flags:        cliutils.GetCommandFlags(cliutils.Upload),
			Aliases:      []string{"u"},
			Description:  upload.Description,
			HelpName:     corecommon.CreateUsage("rt upload", upload.Description, upload.Usage),
			UsageText:    upload.Arguments,
			ArgsUsage:    common.CreateEnvVars(upload.EnvVar),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return uploadCmd(c)
			},
		},
		{
			Name:         "download",
			Flags:        cliutils.GetCommandFlags(cliutils.Download),
			Aliases:      []string{"dl"},
			Description:  download.Description,
			HelpName:     corecommon.CreateUsage("rt download", download.Description, download.Usage),
			UsageText:    download.Arguments,
			ArgsUsage:    common.CreateEnvVars(download.EnvVar),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return downloadCmd(c)
			},
		},
		{
			Name:         "move",
			Flags:        cliutils.GetCommandFlags(cliutils.Move),
			Aliases:      []string{"mv"},
			Description:  move.Description,
			HelpName:     corecommon.CreateUsage("rt move", move.Description, move.Usage),
			UsageText:    move.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return moveCmd(c)
			},
		},
		{
			Name:         "copy",
			Flags:        cliutils.GetCommandFlags(cliutils.Copy),
			Aliases:      []string{"cp"},
			Description:  copydocs.Description,
			HelpName:     corecommon.CreateUsage("rt copy", copydocs.Description, copydocs.Usage),
			UsageText:    copydocs.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return copyCmd(c)
			},
		},
		{
			Name:         "delete",
			Flags:        cliutils.GetCommandFlags(cliutils.Delete),
			Aliases:      []string{"del"},
			Description:  delete.Description,
			HelpName:     corecommon.CreateUsage("rt delete", delete.Description, delete.Usage),
			UsageText:    delete.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return deleteCmd(c)
			},
		},
		{
			Name:         "search",
			Flags:        cliutils.GetCommandFlags(cliutils.Search),
			Aliases:      []string{"s"},
			Description:  search.Description,
			HelpName:     corecommon.CreateUsage("rt search", search.Description, search.Usage),
			UsageText:    search.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return searchCmd(c)
			},
		},
		{
			Name:         "set-props",
			Flags:        cliutils.GetCommandFlags(cliutils.Properties),
			Aliases:      []string{"sp"},
			Description:  setprops.Description,
			HelpName:     corecommon.CreateUsage("rt set-props", setprops.Description, setprops.Usage),
			UsageText:    setprops.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return setPropsCmd(c)
			},
		},
		{
			Name:         "delete-props",
			Flags:        cliutils.GetCommandFlags(cliutils.Properties),
			Aliases:      []string{"delp"},
			Description:  deleteprops.Description,
			HelpName:     corecommon.CreateUsage("rt delete-props", deleteprops.Description, deleteprops.Usage),
			UsageText:    deleteprops.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return deletePropsCmd(c)
			},
		},
		{
			Name:         "build-publish",
			Flags:        cliutils.GetCommandFlags(cliutils.BuildPublish),
			Aliases:      []string{"bp"},
			Description:  buildpublish.Description,
			HelpName:     corecommon.CreateUsage("rt build-publish", buildpublish.Description, buildpublish.Usage),
			UsageText:    buildpublish.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildPublishCmd(c)
			},
		},
		{
			Name:         "build-collect-env",
			Aliases:      []string{"bce"},
			Flags:        cliutils.GetCommandFlags(cliutils.BuildCollectEnv),
			Description:  buildcollectenv.Description,
			HelpName:     corecommon.CreateUsage("rt build-collect-env", buildcollectenv.Description, buildcollectenv.Usage),
			UsageText:    buildcollectenv.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildCollectEnvCmd(c)
			},
		},
		{
			Name:         "build-append",
			Flags:        cliutils.GetCommandFlags(cliutils.BuildAppend),
			Aliases:      []string{"ba"},
			Description:  buildappend.Description,
			HelpName:     corecommon.CreateUsage("rt build-append", buildappend.Description, buildappend.Usage),
			UsageText:    buildappend.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildAppendCmd(c)
			},
		},
		{
			Name:         "build-add-dependencies",
			Flags:        cliutils.GetCommandFlags(cliutils.BuildAddDependencies),
			Aliases:      []string{"bad"},
			Description:  buildadddependencies.Description,
			HelpName:     corecommon.CreateUsage("rt build-add-dependencies", buildadddependencies.Description, buildadddependencies.Usage),
			UsageText:    buildadddependencies.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildAddDependenciesCmd(c)
			},
		},
		{
			Name:         "build-add-git",
			Flags:        cliutils.GetCommandFlags(cliutils.BuildAddGit),
			Aliases:      []string{"bag"},
			Description:  buildaddgit.Description,
			HelpName:     corecommon.CreateUsage("rt build-add-git", buildaddgit.Description, buildaddgit.Usage),
			UsageText:    buildaddgit.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildAddGitCmd(c)
			},
		},
		{
			Name:         "build-scan",
			Flags:        cliutils.GetCommandFlags(cliutils.BuildScan),
			Aliases:      []string{"bs"},
			Description:  buildscan.Description,
			HelpName:     corecommon.CreateUsage("rt build-scan", buildscan.Description, buildscan.Usage),
			UsageText:    buildscan.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildScanCmd(c)
			},
		},
		{
			Name:         "build-clean",
			Aliases:      []string{"bc"},
			Description:  buildclean.Description,
			HelpName:     corecommon.CreateUsage("rt build-clean", buildclean.Description, buildclean.Usage),
			UsageText:    buildclean.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildCleanCmd(c)
			},
		},
		{
			Name:         "build-promote",
			Flags:        cliutils.GetCommandFlags(cliutils.BuildPromote),
			Aliases:      []string{"bpr"},
			Description:  buildpromote.Description,
			HelpName:     corecommon.CreateUsage("rt build-promote", buildpromote.Description, buildpromote.Usage),
			UsageText:    buildpromote.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildPromoteCmd(c)
			},
		},
		{
			Name:         "build-discard",
			Flags:        cliutils.GetCommandFlags(cliutils.BuildDiscard),
			Aliases:      []string{"bdi"},
			Description:  builddiscard.Description,
			HelpName:     corecommon.CreateUsage("rt build-discard", builddiscard.Description, builddiscard.Usage),
			UsageText:    builddiscard.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return buildDiscardCmd(c)
			},
		},
		{
			Name:         "git-lfs-clean",
			Flags:        cliutils.GetCommandFlags(cliutils.GitLfsClean),
			Aliases:      []string{"glc"},
			Description:  gitlfsclean.Description,
			HelpName:     corecommon.CreateUsage("rt git-lfs-clean", gitlfsclean.Description, gitlfsclean.Usage),
			UsageText:    gitlfsclean.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return gitLfsCleanCmd(c)
			},
		},
		{
			Name:         "mvn-config",
			Aliases:      []string{"mvnc"},
			Flags:        cliutils.GetCommandFlags(cliutils.MvnConfig),
			Description:  mvnconfig.Description,
			HelpName:     corecommon.CreateUsage("rt mvn-config", mvnconfig.Description, mvnconfig.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createMvnConfigCmd(c)
			},
		},
		{
			Name:            "mvn",
			Flags:           cliutils.GetCommandFlags(cliutils.Mvn),
			Description:     mvndoc.Description,
			HelpName:        corecommon.CreateUsage("rt mvn", mvndoc.Description, mvndoc.Usage),
			UsageText:       mvndoc.Arguments,
			ArgsUsage:       common.CreateEnvVars(mvndoc.EnvVar),
			SkipFlagParsing: shouldSkipMavenFlagParsing(),
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return mvnCmd(c)
			},
		},
		{
			Name:         "gradle-config",
			Aliases:      []string{"gradlec"},
			Flags:        cliutils.GetCommandFlags(cliutils.GradleConfig),
			Description:  gradleconfig.Description,
			HelpName:     corecommon.CreateUsage("rt gradle-config", gradleconfig.Description, gradleconfig.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createGradleConfigCmd(c)
			},
		},
		{
			Name:            "gradle",
			Flags:           cliutils.GetCommandFlags(cliutils.Gradle),
			Description:     gradledoc.Description,
			HelpName:        corecommon.CreateUsage("rt gradle", gradledoc.Description, gradledoc.Usage),
			UsageText:       gradledoc.Arguments,
			ArgsUsage:       common.CreateEnvVars(gradledoc.EnvVar),
			SkipFlagParsing: shouldSkipGradleFlagParsing(),
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return gradleCmd(c)
			},
		},
		{
			Name:         "docker-promote",
			Flags:        cliutils.GetCommandFlags(cliutils.DockerPromote),
			Aliases:      []string{"dpr"},
			Description:  dockerpromote.Description,
			HelpName:     corecommon.CreateUsage("rt docker-promote", dockerpromote.Description, dockerpromote.Usage),
			UsageText:    dockerpromote.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return dockerPromoteCmd(c)
			},
		},
		{
			Name:         "docker-push",
			Flags:        cliutils.GetCommandFlags(cliutils.ContainerPush),
			Aliases:      []string{"dp"},
			Description:  dockerpush.Description,
			HelpName:     corecommon.CreateUsage("rt docker-push", dockerpush.Description, dockerpush.Usage),
			UsageText:    dockerpush.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return containerPushCmd(c, containerutils.DockerClient)
			},
		},
		{
			Name:         "docker-pull",
			Flags:        cliutils.GetCommandFlags(cliutils.ContainerPull),
			Aliases:      []string{"dpl"},
			Description:  dockerpull.Description,
			HelpName:     corecommon.CreateUsage("rt docker-pull", dockerpull.Description, dockerpull.Usage),
			UsageText:    dockerpull.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return containerPullCmd(c, containerutils.DockerClient)
			},
		},
		{
			Name:         "podman-push",
			Flags:        cliutils.GetCommandFlags(cliutils.ContainerPush),
			Aliases:      []string{"pp"},
			Description:  podmanpush.Description,
			HelpName:     corecommon.CreateUsage("rt podman-push", podmanpush.Description, podmanpush.Usage),
			UsageText:    podmanpush.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return containerPushCmd(c, containerutils.Podman)
			},
		},
		{
			Name:         "podman-pull",
			Flags:        cliutils.GetCommandFlags(cliutils.ContainerPull),
			Aliases:      []string{"ppl"},
			Description:  podmanpull.Description,
			HelpName:     corecommon.CreateUsage("rt podman-pull", podmanpull.Description, podmanpull.Usage),
			UsageText:    podmanpull.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return containerPullCmd(c, containerutils.Podman)
			},
		},
		{
			Name:         "build-docker-create",
			Flags:        cliutils.GetCommandFlags(cliutils.BuildDockerCreate),
			Aliases:      []string{"bdc"},
			Description:  builddockercreate.Description,
			HelpName:     corecommon.CreateUsage("rt build-docker-create", builddockercreate.Description, builddockercreate.Usage),
			UsageText:    builddockercreate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return BuildDockerCreateCmd(c)
			},
		},
		{
			Name:         "npm-config",
			Flags:        cliutils.GetCommandFlags(cliutils.NpmConfig),
			Aliases:      []string{"npmc"},
			Description:  npmconfig.Description,
			HelpName:     corecommon.CreateUsage("rt npm-config", npmconfig.Description, npmconfig.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createNpmConfigCmd(c)
			},
		},
		{
			Name:            "npm-install",
			Flags:           cliutils.GetCommandFlags(cliutils.Npm),
			Aliases:         []string{"npmi"},
			Description:     npminstall.Description,
			HelpName:        corecommon.CreateUsage("rt npm-install", npminstall.Description, npminstall.Usage),
			UsageText:       npminstall.Arguments,
			ArgsUsage:       common.CreateEnvVars(),
			SkipFlagParsing: shouldSkipNpmFlagParsing(),
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return npmInstallOrCiCmd(c)
			},
		},
		{
			Name:            "npm-ci",
			Flags:           cliutils.GetCommandFlags(cliutils.Npm),
			Aliases:         []string{"npmci"},
			Description:     npmci.Description,
			HelpName:        corecommon.CreateUsage("rt npm-ci", npmci.Description, npmci.Usage),
			UsageText:       npmci.Arguments,
			ArgsUsage:       common.CreateEnvVars(),
			SkipFlagParsing: shouldSkipNpmFlagParsing(),
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return npmInstallOrCiCmd(c)
			},
		},
		{
			Name:            "npm-publish",
			Flags:           cliutils.GetCommandFlags(cliutils.NpmPublish),
			Aliases:         []string{"npmp"},
			Description:     npmpublish.Description,
			HelpName:        corecommon.CreateUsage("rt npm-publish", npmpublish.Description, npmpublish.Usage),
			ArgsUsage:       common.CreateEnvVars(),
			SkipFlagParsing: shouldSkipNpmFlagParsing(),
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return npmPublishCmd(c)
			},
		},
		{
			Name:         "yarn-config",
			Aliases:      []string{"yarnc"},
			Flags:        cliutils.GetCommandFlags(cliutils.YarnConfig),
			Description:  yarnconfig.Description,
			HelpName:     corecommon.CreateUsage("rt yarn-config", yarnconfig.Description, yarnconfig.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createYarnConfigCmd(c)
			},
		},
		{
			Name:            "yarn",
			Flags:           cliutils.GetCommandFlags(cliutils.Yarn),
			Description:     yarndocs.Description,
			HelpName:        corecommon.CreateUsage("rt yarn", yarndocs.Description, yarndocs.Usage),
			ArgsUsage:       common.CreateEnvVars(),
			SkipFlagParsing: true,
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return yarnCmd(c)
			},
		},
		{
			Name:         "nuget-config",
			Flags:        cliutils.GetCommandFlags(cliutils.NugetConfig),
			Aliases:      []string{"nugetc"},
			Description:  nugetconfig.Description,
			HelpName:     corecommon.CreateUsage("rt nuget-config", nugetconfig.Description, nugetconfig.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createNugetConfigCmd(c)
			},
		},
		{
			Name:            "nuget",
			Flags:           cliutils.GetCommandFlags(cliutils.Nuget),
			Description:     nugetdocs.Description,
			HelpName:        corecommon.CreateUsage("rt nuget", nugetdocs.Description, nugetdocs.Usage),
			UsageText:       nugetdocs.Arguments,
			ArgsUsage:       common.CreateEnvVars(),
			SkipFlagParsing: shouldSkipNugetFlagParsing(),
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return nugetCmd(c)
			},
		},
		{
			Name:         "nuget-deps-tree",
			Aliases:      []string{"ndt"},
			Description:  nugettree.Description,
			HelpName:     corecommon.CreateUsage("rt nuget-deps-tree", nugettree.Description, nugettree.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return nugetDepsTreeCmd(c)
			},
		},
		{
			Name:         "dotnet-config",
			Flags:        cliutils.GetCommandFlags(cliutils.DotnetConfig),
			Aliases:      []string{"dotnetc"},
			Description:  dotnetconfig.Description,
			HelpName:     corecommon.CreateUsage("rt dotnet-config", dotnetconfig.Description, dotnetconfig.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createDotnetConfigCmd(c)
			},
		},
		{
			Name:            "dotnet",
			Flags:           cliutils.GetCommandFlags(cliutils.Dotnet),
			Description:     dotnetdocs.Description,
			HelpName:        corecommon.CreateUsage("rt dotnet", dotnetdocs.Description, dotnetdocs.Usage),
			UsageText:       dotnetdocs.Arguments,
			ArgsUsage:       common.CreateEnvVars(),
			SkipFlagParsing: true,
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return dotnetCmd(c)
			},
		},
		{
			Name:         "go-config",
			Flags:        cliutils.GetCommandFlags(cliutils.GoConfig),
			Description:  goconfig.Description,
			HelpName:     corecommon.CreateUsage("rt go-config", goconfig.Description, goconfig.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createGoConfigCmd(c)
			},
		},
		{
			Name:         "go-publish",
			Flags:        cliutils.GetCommandFlags(cliutils.GoPublish),
			Aliases:      []string{"gp"},
			Description:  gopublish.Description,
			HelpName:     corecommon.CreateUsage("rt go-publish", gopublish.Description, gopublish.Usage),
			UsageText:    gopublish.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return goPublishCmd(c)
			},
		},
		{
			Name:            "go",
			Flags:           cliutils.GetCommandFlags(cliutils.Go),
			Aliases:         []string{"go"},
			Description:     gocommand.Description,
			HelpName:        corecommon.CreateUsage("rt go", gocommand.Description, gocommand.Usage),
			UsageText:       gocommand.Arguments,
			ArgsUsage:       common.CreateEnvVars(),
			SkipFlagParsing: shouldSkipGoFlagParsing(),
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return goCmd(c)
			},
		},
		{
			Name:         "ping",
			Flags:        cliutils.GetCommandFlags(cliutils.Ping),
			Aliases:      []string{"p"},
			Description:  ping.Description,
			HelpName:     corecommon.CreateUsage("rt ping", ping.Description, ping.Usage),
			UsageText:    ping.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return pingCmd(c)
			},
		},
		{
			Name:            "curl",
			Flags:           cliutils.GetCommandFlags(cliutils.RtCurl),
			Aliases:         []string{"cl"},
			Description:     curldocs.Description,
			HelpName:        corecommon.CreateUsage("rt curl", curldocs.Description, curldocs.Usage),
			UsageText:       curldocs.Arguments,
			ArgsUsage:       common.CreateEnvVars(),
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			SkipFlagParsing: true,
			Action: func(c *cli.Context) error {
				return curlCmd(c)
			},
		},
		{
			Name:         "pip-config",
			Flags:        cliutils.GetCommandFlags(cliutils.PipConfig),
			Aliases:      []string{"pipc"},
			Description:  pipconfig.Description,
			HelpName:     corecommon.CreateUsage("rt pipc", pipconfig.Description, pipconfig.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createPipConfigCmd(c)
			},
		},
		{
			Name:            "pip-install",
			Flags:           cliutils.GetCommandFlags(cliutils.PipInstall),
			Aliases:         []string{"pipi"},
			Description:     pipinstall.Description,
			HelpName:        corecommon.CreateUsage("rt pipi", pipinstall.Description, pipinstall.Usage),
			UsageText:       pipinstall.Arguments,
			ArgsUsage:       common.CreateEnvVars(),
			SkipFlagParsing: true,
			BashComplete:    corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return pipInstallCmd(c)
			},
		},
		{
			Name:         "repo-template",
			Aliases:      []string{"rpt"},
			Description:  repotemplate.Description,
			HelpName:     corecommon.CreateUsage("rt rpt", repotemplate.Description, repotemplate.Usage),
			UsageText:    repotemplate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return repoTemplateCmd(c)
			},
		},
		{
			Name:         "repo-create",
			Aliases:      []string{"rc"},
			Flags:        cliutils.GetCommandFlags(cliutils.TemplateConsumer),
			Description:  repocreate.Description,
			HelpName:     corecommon.CreateUsage("rt rc", repocreate.Description, repocreate.Usage),
			UsageText:    repocreate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return repoCreateCmd(c)
			},
		},
		{
			Name:         "repo-update",
			Aliases:      []string{"ru"},
			Flags:        cliutils.GetCommandFlags(cliutils.TemplateConsumer),
			Description:  repoupdate.Description,
			HelpName:     corecommon.CreateUsage("rt ru", repoupdate.Description, repoupdate.Usage),
			UsageText:    repoupdate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return repoUpdateCmd(c)
			},
		},
		{
			Name:         "repo-delete",
			Aliases:      []string{"rdel"},
			Flags:        cliutils.GetCommandFlags(cliutils.RepoDelete),
			Description:  repodelete.Description,
			HelpName:     corecommon.CreateUsage("rt rdel", repodelete.Description, repodelete.Usage),
			UsageText:    repodelete.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return repoDeleteCmd(c)
			},
		},
		{
			Name:         "replication-template",
			Aliases:      []string{"rplt"},
			Flags:        cliutils.GetCommandFlags(cliutils.TemplateConsumer),
			Description:  replicationtemplate.Description,
			HelpName:     corecommon.CreateUsage("rt rplt", replicationtemplate.Description, replicationtemplate.Usage),
			UsageText:    replicationtemplate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return replicationTemplateCmd(c)
			},
		},
		{
			Name:         "replication-create",
			Aliases:      []string{"rplc"},
			Flags:        cliutils.GetCommandFlags(cliutils.TemplateConsumer),
			Description:  replicationcreate.Description,
			HelpName:     corecommon.CreateUsage("rt rplc", replicationcreate.Description, replicationcreate.Usage),
			UsageText:    replicationcreate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return replicationCreateCmd(c)
			},
		},
		{
			Name:         "replication-delete",
			Aliases:      []string{"rpldel"},
			Flags:        cliutils.GetCommandFlags(cliutils.ReplicationDelete),
			Description:  replicationdelete.Description,
			HelpName:     corecommon.CreateUsage("rt rpldel", replicationdelete.Description, replicationdelete.Usage),
			UsageText:    replicationdelete.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return replicationDeleteCmd(c)
			},
		},
		{
			Name:         "permission-target-template",
			Aliases:      []string{"ptt"},
			Description:  permissiontargettemplate.Description,
			HelpName:     corecommon.CreateUsage("rt ptt", permissiontargettemplate.Description, permissiontargettemplate.Usage),
			UsageText:    permissiontargettemplate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return permissionTargrtTemplateCmd(c)
			},
		},
		{
			Name:         "permission-target-create",
			Aliases:      []string{"ptc"},
			Flags:        cliutils.GetCommandFlags(cliutils.TemplateConsumer),
			Description:  permissiontargetcreate.Description,
			HelpName:     corecommon.CreateUsage("rt ptc", permissiontargetcreate.Description, permissiontargetcreate.Usage),
			UsageText:    permissiontargetcreate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return permissionTargetCreateCmd(c)
			},
		},
		{
			Name:         "permission-target-update",
			Aliases:      []string{"ptu"},
			Flags:        cliutils.GetCommandFlags(cliutils.TemplateConsumer),
			Description:  permissiontargetupdate.Description,
			HelpName:     corecommon.CreateUsage("rt ptu", permissiontargetupdate.Description, permissiontargetupdate.Usage),
			UsageText:    permissiontargetupdate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return permissionTargetUpdateCmd(c)
			},
		},
		{
			Name:         "permission-target-delete",
			Aliases:      []string{"ptdel"},
			Flags:        cliutils.GetCommandFlags(cliutils.PermissionTargetDelete),
			Description:  permissiontargetdelete.Description,
			HelpName:     corecommon.CreateUsage("rt ptdel", permissiontargetdelete.Description, permissiontargetdelete.Usage),
			UsageText:    permissiontargetdelete.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return permissionTargetDeleteCmd(c)
			},
		},
		{
			Name:         "user-create",
			Flags:        cliutils.GetCommandFlags(cliutils.UserCreate),
			Description:  usercreate.Description,
			HelpName:     corecommon.CreateUsage("rt user-create", usercreate.Description, usercreate.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return userCreateCmd(c)
			},
		},
		{
			Name:         "users-create",
			Aliases:      []string{"uc"},
			Flags:        cliutils.GetCommandFlags(cliutils.UsersCreate),
			Description:  userscreate.Description,
			HelpName:     corecommon.CreateUsage("rt uc", userscreate.Description, userscreate.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return usersCreateCmd(c)
			},
		},
		{
			Name:         "users-delete",
			Aliases:      []string{"udel"},
			Flags:        cliutils.GetCommandFlags(cliutils.UsersDelete),
			Description:  usersdelete.Description,
			HelpName:     corecommon.CreateUsage("rt udel", usersdelete.Description, usersdelete.Usage),
			UsageText:    usersdelete.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return usersDeleteCmd(c)
			},
		},
		{
			Name:         "group-create",
			Aliases:      []string{"gc"},
			Flags:        cliutils.GetCommandFlags(cliutils.GroupCreate),
			Description:  groupcreate.Description,
			HelpName:     corecommon.CreateUsage("rt gc", groupcreate.Description, groupcreate.Usage),
			UsageText:    groupcreate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return groupCreateCmd(c)
			},
		},
		{
			Name:         "group-add-users",
			Aliases:      []string{"gau"},
			Flags:        cliutils.GetCommandFlags(cliutils.GroupAddUsers),
			Description:  groupaddusers.Description,
			HelpName:     corecommon.CreateUsage("rt gau", groupaddusers.Description, groupaddusers.Usage),
			UsageText:    groupaddusers.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return groupAddUsersCmd(c)
			},
		},
		{
			Name:         "group-delete",
			Aliases:      []string{"gdel"},
			Flags:        cliutils.GetCommandFlags(cliutils.GroupDelete),
			Description:  groupdelete.Description,
			HelpName:     corecommon.CreateUsage("rt gdel", groupdelete.Description, groupdelete.Usage),
			UsageText:    groupdelete.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return groupDeleteCmd(c)
			},
		},
		{
			Name:         "access-token-create",
			Aliases:      []string{"atc"},
			Flags:        cliutils.GetCommandFlags(cliutils.AccessTokenCreate),
			Description:  accesstokencreate.Description,
			HelpName:     corecommon.CreateUsage("rt atc", accesstokencreate.Description, accesstokencreate.Usage),
			UsageText:    accesstokencreate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommon.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return accessTokenCreateCmd(c)
			},
		},
	})
}

func createArtifactoryDetailsByFlags(c *cli.Context) (*coreConfig.ServerDetails, error) {
	artDetails, err := cliutils.CreateServerDetailsWithConfigOffer(c, false, cliutils.Rt)
	if err != nil {
		return nil, err
	}
	if artDetails.ArtifactoryUrl == "" {
		return nil, errors.New("the --url option is mandatory")
	}
	return artDetails, nil
}

func getSplitCount(c *cli.Context) (splitCount int, err error) {
	splitCount = cliutils.DownloadSplitCount
	err = nil
	if c.String("split-count") != "" {
		splitCount, err = strconv.Atoi(c.String("split-count"))
		if err != nil {
			err = errors.New("The '--split-count' option should have a numeric value. " + cliutils.GetDocumentationMessage())
		}
		if splitCount > cliutils.DownloadMaxSplitCount {
			err = errors.New("The '--split-count' option value is limited to a maximum of " + strconv.Itoa(cliutils.DownloadMaxSplitCount) + ".")
		}
		if splitCount < 0 {
			err = errors.New("the '--split-count' option cannot have a negative value")
		}
	}
	return
}

func getMinSplit(c *cli.Context) (minSplitSize int64, err error) {
	minSplitSize = cliutils.DownloadMinSplitKb
	err = nil
	if c.String("min-split") != "" {
		minSplitSize, err = strconv.ParseInt(c.String("min-split"), 10, 64)
		if err != nil {
			err = errors.New("The '--min-split' option should have a numeric value. " + cliutils.GetDocumentationMessage())
			return 0, err
		}
	}

	return minSplitSize, nil
}

func getRetries(c *cli.Context) (retries int, err error) {
	retries = cliutils.Retries
	err = nil
	if c.String("retries") != "" {
		retries, err = strconv.Atoi(c.String("retries"))
		if err != nil {
			err = errors.New("The '--retries' option should have a numeric value. " + cliutils.GetDocumentationMessage())
			return 0, err
		}
	}

	return retries, nil
}

func mvnCmd(c *cli.Context) error {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return err
	}

	configFilePath, exists, err := utils.GetProjectConfFilePath(utils.Maven)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("No config file was found! Before running the mvn command on a project for the first time, the project should be configured with the mvn-config command. ")
	}
	if c.NArg() < 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	args := cliutils.ExtractCommand(c)
	filteredMavenArgs, insecureTls, err := coreutils.ExtractInsecureTlsFromArgs(args)
	if err != nil {
		return err
	}
	filteredMavenArgs, buildConfiguration, err := utils.ExtractBuildDetailsFromArgs(filteredMavenArgs)
	if err != nil {
		return err
	}
	filteredMavenArgs, threads, err := extractThreadsFlag(filteredMavenArgs)
	if err != nil {
		return err
	}
	filteredMavenArgs, detailedSummary, err := coreutils.ExtractDetailedSummaryFromArgs(filteredMavenArgs)
	if err != nil {
		return err
	}
	filteredMavenArgs, xrayScan, err := coreutils.ExtractXrayScanFromArgs(filteredMavenArgs)
	if err != nil {
		return err
	}
	mvnCmd := mvn.NewMvnCommand().SetConfiguration(buildConfiguration).SetConfigPath(configFilePath).SetGoals(filteredMavenArgs).SetThreads(threads).SetInsecureTls(insecureTls).SetDetailedSummary(detailedSummary).SetXrayScan(xrayScan)
	err = commands.Exec(mvnCmd)
	if err != nil {
		return err
	}
	if mvnCmd.IsDetailedSummary() {
		return PrintDetailedSummaryReport(c, err, mvnCmd.Result())
	}
	return nil
}

func gradleCmd(c *cli.Context) error {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return err
	}

	configFilePath, exists, err := utils.GetProjectConfFilePath(utils.Gradle)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("No config file was found! Before running the gradle command on a project for the first time, the project should be configured with the gradle-config command. ")
	}
	// Found a config file. Continue as native command.
	if c.NArg() < 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	args := cliutils.ExtractCommand(c)
	filteredGradleArgs, buildConfiguration, err := utils.ExtractBuildDetailsFromArgs(args)
	if err != nil {
		return err
	}
	filteredGradleArgs, threads, err := extractThreadsFlag(filteredGradleArgs)
	if err != nil {
		return err
	}
	filteredGradleArgs, detailedSummary, err := coreutils.ExtractDetailedSummaryFromArgs(filteredGradleArgs)
	if err != nil {
		return err
	}
	filteredGradleArgs, xrayScan, err := coreutils.ExtractXrayScanFromArgs(filteredGradleArgs)
	if err != nil {
		return err
	}
	gradleCmd := gradle.NewGradleCommand().SetConfiguration(buildConfiguration).SetTasks(strings.Join(filteredGradleArgs, " ")).SetConfigPath(configFilePath).SetThreads(threads).SetDetailedSummary(detailedSummary).SetXrayScan(xrayScan)
	err = commands.Exec(gradleCmd)
	if err != nil {
		return err
	}
	if gradleCmd.IsDetailedSummary() {
		return PrintDetailedSummaryReport(c, err, gradleCmd.Result())
	}
	return nil
}

func PrintDetailedSummaryReport(c *cli.Context, originalErr error, result *commandsutils.Result) error {
	if len(result.Reader().GetFilesPaths()) == 0 {
		return errorutils.CheckError(errors.New("Empty reader - no files paths."))
	}
	defer os.Remove(result.Reader().GetFilesPaths()[0])
	err := cliutils.PrintDetailedSummaryReport(result.SuccessCount(), result.FailCount(), result.Reader(), true, originalErr)
	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

func dockerPromoteCmd(c *cli.Context) error {
	if c.NArg() != 3 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	artDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	params := services.NewDockerPromoteParams(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	params.TargetDockerImage = c.String("target-docker-image")
	params.SourceTag = c.String("source-tag")
	params.TargetTag = c.String("target-tag")
	params.Copy = c.Bool("copy")
	dockerPromoteCommand := container.NewDockerPromoteCommand()
	dockerPromoteCommand.SetParams(params).SetServerDetails(artDetails)

	return commands.Exec(dockerPromoteCommand)
}

func containerPushCmd(c *cli.Context, containerManagerType containerutils.ContainerManagerType) error {
	if c.NArg() != 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	artDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	imageTag := c.Args().Get(0)
	targetRepo := c.Args().Get(1)
	skipLogin := c.Bool("skip-login")

	buildConfiguration, err := createBuildConfigurationWithModule(c)
	if err != nil {
		return err
	}
	dockerPushCommand := container.NewPushCommand(containerManagerType)
	threads, err := cliutils.GetThreadsCount(c)
	if err != nil {
		return err
	}
	dockerPushCommand.SetThreads(threads).SetDetailedSummary(c.Bool("detailed-summary")).SetBuildConfiguration(buildConfiguration).SetRepo(targetRepo).SetSkipLogin(skipLogin).SetServerDetails(artDetails).SetImageTag(imageTag)

	err = commands.Exec(dockerPushCommand)
	if err != nil {
		return err
	}
	if dockerPushCommand.IsDetailedSummary() {
		result := dockerPushCommand.Result()
		return cliutils.PrintDetailedSummaryReport(result.SuccessCount(), result.FailCount(), result.Reader(), true, err)
	}
	return nil
}

func containerPullCmd(c *cli.Context, containerManagerType containerutils.ContainerManagerType) error {
	if c.NArg() != 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	artDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	imageTag := c.Args().Get(0)
	sourceRepo := c.Args().Get(1)
	skipLogin := c.Bool("skip-login")
	buildConfiguration, err := createBuildConfigurationWithModule(c)
	if err != nil {
		return err
	}
	dockerPullCommand := container.NewPullCommand(containerManagerType)
	dockerPullCommand.SetImageTag(imageTag).SetRepo(sourceRepo).SetSkipLogin(skipLogin).SetServerDetails(artDetails).SetBuildConfiguration(buildConfiguration)

	return commands.Exec(dockerPullCommand)
}

func BuildDockerCreateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	artDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	sourceRepo := c.Args().Get(0)
	imageNameWithDigestFile := c.String("image-file")
	if imageNameWithDigestFile == "" {
		return cliutils.PrintHelpAndReturnError("The '--image-file' command option was not provided.", c)
	}
	buildConfiguration, err := createBuildConfigurationWithModule(c)
	if err != nil {
		return err
	}
	buildDockerCreateCommand := container.NewBuildDockerCreateCommand()
	if err := buildDockerCreateCommand.SetImageNameWithDigest(imageNameWithDigestFile); err != nil {
		return err
	}
	buildDockerCreateCommand.SetRepo(sourceRepo).SetServerDetails(artDetails).SetBuildConfiguration(buildConfiguration)
	return commands.Exec(buildDockerCreateCommand)
}

func nugetCmd(c *cli.Context) error {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return err
	}
	if c.NArg() < 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	configFilePath, exists, err := utils.GetProjectConfFilePath(utils.Nuget)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(fmt.Sprintf("No config file was found! Before running the nuget command on a project for the first time, the project should be configured using the nuget-config command."))
	}

	rtDetails, targetRepo, useNugetV2, err := getNugetAndDotnetConfigFields(configFilePath)
	if err != nil {
		return err
	}
	args := cliutils.ExtractCommand(c)
	filteredNugetArgs, buildConfiguration, err := utils.ExtractBuildDetailsFromArgs(args)
	if err != nil {
		return err
	}

	nugetCmd := dotnet.NewNugetCommand()
	nugetCmd.SetServerDetails(rtDetails).SetRepoName(targetRepo).SetBuildConfiguration(buildConfiguration).
		SetBasicCommand(filteredNugetArgs[0]).SetUseNugetV2(useNugetV2)
	// Since we are using the values of the command's arguments and flags along the buildInfo collection process,
	// we want to separate the actual NuGet basic command (restore/build...) from the arguments and flags
	if len(filteredNugetArgs) > 1 {
		nugetCmd.SetArgAndFlags(filteredNugetArgs[1:])
	}
	return commands.Exec(nugetCmd)
}

func nugetDepsTreeCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	return dotnet.DependencyTreeCmd()
}

func dotnetCmd(c *cli.Context) error {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return err
	}

	if c.NArg() < 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	// Get configuration file path.
	configFilePath, exists, err := utils.GetProjectConfFilePath(utils.Dotnet)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(fmt.Sprintf("Error occurred while attempting to read dotnet-configuration file.\n" +
			"Please run 'jfrog rt dotnet-config' command prior to running 'jfrog rt dotnet'."))
	}

	rtDetails, targetRepo, useNugetV2, err := getNugetAndDotnetConfigFields(configFilePath)
	if err != nil {
		return err
	}

	args := cliutils.ExtractCommand(c)

	filteredDotnetArgs, buildConfiguration, err := utils.ExtractBuildDetailsFromArgs(args)
	if err != nil {
		return err
	}

	// Run command.
	dotnetCmd := dotnet.NewDotnetCoreCliCommand()
	dotnetCmd.SetServerDetails(rtDetails).SetRepoName(targetRepo).SetBuildConfiguration(buildConfiguration).
		SetBasicCommand(filteredDotnetArgs[0]).SetUseNugetV2(useNugetV2)
	// Since we are using the values of the command's arguments and flags along the buildInfo collection process,
	// we want to separate the actual .NET basic command (restore/build...) from the arguments and flags
	if len(filteredDotnetArgs) > 1 {
		dotnetCmd.SetArgAndFlags(filteredDotnetArgs[1:])
	}
	return commands.Exec(dotnetCmd)
}

func getNugetAndDotnetConfigFields(configFilePath string) (rtDetails *coreConfig.ServerDetails, targetRepo string, useNugetV2 bool, err error) {
	vConfig, err := utils.ReadConfigFile(configFilePath, utils.YAML)
	if err != nil {
		return nil, "", false, errors.New(fmt.Sprintf("Error occurred while attempting to read nuget-configuration file: %s", err.Error()))
	}
	projectConfig, err := utils.GetRepoConfigByPrefix(configFilePath, utils.ProjectConfigResolverPrefix, vConfig)
	if err != nil {
		return nil, "", false, err
	}
	rtDetails, err = projectConfig.ServerDetails()
	if err != nil {
		return nil, "", false, err
	}
	targetRepo = projectConfig.TargetRepo()
	useNugetV2 = vConfig.GetBool(utils.ProjectConfigResolverPrefix + "." + "nugetV2")
	return
}

func npmInstallOrCiCmd(c *cli.Context) error {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return err
	}

	configFilePath, exists, err := utils.GetProjectConfFilePath(utils.Npm)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("No config file was found! Before running the npm-install or npm-ci command on a project for the first time, the project should be configured using the npm-config command. ")
	}
	npmCmd := npm.NewNpmInstallCommand()
	args := cliutils.ExtractCommand(c)
	npmCmd.SetConfigFilePath(configFilePath).SetArgs(args)
	return commands.Exec(npmCmd)
}

func npmPublishCmd(c *cli.Context) error {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return err
	}

	configFilePath, exists, err := utils.GetProjectConfFilePath(utils.Npm)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("No config file was found! Before running the npm-publish command on a project for the first time, the project should be configured using the npm-config command.\nThis configuration includes the Artifactory server and repository to which the package should deployed. ")
	}
	args := cliutils.ExtractCommand(c)
	npmCmd := npm.NewNpmPublishCommand()
	npmCmd.SetConfigFilePath(configFilePath).SetArgs(args)
	err = commands.Exec(npmCmd)
	if err != nil {
		return err
	}
	if npmCmd.IsDetailedSummary() {
		result := npmCmd.Result()
		return cliutils.PrintDetailedSummaryReport(result.SuccessCount(), result.FailCount(), result.Reader(), true, err)
	}
	return nil
}

func yarnCmd(c *cli.Context) error {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return err
	}

	configFilePath, exists, err := utils.GetProjectConfFilePath(utils.Yarn)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(fmt.Sprintf("JFrog CLI's Yarn configuration file was not found.\n" +
			"Run 'jfrog rt yarn-config' command to create it prior to running 'jfrog rt yarn'."))
	}

	yarnCmd := yarn.NewYarnCommand().SetConfigFilePath(configFilePath).SetArgs(c.Args())
	return commands.Exec(yarnCmd)
}

func shouldSkipGoFlagParsing() bool {
	// This function is executed by code-gangsta, regardless of the CLI command being executed.
	// There's no need to run the code of this function, if the command is not "jfrog rt go*".
	if len(os.Args) < 3 || os.Args[2] != "go" {
		return false
	}

	_, exists, err := utils.GetProjectConfFilePath(utils.Go)
	if err != nil {
		coreutils.ExitOnErr(err)
	}
	return exists
}

func shouldSkipNpmFlagParsing() bool {
	// This function is executed by code-gangsta, regardless of the CLI command being executed.
	// There's no need to run the code of this function, if the command is not "jfrog rt npm*".
	if len(os.Args) < 3 || !npmUtils.IsNpmCommand(os.Args[2]) {
		return false
	}

	_, exists, err := utils.GetProjectConfFilePath(utils.Npm)
	if err != nil {
		coreutils.ExitOnErr(err)
	}
	return exists
}

func shouldSkipNugetFlagParsing() bool {
	// This function is executed by code-gangsta, regardless of the CLI command being executed.
	// There's no need to run the code of this function, if the command is not "jfrog rt nuget*".
	if len(os.Args) < 3 || os.Args[2] != "nuget" {
		return false
	}

	_, exists, err := utils.GetProjectConfFilePath(utils.Nuget)
	if err != nil {
		coreutils.ExitOnErr(err)
	}
	return exists
}

func shouldSkipMavenFlagParsing() bool {
	// This function is executed by code-gangsta, regardless of the CLI command being executed.
	// There's no need to run the code of this function, if the command is not "jfrog rt mvn*".
	if len(os.Args) < 3 || os.Args[2] != "mvn" {
		return false
	}
	_, exists, err := utils.GetProjectConfFilePath(utils.Maven)
	if err != nil {
		coreutils.ExitOnErr(err)
	}
	return exists
}

func shouldSkipGradleFlagParsing() bool {
	// This function is executed by code-gangsta, regardless of the CLI command being executed.
	// There's no need to run the code of this function, if the command is not "jfrog rt gradle*".
	if len(os.Args) < 3 || os.Args[2] != "gradle" {
		return false
	}
	_, exists, err := utils.GetProjectConfFilePath(utils.Gradle)
	if err != nil {
		coreutils.ExitOnErr(err)
	}
	return exists
}

func goCmd(c *cli.Context) error {
	configFilePath, err := goCmdVerification(c)
	if err != nil {
		return err
	}
	args := cliutils.ExtractCommand(c)
	goCommand := golang.NewGoCommand()
	goCommand.SetConfigFilePath(configFilePath).SetGoArg(args)
	return commands.Exec(goCommand)
}

func goPublishCmd(c *cli.Context) error {
	configFilePath, err := goCmdVerification(c)
	if err != nil {
		return err
	}
	buildConfiguration, err := createBuildConfigurationWithModule(c)
	if err != nil {
		return err
	}
	version := c.Args().Get(0)
	goPublishCmd := golang.NewGoPublishCommand()
	goPublishCmd.SetConfigFilePath(configFilePath).SetBuildConfiguration(buildConfiguration).SetVersion(version).SetDetailedSummary(c.Bool("detailed-summary"))
	err = commands.Exec(goPublishCmd)
	result := goPublishCmd.Result()
	return cliutils.PrintDetailedSummaryReport(result.SuccessCount(), result.FailCount(), result.Reader(), true, err)
}

func goCmdVerification(c *cli.Context) (string, error) {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return "", err
	}
	if c.NArg() < 1 {
		return "", cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	configFilePath, exists, err := utils.GetProjectConfFilePath(utils.Go)
	if err != nil {
		return "", err
	}
	// Verify config file is found.
	if !exists {
		return "", errors.New(fmt.Sprintf("No config file was found! Before running the go command on a project for the first time, the project should be configured using the go-config command."))
	}
	log.Debug("Go config file was found in:", configFilePath)
	return configFilePath, nil
}

func createGradleConfigCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	return commandUtils.CreateBuildConfig(c, utils.Gradle)
}

func createMvnConfigCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	return commandUtils.CreateBuildConfig(c, utils.Maven)
}

func createGoConfigCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	return commandUtils.CreateBuildConfig(c, utils.Go)
}

func createNpmConfigCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	return commandUtils.CreateBuildConfig(c, utils.Npm)
}

func createYarnConfigCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	return commandUtils.CreateBuildConfig(c, utils.Yarn)
}

func createNugetConfigCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	return commandUtils.CreateBuildConfig(c, utils.Nuget)
}

func createDotnetConfigCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	return commandUtils.CreateBuildConfig(c, utils.Dotnet)
}

func createPipConfigCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	return commandUtils.CreateBuildConfig(c, utils.Pip)
}

func pingCmd(c *cli.Context) error {
	if c.NArg() > 0 {
		return cliutils.PrintHelpAndReturnError("No arguments should be sent.", c)
	}
	artDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	pingCmd := generic.NewPingCommand()
	pingCmd.SetServerDetails(artDetails)
	err = commands.Exec(pingCmd)
	resString := clientutils.IndentJson(pingCmd.Response())
	if err != nil {
		return errors.New(err.Error() + "\n" + resString)
	}
	log.Output(resString)

	return err
}

func prepareDownloadCommand(c *cli.Context) (*spec.SpecFiles, error) {
	if c.NArg() > 0 && c.IsSet("spec") {
		return nil, cliutils.PrintHelpAndReturnError("No arguments should be sent when the spec option is used.", c)
	}
	if !(c.NArg() == 1 || c.NArg() == 2 || (c.NArg() == 0 && (c.IsSet("spec") || c.IsSet("build") || c.IsSet("bundle")))) {
		return nil, cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	var downloadSpec *spec.SpecFiles
	var err error
	if c.IsSet("spec") {
		downloadSpec, err = cliutils.GetSpec(c, true)
	} else {
		downloadSpec, err = createDefaultDownloadSpec(c)
	}
	if err != nil {
		return nil, err
	}
	err = spec.ValidateSpec(downloadSpec.Files, false, true, false)
	if err != nil {
		return nil, err
	}
	return downloadSpec, nil
}

func downloadCmd(c *cli.Context) error {
	downloadSpec, err := prepareDownloadCommand(c)
	if err != nil {
		return err
	}
	fixWinPathsForDownloadCmd(downloadSpec, c)
	configuration, err := createDownloadConfiguration(c)
	if err != nil {
		return err
	}
	serverDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	buildConfiguration, err := createBuildConfigurationWithModule(c)
	if err != nil {
		return err
	}
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	downloadCommand := generic.NewDownloadCommand()
	downloadCommand.SetConfiguration(configuration).SetBuildConfiguration(buildConfiguration).SetSpec(downloadSpec).SetServerDetails(serverDetails).SetDryRun(c.Bool("dry-run")).SetSyncDeletesPath(c.String("sync-deletes")).SetQuiet(cliutils.GetQuietValue(c)).SetDetailedSummary(c.Bool("detailed-summary")).SetRetries(retries)

	if downloadCommand.ShouldPrompt() && !coreutils.AskYesNo("Sync-deletes may delete some files in your local file system. Are you sure you want to continue?\n"+
		"You can avoid this confirmation message by adding --quiet to the command.", false) {
		return nil
	}

	err = execWithProgress(downloadCommand)
	result := downloadCommand.Result()
	err = cliutils.PrintDetailedSummaryReport(result.SuccessCount(), result.FailCount(), result.Reader(), false, err)

	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

func uploadCmd(c *cli.Context) error {
	if c.NArg() > 0 && c.IsSet("spec") {
		return cliutils.PrintHelpAndReturnError("No arguments should be sent when the spec option is used.", c)
	}
	if !(c.NArg() == 2 || (c.NArg() == 0 && c.IsSet("spec"))) {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	var uploadSpec *spec.SpecFiles
	var err error
	if c.IsSet("spec") {
		uploadSpec, err = cliutils.GetFileSystemSpec(c)
	} else {
		uploadSpec, err = createDefaultUploadSpec(c)
	}
	if err != nil {
		return err
	}
	err = spec.ValidateSpec(uploadSpec.Files, true, false, true)
	if err != nil {
		return err
	}
	cliutils.FixWinPathsForFileSystemSourcedCmds(uploadSpec, c)
	configuration, err := createUploadConfiguration(c)
	if err != nil {
		return err
	}
	buildConfiguration, err := createBuildConfigurationWithModule(c)
	if err != nil {
		return err
	}
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	uploadCmd := generic.NewUploadCommand()
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	uploadCmd.SetUploadConfiguration(configuration).SetBuildConfiguration(buildConfiguration).SetSpec(uploadSpec).SetServerDetails(rtDetails).SetDryRun(c.Bool("dry-run")).SetSyncDeletesPath(c.String("sync-deletes")).SetQuiet(cliutils.GetQuietValue(c)).SetDetailedSummary(c.Bool("detailed-summary")).SetRetries(retries)

	if uploadCmd.ShouldPrompt() && !coreutils.AskYesNo("Sync-deletes may delete some artifacts in Artifactory. Are you sure you want to continue?\n"+
		"You can avoid this confirmation message by adding --quiet to the command.", false) {
		return nil
	}
	err = execWithProgress(uploadCmd)
	result := uploadCmd.Result()
	err = cliutils.PrintDetailedSummaryReport(result.SuccessCount(), result.FailCount(), result.Reader(), true, err)

	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

type CommandWithProgress interface {
	commands.Command
	SetProgress(ioUtils.ProgressMgr)
}

func execWithProgress(cmd CommandWithProgress) error {
	// Init progress bar.
	progressBar, logFile, err := progressbar.InitProgressBarIfPossible()
	if err != nil {
		return err
	}
	if progressBar != nil {
		cmd.SetProgress(progressBar)
		defer logUtils.CloseLogFile(logFile)
		defer progressBar.Quit()
	}
	return commands.Exec(cmd)
}

func prepareCopyMoveCommand(c *cli.Context) (*spec.SpecFiles, error) {
	if c.NArg() > 0 && c.IsSet("spec") {
		return nil, cliutils.PrintHelpAndReturnError("No arguments should be sent when the spec option is used.", c)
	}
	if !(c.NArg() == 2 || (c.NArg() == 0 && (c.IsSet("spec")))) {
		return nil, cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	var copyMoveSpec *spec.SpecFiles
	var err error
	if c.IsSet("spec") {
		copyMoveSpec, err = cliutils.GetSpec(c, false)
	} else {
		copyMoveSpec, err = createDefaultCopyMoveSpec(c)
	}
	if err != nil {
		return nil, err
	}
	err = spec.ValidateSpec(copyMoveSpec.Files, true, true, false)
	if err != nil {
		return nil, err
	}
	return copyMoveSpec, nil
}

func moveCmd(c *cli.Context) error {
	moveSpec, err := prepareCopyMoveCommand(c)
	if err != nil {
		return err
	}
	moveCmd := generic.NewMoveCommand()
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	threads, err := cliutils.GetThreadsCount(c)
	if err != nil {
		return err
	}
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	moveCmd.SetThreads(threads).SetDryRun(c.Bool("dry-run")).SetServerDetails(rtDetails).SetSpec(moveSpec).SetRetries(retries)
	err = commands.Exec(moveCmd)
	result := moveCmd.Result()
	err = cliutils.PrintSummaryReport(result.SuccessCount(), result.FailCount(), err)

	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

func copyCmd(c *cli.Context) error {
	copySpec, err := prepareCopyMoveCommand(c)
	if err != nil {
		return err
	}

	copyCommand := generic.NewCopyCommand()
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	threads, err := cliutils.GetThreadsCount(c)
	if err != nil {
		return err
	}
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	copyCommand.SetThreads(threads).SetSpec(copySpec).SetDryRun(c.Bool("dry-run")).SetServerDetails(rtDetails).SetRetries(retries)
	err = commands.Exec(copyCommand)
	result := copyCommand.Result()
	err = cliutils.PrintSummaryReport(result.SuccessCount(), result.FailCount(), err)

	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

func prepareDeleteCommand(c *cli.Context) (*spec.SpecFiles, error) {
	if c.NArg() > 0 && c.IsSet("spec") {
		return nil, cliutils.PrintHelpAndReturnError("No arguments should be sent when the spec option is used.", c)
	}
	if !(c.NArg() == 1 || (c.NArg() == 0 && (c.IsSet("spec") || c.IsSet("build") || c.IsSet("bundle")))) {
		return nil, cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	var deleteSpec *spec.SpecFiles
	var err error
	if c.IsSet("spec") {
		deleteSpec, err = cliutils.GetSpec(c, false)
	} else {
		deleteSpec, err = createDefaultDeleteSpec(c)
	}
	if err != nil {
		return nil, err
	}
	err = spec.ValidateSpec(deleteSpec.Files, false, true, false)
	if err != nil {
		return nil, err
	}
	return deleteSpec, nil
}

func deleteCmd(c *cli.Context) error {
	deleteSpec, err := prepareDeleteCommand(c)
	if err != nil {
		return err
	}

	deleteCommand := generic.NewDeleteCommand()
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	threads, err := cliutils.GetThreadsCount(c)
	if err != nil {
		return err
	}
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	deleteCommand.SetThreads(threads).SetQuiet(cliutils.GetQuietValue(c)).SetDryRun(c.Bool("dry-run")).SetServerDetails(rtDetails).SetSpec(deleteSpec).SetRetries(retries)
	err = commands.Exec(deleteCommand)
	result := deleteCommand.Result()
	err = cliutils.PrintSummaryReport(result.SuccessCount(), result.FailCount(), err)

	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

func prepareSearchCommand(c *cli.Context) (*spec.SpecFiles, error) {
	if c.NArg() > 0 && c.IsSet("spec") {
		return nil, cliutils.PrintHelpAndReturnError("No arguments should be sent when the spec option is used.", c)
	}
	if !(c.NArg() == 1 || (c.NArg() == 0 && (c.IsSet("spec") || c.IsSet("build") || c.IsSet("bundle")))) {
		return nil, cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	var searchSpec *spec.SpecFiles
	var err error
	if c.IsSet("spec") {
		searchSpec, err = cliutils.GetSpec(c, false)
	} else {
		searchSpec, err = createDefaultSearchSpec(c)
	}
	if err != nil {
		return nil, err
	}
	err = spec.ValidateSpec(searchSpec.Files, false, true, false)
	if err != nil {
		return nil, err
	}
	return searchSpec, err
}

func searchCmd(c *cli.Context) error {
	searchSpec, err := prepareSearchCommand(c)
	if err != nil {
		return err
	}
	artDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	searchCmd := generic.NewSearchCommand()
	searchCmd.SetServerDetails(artDetails).SetSpec(searchSpec).SetRetries(retries)
	err = commands.Exec(searchCmd)
	if err != nil {
		return err
	}
	reader := searchCmd.Result().Reader()
	defer reader.Close()
	length, err := reader.Length()
	if err != nil {
		return err
	}
	err = cliutils.GetCliError(err, length, 0, isFailNoOp(c))
	if err != nil {
		return err
	}
	if !c.Bool("count") {
		return utils.PrintSearchResults(reader)
	}
	log.Output(length)
	return nil
}

func preparePropsCmd(c *cli.Context) (*generic.PropsCommand, error) {
	if c.NArg() > 1 && c.IsSet("spec") {
		return nil, cliutils.PrintHelpAndReturnError("Only the 'artifact properties' argument should be sent when the spec option is used.", c)
	}
	if !(c.NArg() == 2 || (c.NArg() == 1 && (c.IsSet("spec") || c.IsSet("build") || c.IsSet("bundle")))) {
		return nil, cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	var propsSpec *spec.SpecFiles
	var err error
	var props string
	if c.IsSet("spec") {
		props = c.Args()[0]
		propsSpec, err = cliutils.GetSpec(c, false)
	} else {
		propsSpec, err = createDefaultPropertiesSpec(c)
		if c.NArg() == 1 {
			props = c.Args()[0]
			propsSpec.Get(0).Pattern = "*"
		} else {
			props = c.Args()[1]
		}
	}
	if err != nil {
		return nil, err
	}
	err = spec.ValidateSpec(propsSpec.Files, false, true, false)
	if err != nil {
		return nil, err
	}

	command := generic.NewPropsCommand()
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return nil, err
	}
	threads, err := cliutils.GetThreadsCount(c)
	if err != nil {
		return nil, err
	}

	cmd := command.SetProps(props)
	cmd.SetThreads(threads).SetSpec(propsSpec).SetDryRun(c.Bool("dry-run")).SetServerDetails(rtDetails)
	return cmd, nil
}

func setPropsCmd(c *cli.Context) error {
	cmd, err := preparePropsCmd(c)
	if err != nil {
		return err
	}
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	propsCmd := generic.NewSetPropsCommand().SetPropsCommand(*cmd)
	propsCmd.SetRetries(retries)
	err = commands.Exec(propsCmd)
	result := propsCmd.Result()
	err = cliutils.PrintSummaryReport(result.SuccessCount(), result.FailCount(), err)

	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

func deletePropsCmd(c *cli.Context) error {
	cmd, err := preparePropsCmd(c)
	if err != nil {
		return err
	}
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	propsCmd := generic.NewDeletePropsCommand().DeletePropsCommand(*cmd)
	propsCmd.SetRetries(retries)
	err = commands.Exec(propsCmd)
	result := propsCmd.Result()
	err = cliutils.PrintSummaryReport(result.SuccessCount(), result.FailCount(), err)

	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

func buildPublishCmd(c *cli.Context) error {
	if c.NArg() > 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	buildConfiguration := createBuildConfiguration(c)
	if err := validateBuildConfiguration(c, buildConfiguration); err != nil {
		return err
	}
	buildInfoConfiguration := createBuildInfoConfiguration(c)
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	buildPublishCmd := buildinfo.NewBuildPublishCommand().SetServerDetails(rtDetails).SetBuildConfiguration(buildConfiguration).SetConfig(buildInfoConfiguration).SetDetailedSummary(c.Bool("detailed-summary"))

	err = commands.Exec(buildPublishCmd)
	if buildPublishCmd.IsDetailedSummary() {
		if summary := buildPublishCmd.GetSummary(); summary != nil {
			return cliutils.PrintBuildInfoSummaryReport(summary.IsSucceeded(), summary.GetSha256(), err)
		}
	}
	return err
}

func buildAppendCmd(c *cli.Context) error {
	if c.NArg() != 4 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	buildConfiguration := createBuildConfiguration(c)
	if err := validateBuildConfiguration(c, buildConfiguration); err != nil {
		return err
	}
	buildNameToAppend, buildNumberToAppend := c.Args().Get(2), c.Args().Get(3)
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	buildAppendCmd := buildinfo.NewBuildAppendCommand().SetServerDetails(rtDetails).SetBuildConfiguration(buildConfiguration).SetBuildNameToAppend(buildNameToAppend).SetBuildNumberToAppend(buildNumberToAppend)
	return commands.Exec(buildAppendCmd)
}

func buildAddDependenciesCmd(c *cli.Context) error {
	if c.NArg() > 2 && c.IsSet("spec") {
		return cliutils.PrintHelpAndReturnError("Only path or spec is allowed, not both.", c)
	}
	if c.IsSet("regexp") && c.IsSet("from-rt") {
		return cliutils.PrintHelpAndReturnError("The --regexp option is not supported when --from-rt is set to true.", c)
	}
	buildConfiguration := createBuildConfiguration(c)
	if err := validateBuildConfiguration(c, buildConfiguration); err != nil {
		return err
	}
	// Odd number of args - Use pattern arg
	// Even number of args - Use spec flag
	if c.NArg() > 3 || !(c.NArg()%2 == 1 || (c.NArg()%2 == 0 && c.IsSet("spec"))) {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	var dependenciesSpec *spec.SpecFiles
	var rtDetails *coreConfig.ServerDetails
	var err error
	if c.IsSet("spec") {
		dependenciesSpec, err = cliutils.GetFileSystemSpec(c)
		if err != nil {
			return err
		}
	} else {
		dependenciesSpec = createDefaultBuildAddDependenciesSpec(c)
	}
	if c.Bool("from-rt") {
		rtDetails, err = createArtifactoryDetailsByFlags(c)
		if err != nil {
			return err
		}
	} else {
		cliutils.FixWinPathsForFileSystemSourcedCmds(dependenciesSpec, c)
	}
	buildAddDependenciesCmd := buildinfo.NewBuildAddDependenciesCommand().SetDryRun(c.Bool("dry-run")).SetBuildConfiguration(buildConfiguration).SetDependenciesSpec(dependenciesSpec).SetServerDetails(rtDetails)
	err = commands.Exec(buildAddDependenciesCmd)
	result := buildAddDependenciesCmd.Result()
	err = cliutils.PrintSummaryReport(result.SuccessCount(), result.FailCount(), err)
	if err != nil {
		return err
	}

	return cliutils.GetCliError(err, result.SuccessCount(), result.FailCount(), isFailNoOp(c))
}

func buildCollectEnvCmd(c *cli.Context) error {
	if c.NArg() > 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	buildConfiguration := createBuildConfiguration(c)
	if err := validateBuildConfiguration(c, buildConfiguration); err != nil {
		return err
	}
	buildCollectEnvCmd := buildinfo.NewBuildCollectEnvCommand().SetBuildConfiguration(buildConfiguration)

	return commands.Exec(buildCollectEnvCmd)
}

func buildAddGitCmd(c *cli.Context) error {
	if c.NArg() > 3 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	buildConfiguration := createBuildConfiguration(c)
	if err := validateBuildConfiguration(c, buildConfiguration); err != nil {
		return err
	}

	buildAddGitConfigurationCmd := buildinfo.NewBuildAddGitCommand().SetBuildConfiguration(buildConfiguration).SetConfigFilePath(c.String("config")).SetServerId(c.String("server-id"))
	if c.NArg() == 3 {
		buildAddGitConfigurationCmd.SetDotGitPath(c.Args().Get(2))
	} else if c.NArg() == 1 {
		buildAddGitConfigurationCmd.SetDotGitPath(c.Args().Get(0))
	}
	return commands.Exec(buildAddGitConfigurationCmd)
}

func buildScanCmd(c *cli.Context) error {
	if c.NArg() > 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	buildConfiguration := createBuildConfiguration(c)
	if err := validateBuildConfiguration(c, buildConfiguration); err != nil {
		return err
	}
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	buildScanCmd := buildinfo.NewBuildScanCommand().SetServerDetails(rtDetails).SetFailBuild(c.BoolT("fail")).SetBuildConfiguration(buildConfiguration)
	err = commands.Exec(buildScanCmd)

	return checkBuildScanError(err)
}

func checkBuildScanError(err error) error {
	// If the build was found vulnerable, exit with ExitCodeVulnerableBuild.
	if err == utils.GetBuildScanError() {
		return coreutils.CliError{ExitCode: coreutils.ExitCodeVulnerableBuild, ErrorMsg: err.Error()}
	}
	// If the scan operation failed, for example due to HTTP timeout, exit with ExitCodeError.
	if err != nil {
		return coreutils.CliError{ExitCode: coreutils.ExitCodeError, ErrorMsg: err.Error()}
	}
	return nil
}

func buildCleanCmd(c *cli.Context) error {
	if c.NArg() > 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	buildConfiguration := createBuildConfiguration(c)
	if err := validateBuildConfiguration(c, buildConfiguration); err != nil {
		return err
	}
	buildCleanCmd := buildinfo.NewBuildCleanCommand().SetBuildConfiguration(buildConfiguration)

	return commands.Exec(buildCleanCmd)
}

func buildPromoteCmd(c *cli.Context) error {
	if c.NArg() > 3 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	if err := validateBuildConfiguration(c, createBuildConfiguration(c)); err != nil {
		return err
	}
	configuration := createBuildPromoteConfiguration(c)
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	buildPromotionCmd := buildinfo.NewBuildPromotionCommand().SetDryRun(c.Bool("dry-run")).SetServerDetails(rtDetails).SetPromotionParams(configuration)

	return commands.Exec(buildPromotionCmd)
}

func buildDistributeCmd(c *cli.Context) error {
	if c.NArg() > 3 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	if err := validateBuildConfiguration(c, createBuildConfiguration(c)); err != nil {
		return err
	}
	configuration := createBuildDistributionConfiguration(c)
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	buildDistributeCmd := buildinfo.NewBuildDistributeCommnad().SetDryRun(c.Bool("dry-run")).SetServerDetails(rtDetails).SetBuildDistributionParams(configuration)

	return commands.Exec(buildDistributeCmd)
}

func buildDiscardCmd(c *cli.Context) error {
	if c.NArg() > 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	configuration := createBuildDiscardConfiguration(c)
	if configuration.BuildName == "" {
		return cliutils.PrintHelpAndReturnError("Build name is expected as a command argument or environment variable.", c)
	}
	buildDiscardCmd := buildinfo.NewBuildDiscardCommand()
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	buildDiscardCmd.SetServerDetails(rtDetails).SetDiscardBuildsParams(configuration)

	return commands.Exec(buildDiscardCmd)
}

func gitLfsCleanCmd(c *cli.Context) error {
	if c.NArg() > 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	configuration := createGitLfsCleanConfiguration(c)
	retries, err := getRetries(c)
	if err != nil {
		return err
	}
	gitLfsCmd := generic.NewGitLfsCommand()
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	gitLfsCmd.SetConfiguration(configuration).SetServerDetails(rtDetails).SetDryRun(c.Bool("dry-run")).SetRetries(retries)

	return commands.Exec(gitLfsCmd)
}

func curlCmd(c *cli.Context) error {
	if c.NArg() < 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	rtCurlCommand, err := newRtCurlCommand(c)
	if err != nil {
		return err
	}
	return commands.Exec(rtCurlCommand)
}

func newRtCurlCommand(c *cli.Context) (*curl.RtCurlCommand, error) {
	curlCommand := coreCommonCommands.NewCurlCommand().SetArguments(cliutils.ExtractCommand(c))
	rtCurlCommand := curl.NewRtCurlCommand(*curlCommand)
	rtDetails, err := rtCurlCommand.GetServerDetails()
	if err != nil {
		return nil, err
	}
	rtCurlCommand.SetServerDetails(rtDetails)
	rtCurlCommand.SetUrl(rtDetails.ArtifactoryUrl)
	return rtCurlCommand, err
}

func pipInstallCmd(c *cli.Context) error {
	if show, err := cliutils.ShowCmdHelpIfNeeded(c); show || err != nil {
		return err
	}

	if c.NArg() < 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	// Get pip configuration.
	pipConfig, err := utils.GetResolutionOnlyConfiguration(utils.Pip)
	if err != nil {
		return errors.New(fmt.Sprintf("Error occurred while attempting to read pip-configuration file: %s\n"+
			"Please run 'jfrog rt pip-config' command prior to running 'jfrog rt %s'.", err.Error(), "pip-install"))
	}

	// Set arg values.
	rtDetails, err := pipConfig.ServerDetails()
	if err != nil {
		return err
	}

	// Run command.
	pipCmd := pip.NewPipInstallCommand()
	pipCmd.SetServerDetails(rtDetails).SetRepo(pipConfig.TargetRepo()).SetArgs(cliutils.ExtractCommand(c))
	return commands.Exec(pipCmd)
}

func repoTemplateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	// Run command.
	repoTemplateCmd := repository.NewRepoTemplateCommand()
	repoTemplateCmd.SetTemplatePath(c.Args().Get(0))
	return commands.Exec(repoTemplateCmd)
}

func repoCreateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	// Run command.
	repoCreateCmd := repository.NewRepoCreateCommand()
	repoCreateCmd.SetTemplatePath(c.Args().Get(0)).SetServerDetails(rtDetails).SetVars(c.String("vars"))
	return commands.Exec(repoCreateCmd)
}

func repoUpdateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	// Run command.
	repoUpdateCmd := repository.NewRepoUpdateCommand()
	repoUpdateCmd.SetTemplatePath(c.Args().Get(0)).SetServerDetails(rtDetails).SetVars(c.String("vars"))
	return commands.Exec(repoUpdateCmd)
}

func repoDeleteCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	repoDeleteCmd := repository.NewRepoDeleteCommand()
	repoDeleteCmd.SetRepoPattern(c.Args().Get(0)).SetServerDetails(rtDetails).SetQuiet(cliutils.GetQuietValue(c))
	return commands.Exec(repoDeleteCmd)
}

func replicationTemplateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	replicationTemplateCmd := replication.NewReplicationTemplateCommand()
	replicationTemplateCmd.SetTemplatePath(c.Args().Get(0))
	return commands.Exec(replicationTemplateCmd)
}

func replicationCreateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	replicationCreateCmd := replication.NewReplicationCreateCommand()
	replicationCreateCmd.SetTemplatePath(c.Args().Get(0)).SetServerDetails(rtDetails).SetVars(c.String("vars"))
	return commands.Exec(replicationCreateCmd)
}

func replicationDeleteCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	replicationDeleteCmd := replication.NewReplicationDeleteCommand()
	replicationDeleteCmd.SetRepoKey(c.Args().Get(0)).SetServerDetails(rtDetails).SetQuiet(cliutils.GetQuietValue(c))
	return commands.Exec(replicationDeleteCmd)
}

func permissionTargrtTemplateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	// Run command.
	permissionTargetTemplateCmd := permissiontarget.NewPermissionTargetTemplateCommand()
	permissionTargetTemplateCmd.SetTemplatePath(c.Args().Get(0))
	return commands.Exec(permissionTargetTemplateCmd)
}

func permissionTargetCreateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	// Run command.
	permissionTargetCreateCmd := permissiontarget.NewPermissionTargetCreateCommand()
	permissionTargetCreateCmd.SetTemplatePath(c.Args().Get(0)).SetServerDetails(rtDetails).SetVars(c.String("vars"))
	return commands.Exec(permissionTargetCreateCmd)
}

func permissionTargetUpdateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	// Run command.
	permissionTargetUpdateCmd := permissiontarget.NewPermissionTargetUpdateCommand()
	permissionTargetUpdateCmd.SetTemplatePath(c.Args().Get(0)).SetServerDetails(rtDetails).SetVars(c.String("vars"))
	return commands.Exec(permissionTargetUpdateCmd)
}

func permissionTargetDeleteCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	permissionTargetDeleteCmd := permissiontarget.NewPermissionTargetDeleteCommand()
	permissionTargetDeleteCmd.SetPermissionTargetName(c.Args().Get(0)).SetServerDetails(rtDetails).SetQuiet(cliutils.GetQuietValue(c))
	return commands.Exec(permissionTargetDeleteCmd)
}

func userCreateCmd(c *cli.Context) error {
	if c.NArg() != 3 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	usersCreateCmd := usersmanagement.NewUsersCreateCommand()
	userDetails := services.User{}
	userDetails.Name = c.Args().Get(0)
	userDetails.Password = c.Args().Get(1)
	userDetails.Email = c.Args().Get(2)

	user := []services.User{userDetails}
	var usersGroups []string
	if c.String(cliutils.UsersGroups) != "" {
		usersGroups = strings.Split(c.String(cliutils.UsersGroups), ",")
	}
	if c.String(cliutils.Admin) != "" {
		userDetails.Admin = c.Bool(cliutils.Admin)
	}
	// Run command.
	usersCreateCmd.SetServerDetails(rtDetails).SetUsers(user).SetUsersGroups(usersGroups).SetReplaceIfExists(c.Bool(cliutils.Replace))
	return commands.Exec(usersCreateCmd)
}

func usersCreateCmd(c *cli.Context) error {
	if c.NArg() != 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	usersCreateCmd := usersmanagement.NewUsersCreateCommand()
	csvFilePath := c.String("csv")
	if csvFilePath == "" {
		return cliutils.PrintHelpAndReturnError("missing --csv <File Path>", c)
	}
	usersList, err := parseCSVToUsersList(csvFilePath)
	if err != nil {
		return err
	}
	if len(usersList) < 1 {
		return errorutils.CheckError(errors.New("an empty input file was provided"))
	}
	var usersGroups []string
	if c.String(cliutils.UsersGroups) != "" {
		usersGroups = strings.Split(c.String(cliutils.UsersGroups), ",")
	}
	// Run command.
	usersCreateCmd.SetServerDetails(rtDetails).SetUsers(usersList).SetUsersGroups(usersGroups).SetReplaceIfExists(c.Bool(cliutils.Replace))
	return commands.Exec(usersCreateCmd)
}

func usersDeleteCmd(c *cli.Context) error {
	if c.NArg() > 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	usersDeleteCmd := usersmanagement.NewUsersDeleteCommand()
	var usersNamesList = make([]string, 0)
	csvFilePath := c.String("csv")
	if csvFilePath != "" {
		usersList, err := parseCSVToUsersList(csvFilePath)
		if err != nil {
			return err
		}
		// If --csv <users details file path> provided, parse and append its content to the usersNamesList to be deleted.
		usersNamesList = append(usersNamesList, usersToUsersNamesList(usersList)...)
	}
	// If <users list> provided as arg, append its content to the usersNamesList to be deleted.
	if c.NArg() > 0 {
		usersNamesList = append(usersNamesList, strings.Split(c.Args().Get(0), ",")...)
	}

	if len(usersNamesList) < 1 {
		return cliutils.PrintHelpAndReturnError("missing <users list> OR --csv <users details file path>", c)
	}

	if !cliutils.GetQuietValue(c) && !coreutils.AskYesNo("This command will delete users. Are you sure you want to continue?\n"+
		"You can avoid this confirmation message by adding --quiet to the command.", false) {
		return nil
	}

	// Run command.
	usersDeleteCmd.SetServerDetails(rtDetails).SetUsers(usersNamesList)
	return commands.Exec(usersDeleteCmd)
}

func parseCSVToUsersList(csvFilePath string) ([]services.User, error) {
	var usersList []services.User
	csvInput, err := ioutil.ReadFile(csvFilePath)
	if err != nil {
		return usersList, errorutils.CheckError(err)
	}
	if err = csvutil.Unmarshal(csvInput, &usersList); err != nil {
		return usersList, errorutils.CheckError(err)
	}
	return usersList, nil
}

func usersToUsersNamesList(usersList []services.User) (usersNames []string) {
	for _, user := range usersList {
		if user.Name != "" {
			usersNames = append(usersNames, user.Name)
		}
	}
	return
}

func groupCreateCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	// Run command.
	groupCreateCmd := usersmanagement.NewGroupCreateCommand()
	groupCreateCmd.SetName(c.Args().Get(0)).SetServerDetails(rtDetails).SetReplaceIfExists(c.Bool(cliutils.Replace))
	return commands.Exec(groupCreateCmd)
}

func groupAddUsersCmd(c *cli.Context) error {
	if c.NArg() != 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	// Run command.
	groupAddUsersCmd := usersmanagement.NewGroupUpdateCommand()
	groupAddUsersCmd.SetName(c.Args().Get(0)).SetUsers(strings.Split(c.Args().Get(1), ",")).SetServerDetails(rtDetails)
	return commands.Exec(groupAddUsersCmd)
}

func groupDeleteCmd(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	rtDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}

	if !cliutils.GetQuietValue(c) && !coreutils.AskYesNo("This command will delete the group. Are you sure you want to continue?\n"+
		"You can avoid this confirmation message by adding --quiet to the command.", false) {
		return nil
	}

	// Run command.
	groupDeleteCmd := usersmanagement.NewGroupDeleteCommand()
	groupDeleteCmd.SetName(c.Args().Get(0)).SetServerDetails(rtDetails)
	return commands.Exec(groupDeleteCmd)
}

func accessTokenCreateCmd(c *cli.Context) error {
	if c.NArg() > 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	serverDetails, err := createArtifactoryDetailsByFlags(c)
	if err != nil {
		return err
	}
	// If the username is provided as an argument, then it is used when creating the token.
	// If not, then the configured username (or the the value of the --user option) is used.
	var userName string
	if c.NArg() > 0 {
		userName = c.Args().Get(0)
	} else {
		userName = serverDetails.GetUser()
	}
	expiry, err := cliutils.GetIntFlagValue(c, "expiry", cliutils.TokenExpiry)
	if err != nil {
		return err
	}
	accessTokenCreateCmd := generic.NewAccessTokenCreateCommand()
	accessTokenCreateCmd.SetUserName(userName).SetServerDetails(serverDetails).SetRefreshable(c.Bool("refreshable")).SetExpiry(expiry).SetGroups(c.String("groups")).SetAudience(c.String("audience")).SetGrantAdmin(c.Bool("grant-admin"))
	err = commands.Exec(accessTokenCreateCmd)
	if err != nil {
		return err
	}
	resString, err := accessTokenCreateCmd.Response()
	if err != nil {
		return err
	}
	log.Output(clientutils.IndentJson(resString))

	return nil
}

func validateBuildConfiguration(c *cli.Context, buildConfiguration *utils.BuildConfiguration) error {
	if buildConfiguration.BuildName == "" || buildConfiguration.BuildNumber == "" {
		return cliutils.PrintHelpAndReturnError("Build name and build number are expected as command arguments or environment variables.", c)
	}
	return nil
}

func getDebFlag(c *cli.Context) (deb string, err error) {
	deb = c.String("deb")
	slashesCount := strings.Count(deb, "/") - strings.Count(deb, "\\/")
	if deb != "" && slashesCount != 2 {
		return "", errors.New("the --deb option should be in the form of distribution/component/architecture")
	}
	return deb, nil
}

func createDefaultCopyMoveSpec(c *cli.Context) (*spec.SpecFiles, error) {
	offset, limit, err := getOffsetAndLimitValues(c)
	if err != nil {
		return nil, err
	}
	return spec.NewBuilder().
		Pattern(c.Args().Get(0)).
		Props(c.String("props")).
		ExcludeProps(c.String("exclude-props")).
		Build(c.String("build")).
		ExcludeArtifacts(c.Bool("exclude-artifacts")).
		IncludeDeps(c.Bool("include-deps")).
		Bundle(c.String("bundle")).
		Offset(offset).
		Limit(limit).
		SortOrder(c.String("sort-order")).
		SortBy(cliutils.GetStringsArrFlagValue(c, "sort-by")).
		Recursive(c.BoolT("recursive")).
		Exclusions(cliutils.GetStringsArrFlagValue(c, "exclusions")).
		Flat(c.Bool("flat")).
		IncludeDirs(true).
		Target(c.Args().Get(1)).
		ArchiveEntries(c.String("archive-entries")).
		BuildSpec(), nil
}

func createDefaultDeleteSpec(c *cli.Context) (*spec.SpecFiles, error) {
	offset, limit, err := getOffsetAndLimitValues(c)
	if err != nil {
		return nil, err
	}
	return spec.NewBuilder().
		Pattern(c.Args().Get(0)).
		Props(c.String("props")).
		ExcludeProps(c.String("exclude-props")).
		Build(c.String("build")).
		ExcludeArtifacts(c.Bool("exclude-artifacts")).
		IncludeDeps(c.Bool("include-deps")).
		Bundle(c.String("bundle")).
		Offset(offset).
		Limit(limit).
		SortOrder(c.String("sort-order")).
		SortBy(cliutils.GetStringsArrFlagValue(c, "sort-by")).
		Recursive(c.BoolT("recursive")).
		Exclusions(cliutils.GetStringsArrFlagValue(c, "exclusions")).
		ArchiveEntries(c.String("archive-entries")).
		BuildSpec(), nil
}

func createDefaultSearchSpec(c *cli.Context) (*spec.SpecFiles, error) {
	offset, limit, err := getOffsetAndLimitValues(c)
	if err != nil {
		return nil, err
	}
	return spec.NewBuilder().
		Pattern(c.Args().Get(0)).
		Props(c.String("props")).
		ExcludeProps(c.String("exclude-props")).
		Build(c.String("build")).
		ExcludeArtifacts(c.Bool("exclude-artifacts")).
		IncludeDeps(c.Bool("include-deps")).
		Bundle(c.String("bundle")).
		Offset(offset).
		Limit(limit).
		SortOrder(c.String("sort-order")).
		SortBy(cliutils.GetStringsArrFlagValue(c, "sort-by")).
		Recursive(c.BoolT("recursive")).
		Exclusions(cliutils.GetStringsArrFlagValue(c, "exclusions")).
		IncludeDirs(c.Bool("include-dirs")).
		ArchiveEntries(c.String("archive-entries")).
		Transitive(c.Bool("transitive")).
		BuildSpec(), nil
}

func createDefaultPropertiesSpec(c *cli.Context) (*spec.SpecFiles, error) {
	offset, limit, err := getOffsetAndLimitValues(c)
	if err != nil {
		return nil, err
	}
	return spec.NewBuilder().
		Pattern(c.Args().Get(0)).
		Props(c.String("props")).
		ExcludeProps(c.String("exclude-props")).
		Build(c.String("build")).
		ExcludeArtifacts(c.Bool("exclude-artifacts")).
		IncludeDeps(c.Bool("include-deps")).
		Bundle(c.String("bundle")).
		Offset(offset).
		Limit(limit).
		SortOrder(c.String("sort-order")).
		SortBy(cliutils.GetStringsArrFlagValue(c, "sort-by")).
		Recursive(c.BoolT("recursive")).
		Exclusions(cliutils.GetStringsArrFlagValue(c, "exclusions")).
		IncludeDirs(c.Bool("include-dirs")).
		ArchiveEntries(c.String("archive-entries")).
		BuildSpec(), nil
}

func createBuildInfoConfiguration(c *cli.Context) *buildinfocmd.Configuration {
	flags := new(buildinfocmd.Configuration)
	flags.BuildUrl = cliutils.GetBuildUrl(c.String("build-url"))
	flags.DryRun = c.Bool("dry-run")
	flags.EnvInclude = c.String("env-include")
	flags.EnvExclude = cliutils.GetEnvExclude(c.String("env-exclude"))
	if flags.EnvInclude == "" {
		flags.EnvInclude = "*"
	}
	// Allow to use `env-exclude=""` and get no filters
	if flags.EnvExclude == "" {
		flags.EnvExclude = "*password*;*psw*;*secret*;*key*;*token*"
	}
	return flags
}

func createBuildPromoteConfiguration(c *cli.Context) services.PromotionParams {
	promotionParamsImpl := services.NewPromotionParams()
	promotionParamsImpl.Comment = c.String("comment")
	promotionParamsImpl.SourceRepo = c.String("source-repo")
	promotionParamsImpl.Status = c.String("status")
	promotionParamsImpl.IncludeDependencies = c.Bool("include-dependencies")
	promotionParamsImpl.Copy = c.Bool("copy")
	promotionParamsImpl.Properties = c.String("props")
	promotionParamsImpl.ProjectKey = utils.GetBuildProject(c.String("project"))
	promotionParamsImpl.FailFast = c.BoolT("fail-fast")

	// If the command received 3 args, read the build name, build number
	// and target repo as ags.
	buildName, buildNumber, targetRepo := c.Args().Get(0), c.Args().Get(1), c.Args().Get(2)
	// But if the command received only one arg, the build name and build number
	// are expected as env vars, and only the target repo is received as an arg.
	if len(c.Args()) == 1 {
		buildName, buildNumber, targetRepo = "", "", c.Args().Get(0)
	}

	promotionParamsImpl.BuildName, promotionParamsImpl.BuildNumber = utils.GetBuildNameAndNumber(buildName, buildNumber)
	promotionParamsImpl.TargetRepo = targetRepo
	return promotionParamsImpl
}

func createBuildDiscardConfiguration(c *cli.Context) services.DiscardBuildsParams {
	discardParamsImpl := services.NewDiscardBuildsParams()
	discardParamsImpl.DeleteArtifacts = c.Bool("delete-artifacts")
	discardParamsImpl.MaxBuilds = c.String("max-builds")
	discardParamsImpl.MaxDays = c.String("max-days")
	discardParamsImpl.ExcludeBuilds = c.String("exclude-builds")
	discardParamsImpl.Async = c.Bool("async")
	discardParamsImpl.BuildName = cliutils.GetBuildName(c.Args().Get(0))
	return discardParamsImpl
}

func createBuildDistributionConfiguration(c *cli.Context) services.BuildDistributionParams {
	distributeParamsImpl := services.NewBuildDistributionParams()
	distributeParamsImpl.Publish = c.BoolT("publish")
	distributeParamsImpl.OverrideExistingFiles = c.Bool("override")
	distributeParamsImpl.GpgPassphrase = c.String("passphrase")
	distributeParamsImpl.Async = c.Bool("async")
	distributeParamsImpl.SourceRepos = c.String("source-repos")
	distributeParamsImpl.BuildName, distributeParamsImpl.BuildNumber = utils.GetBuildNameAndNumber(c.Args().Get(0), c.Args().Get(1))
	distributeParamsImpl.TargetRepo = c.Args().Get(2)
	return distributeParamsImpl
}

func createGitLfsCleanConfiguration(c *cli.Context) (gitLfsCleanConfiguration *generic.GitLfsCleanConfiguration) {
	gitLfsCleanConfiguration = new(generic.GitLfsCleanConfiguration)

	gitLfsCleanConfiguration.Refs = c.String("refs")
	if len(gitLfsCleanConfiguration.Refs) == 0 {
		gitLfsCleanConfiguration.Refs = "refs/remotes/*"
	}

	gitLfsCleanConfiguration.Repo = c.String("repo")
	gitLfsCleanConfiguration.Quiet = cliutils.GetQuietValue(c)
	dotGitPath := ""
	if c.NArg() == 1 {
		dotGitPath = c.Args().Get(0)
	}
	gitLfsCleanConfiguration.GitPath = dotGitPath
	return
}

func createDefaultDownloadSpec(c *cli.Context) (*spec.SpecFiles, error) {
	offset, limit, err := getOffsetAndLimitValues(c)
	if err != nil {
		return nil, err
	}
	return spec.NewBuilder().
		Pattern(strings.TrimPrefix(c.Args().Get(0), "/")).
		Props(c.String("props")).
		ExcludeProps(c.String("exclude-props")).
		Build(c.String("build")).
		ExcludeArtifacts(c.Bool("exclude-artifacts")).
		IncludeDeps(c.Bool("include-deps")).
		Bundle(c.String("bundle")).
		Offset(offset).
		Limit(limit).
		SortOrder(c.String("sort-order")).
		SortBy(cliutils.GetStringsArrFlagValue(c, "sort-by")).
		Recursive(c.BoolT("recursive")).
		Exclusions(cliutils.GetStringsArrFlagValue(c, "exclusions")).
		Flat(c.Bool("flat")).
		Explode(c.String("explode")).
		IncludeDirs(c.Bool("include-dirs")).
		Target(c.Args().Get(1)).
		ArchiveEntries(c.String("archive-entries")).
		ValidateSymlinks(c.Bool("validate-symlinks")).
		BuildSpec(), nil
}

func createDownloadConfiguration(c *cli.Context) (downloadConfiguration *utils.DownloadConfiguration, err error) {
	downloadConfiguration = new(utils.DownloadConfiguration)
	downloadConfiguration.MinSplitSize, err = getMinSplit(c)
	if err != nil {
		return nil, err
	}
	downloadConfiguration.SplitCount, err = getSplitCount(c)
	if err != nil {
		return nil, err
	}
	downloadConfiguration.Threads, err = cliutils.GetThreadsCount(c)
	if err != nil {
		return nil, err
	}
	downloadConfiguration.Symlink = true
	return
}

func createDefaultUploadSpec(c *cli.Context) (*spec.SpecFiles, error) {
	offset, limit, err := getOffsetAndLimitValues(c)
	if err != nil {
		return nil, err
	}
	return spec.NewBuilder().
		Pattern(c.Args().Get(0)).
		Props(c.String("props")).
		TargetProps(c.String("target-props")).
		Build(c.String("build")).
		Offset(offset).
		Limit(limit).
		SortOrder(c.String("sort-order")).
		SortBy(cliutils.GetStringsArrFlagValue(c, "sort-by")).
		Recursive(c.BoolT("recursive")).
		Exclusions(cliutils.GetStringsArrFlagValue(c, "exclusions")).
		Flat(c.Bool("flat")).
		Explode(c.String("explode")).
		Regexp(c.Bool("regexp")).
		Ant(c.Bool("ant")).
		IncludeDirs(c.Bool("include-dirs")).
		Target(strings.TrimPrefix(c.Args().Get(1), "/")).
		Symlinks(c.Bool("symlinks")).
		Archive(c.String("archive")).
		BuildSpec(), nil
}

func createDefaultBuildAddDependenciesSpec(c *cli.Context) *spec.SpecFiles {
	pattern := c.Args().Get(2)
	if pattern == "" {
		// Build name and build number from env
		pattern = c.Args().Get(0)
	}
	return spec.NewBuilder().
		Pattern(pattern).
		Recursive(c.BoolT("recursive")).
		Exclusions(cliutils.GetStringsArrFlagValue(c, "exclusions")).
		Regexp(c.Bool("regexp")).
		Ant(c.Bool("ant")).
		BuildSpec()
}

func getFileSystemSpec(c *cli.Context) (fsSpec *spec.SpecFiles, err error) {
	fsSpec, err = spec.CreateSpecFromFile(c.String("spec"), coreutils.SpecVarsStringToMap(c.String("spec-vars")))
	if err != nil {
		return
	}
	// Override spec with CLI options
	for i := 0; i < len(fsSpec.Files); i++ {
		fsSpec.Get(i).Target = strings.TrimPrefix(fsSpec.Get(i).Target, "/")
		cliutils.OverrideFieldsIfSet(fsSpec.Get(i), c)
	}
	return
}

func fixWinPathsForFileSystemSourcedCmds(uploadSpec *spec.SpecFiles, c *cli.Context) {
	if coreutils.IsWindows() {
		for i, file := range uploadSpec.Files {
			uploadSpec.Files[i].Pattern = fixWinPathBySource(file.Pattern, c.IsSet("spec"))
			for j, exclusion := range uploadSpec.Files[i].Exclusions {
				// If exclusions are set, they override the spec value
				uploadSpec.Files[i].Exclusions[j] = fixWinPathBySource(exclusion, c.IsSet("spec") && !c.IsSet("exclusions"))
			}
		}
	}
}

func fixWinPathsForDownloadCmd(uploadSpec *spec.SpecFiles, c *cli.Context) {
	if coreutils.IsWindows() {
		for i, file := range uploadSpec.Files {
			uploadSpec.Files[i].Target = fixWinPathBySource(file.Target, c.IsSet("spec"))
		}
	}
}

func fixWinPathBySource(path string, fromSpec bool) string {
	if strings.Count(path, "/") > 0 {
		// Assuming forward slashes - not doubling backslash to allow regexp escaping
		return ioutils.UnixToWinPathSeparator(path)
	}
	if fromSpec {
		// Doubling backslash only for paths from spec files (that aren't forward slashed)
		return ioutils.DoubleWinPathSeparator(path)
	}
	return path
}

func createUploadConfiguration(c *cli.Context) (uploadConfiguration *utils.UploadConfiguration, err error) {
	uploadConfiguration = new(utils.UploadConfiguration)
	uploadConfiguration.Threads, err = cliutils.GetThreadsCount(c)
	if err != nil {
		return nil, err
	}
	uploadConfiguration.Deb, err = getDebFlag(c)
	if err != nil {
		return
	}
	return
}

func createBuildConfigurationWithModule(c *cli.Context) (buildConfigConfiguration *utils.BuildConfiguration, err error) {
	buildConfigConfiguration = new(utils.BuildConfiguration)
	buildConfigConfiguration.BuildName, buildConfigConfiguration.BuildNumber = utils.GetBuildNameAndNumber(c.String("build-name"), c.String("build-number"))
	buildConfigConfiguration.Project = utils.GetBuildProject(c.String("project"))
	buildConfigConfiguration.Module = c.String("module")
	err = utils.ValidateBuildAndModuleParams(buildConfigConfiguration)
	return
}

func validateConfigFlags(configCommandConfiguration *coreCommonCommands.ConfigCommandConfiguration) error {
	if !configCommandConfiguration.Interactive && configCommandConfiguration.ServerDetails.ArtifactoryUrl == "" {
		return errors.New("the --url option is mandatory when the --interactive option is set to false or the CI environment variable is set to true.")
	}
	// Validate the option is not used along with an access token
	if configCommandConfiguration.BasicAuthOnly && configCommandConfiguration.ServerDetails.AccessToken != "" {
		return errors.New("the --basic-auth-only option is only supported when username and password/API key are provided")
	}
	return nil
}

func getOffsetAndLimitValues(c *cli.Context) (offset, limit int, err error) {
	offset, err = cliutils.GetIntFlagValue(c, "offset", 0)
	if err != nil {
		return 0, 0, err
	}
	limit, err = cliutils.GetIntFlagValue(c, "limit", 0)
	if err != nil {
		return 0, 0, err
	}

	return
}

func isFailNoOp(context *cli.Context) bool {
	if context == nil {
		return false
	}
	return context.Bool("fail-no-op")
}

// Returns build configuration struct using the params provided from the console.
func createBuildConfiguration(c *cli.Context) *utils.BuildConfiguration {
	buildConfiguration := new(utils.BuildConfiguration)
	buildNameArg, buildNumberArg := c.Args().Get(0), c.Args().Get(1)
	if buildNameArg == "" || buildNumberArg == "" {
		buildNameArg = ""
		buildNumberArg = ""
	}
	buildConfiguration.BuildName, buildConfiguration.BuildNumber = utils.GetBuildNameAndNumber(buildNameArg, buildNumberArg)
	buildConfiguration.Project = utils.GetBuildProject(c.String("project"))
	return buildConfiguration
}

func extractThreadsFlag(args []string) (cleanArgs []string, threadsCount int, err error) {
	// Extract threads flag.
	cleanArgs = append([]string(nil), args...)
	threadsFlagIndex, threadsValueIndex, threads, err := coreutils.FindFlag("--threads", cleanArgs)
	if err != nil || threadsFlagIndex < 0 {
		return
	}
	coreutils.RemoveFlagFromCommand(&cleanArgs, threadsFlagIndex, threadsValueIndex)

	// Convert flag value to int.
	threadsCount, err = strconv.Atoi(threads)
	if err != nil {
		err = errors.New("The '--threads' option should have a numeric value. " + cliutils.GetDocumentationMessage())
	}

	return
}
