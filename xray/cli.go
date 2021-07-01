package xray

import (
	"errors"
	"time"

	"github.com/jfrog/jfrog-cli-core/utils/config"

	"github.com/codegangsta/cli"
	"github.com/jfrog/jfrog-cli-core/common/commands"
	corecommon "github.com/jfrog/jfrog-cli-core/common/commands"
	corecommondocs "github.com/jfrog/jfrog-cli-core/docs/common"
	"github.com/jfrog/jfrog-cli-core/xray/commands/audit"
	"github.com/jfrog/jfrog-cli-core/xray/commands/curl"
	"github.com/jfrog/jfrog-cli-core/xray/commands/offlineupdate"
	"github.com/jfrog/jfrog-cli/docs/common"
	"github.com/jfrog/jfrog-cli/docs/xray/auditgradle"
	"github.com/jfrog/jfrog-cli/docs/xray/auditmvn"
	curldocs "github.com/jfrog/jfrog-cli/docs/xray/curl"
	offlineupdatedocs "github.com/jfrog/jfrog-cli/docs/xray/offlineupdate"
	"github.com/jfrog/jfrog-cli/utils/cliutils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const DATE_FORMAT = "2006-01-02"

func GetCommands() []cli.Command {
	return cliutils.GetSortedCommands(cli.CommandsByName{
		{
			Name:            "curl",
			Flags:           cliutils.GetCommandFlags(cliutils.XrCurl),
			Aliases:         []string{"cl"},
			Description:     curldocs.Description,
			HelpName:        corecommondocs.CreateUsage("xr curl", curldocs.Description, curldocs.Usage),
			UsageText:       curldocs.Arguments,
			ArgsUsage:       common.CreateEnvVars(),
			BashComplete:    corecommondocs.CreateBashCompletionFunc(),
			SkipFlagParsing: true,
			Action:          curlCmd,
		},
		{
			Name:         "audit-mvn",
			Flags:        cliutils.GetCommandFlags(cliutils.AuditMvn),
			Aliases:      []string{"am"},
			Description:  auditmvn.Description,
			HelpName:     corecommondocs.CreateUsage("xr audit-mvn", auditmvn.Description, auditmvn.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       auditMvnCmd,
		},
		{
			Name:         "audit-gradle",
			Flags:        cliutils.GetCommandFlags(cliutils.AuditGradle),
			Aliases:      []string{"ag"},
			Description:  auditgradle.Description,
			HelpName:     corecommondocs.CreateUsage("xr audit-gradle", auditgradle.Description, auditgradle.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       auditGradleCmd,
		},
		{
			Name:         "offline-update",
			Description:  offlineupdatedocs.Description,
			HelpName:     corecommondocs.CreateUsage("xr offline-update", offlineupdatedocs.Description, offlineupdatedocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			Flags:        cliutils.GetCommandFlags(cliutils.OfflineUpdate),
			Aliases:      []string{"ou"},
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       offlineUpdates,
		},
	})
}

func getOfflineUpdatesFlag(c *cli.Context) (flags *offlineupdate.OfflineUpdatesFlags, err error) {
	flags = new(offlineupdate.OfflineUpdatesFlags)
	flags.Version = c.String("version")
	flags.License = c.String("license-id")
	flags.Target = c.String("target")
	if len(flags.License) < 1 {
		return nil, errors.New("The --license-id option is mandatory.")
	}
	from := c.String("from")
	to := c.String("to")
	if len(to) > 0 && len(from) < 1 {
		return nil, errors.New("The --from option is mandatory, when the --to option is sent.")
	}
	if len(from) > 0 && len(to) < 1 {
		return nil, errors.New("The --to option is mandatory, when the --from option is sent.")
	}
	if len(from) > 0 && len(to) > 0 {
		flags.From, err = dateToMilliseconds(from)
		errorutils.CheckError(err)
		if err != nil {
			return
		}
		flags.To, err = dateToMilliseconds(to)
		errorutils.CheckError(err)
	}
	return
}

func dateToMilliseconds(date string) (dateInMillisecond int64, err error) {
	t, err := time.Parse(DATE_FORMAT, date)
	if err != nil {
		errorutils.CheckError(err)
		return
	}
	dateInMillisecond = t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	return
}

func offlineUpdates(c *cli.Context) error {
	offlineUpdateFlags, err := getOfflineUpdatesFlag(c)
	if err != nil {
		return err
	}

	return offlineupdate.OfflineUpdate(offlineUpdateFlags)
}

func curlCmd(c *cli.Context) error {
	if c.NArg() < 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	xrCurlCmd, err := newXrCurlCommand(c)
	if err != nil {
		return err
	}
	return commands.Exec(xrCurlCmd)
}

func newXrCurlCommand(c *cli.Context) (*curl.XrCurlCommand, error) {
	curlCommand := corecommon.NewCurlCommand().SetArguments(cliutils.ExtractCommand(c))
	xrCurlCommand := curl.NewXrCurlCommand(*curlCommand)
	xrDetails, err := xrCurlCommand.GetServerDetails()
	if err != nil {
		return nil, err
	}
	xrCurlCommand.SetServerDetails(xrDetails)
	xrCurlCommand.SetUrl(xrDetails.XrayUrl)
	return xrCurlCommand, err
}

func auditMvnCmd(c *cli.Context) error {
	serverDetailes, err := createXrayServerDetails(c)
	if err != nil {
		return err
	}
	xrAuditNpmCmd := audit.NewAuditMvnCommand().SetExcludeTestDeps(c.Bool(cliutils.ExcludeTestDeps)).SetInsecureTls(c.Bool(cliutils.InsecureTls)).SetServerDetails(serverDetailes)
	return commands.Exec(xrAuditNpmCmd)
}

func auditGradleCmd(c *cli.Context) error {
	serverDetailes, err := createXrayServerDetails(c)
	if err != nil {
		return err
	}
	xrAuditNpmCmd := audit.NewAuditGradleCommand().SetExcludeTestDeps(c.Bool(cliutils.ExcludeTestDeps)).SetUseWrapper(c.Bool(cliutils.UseWrapper)).SetServerDetails(serverDetailes)
	return commands.Exec(xrAuditNpmCmd)
}

func createXrayServerDetails(c *cli.Context) (*config.ServerDetails, error) {
	return cliutils.CreateServerDetailsWithConfigOffer(c, false, cliutils.Xr)
}
