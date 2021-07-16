package cliutils

import (
	"sort"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	// Artifactory's Commands Keys
	DeleteConfig            = "delete-config"
	Upload                  = "upload"
	Download                = "download"
	Move                    = "move"
	Copy                    = "copy"
	Delete                  = "delete"
	Properties              = "properties"
	Search                  = "search"
	BuildPublish            = "build-publish"
	BuildAppend             = "build-append"
	BuildScan               = "build-scan"
	BuildPromote            = "build-promote"
	BuildDiscard            = "build-discard"
	BuildAddDependencies    = "build-add-dependencies"
	BuildAddGit             = "build-add-git"
	BuildCollectEnv         = "build-collect-env"
	GitLfsClean             = "git-lfs-clean"
	Mvn                     = "mvn"
	MvnConfig               = "mvn-config"
	Gradle                  = "gradle"
	GradleConfig            = "gradle-config"
	DockerPromote           = "docker-promote"
	ContainerPull           = "container-pull"
	ContainerPush           = "container-push"
	BuildDockerCreate       = "build-docker-create"
	NpmConfig               = "npm-config"
	Npm                     = "npm"
	NpmPublish              = "npmPublish"
	YarnConfig              = "yarn-config"
	Yarn                    = "yarn"
	NugetConfig             = "nuget-config"
	Nuget                   = "nuget"
	Dotnet                  = "dotnet"
	DotnetConfig            = "dotnet-config"
	Go                      = "go"
	GoConfig                = "go-config"
	GoPublish               = "go-publish"
	PipInstall              = "pip-install"
	PipConfig               = "pip-config"
	Ping                    = "ping"
	RtCurl                  = "rt-curl"
	ReleaseBundleCreate     = "release-bundle-create"
	ReleaseBundleUpdate     = "release-bundle-update"
	ReleaseBundleSign       = "release-bundle-sign"
	ReleaseBundleDistribute = "release-bundle-distribute"
	ReleaseBundleDelete     = "release-bundle-delete"
	TemplateConsumer        = "template-consumer"
	RepoDelete              = "repo-delete"
	ReplicationDelete       = "replication-delete"
	PermissionTargetDelete  = "permission-target-delete"
	AccessTokenCreate       = "access-token-create"
	UserCreate              = "user-create"
	UsersCreate             = "users-create"
	UsersDelete             = "users-delete"
	GroupCreate             = "group-create"
	GroupAddUsers           = "group-add-users"
	GroupDelete             = "group-delete"
	passphrase              = "passphrase"

	// MC's Commands Keys
	McConfig       = "mc-config"
	LicenseAcquire = "license-acquire"
	LicenseDeploy  = "license-deploy"
	LicenseRelease = "license-release"
	JpdAdd         = "jpd-add"
	JpdDelete      = "jpd-delete"
	// XRay's Commands Keys
	XrCurl        = "xr-curl"
	AuditMvn      = "audit-maven"
	AuditGradle   = "audit-gradle"
	AuditNpm      = "audit-npm"
	XrScan        = "xr-scan"
	OfflineUpdate = "offline-update"

	// Config commands keys
	AddConfig  = "config-add"
	EditConfig = "config-edit"

	// *** Artifactory Commands' flags ***
	// Base flags
	url         = "url"
	distUrl     = "dist-url"
	user        = "user"
	password    = "password"
	accessToken = "access-token"
	serverId    = "server-id"

	// Ssh flags
	sshKeyPath    = "ssh-key-path"
	sshPassPhrase = "ssh-passphrase"

	// Client certification flags
	clientCertPath    = "client-cert-path"
	clientCertKeyPath = "client-cert-key-path"
	InsecureTls       = "insecure-tls"

	// Sort & limit flags
	sortBy    = "sort-by"
	sortOrder = "sort-order"
	limit     = "limit"
	offset    = "offset"

	// Spec flags
	specFlag = "spec"
	specVars = "spec-vars"

	// Build info flags
	buildName   = "build-name"
	buildNumber = "build-number"
	module      = "module"

	// Generic commands flags
	exclusions       = "exclusions"
	recursive        = "recursive"
	flat             = "flat"
	build            = "build"
	excludeArtifacts = "exclude-artifacts"
	includeDeps      = "include-deps"
	regexpFlag       = "regexp"
	retries          = "retries"
	dryRun           = "dry-run"
	explode          = "explode"
	includeDirs      = "include-dirs"
	props            = "props"
	targetProps      = "target-props"
	excludeProps     = "exclude-props"
	failNoOp         = "fail-no-op"
	threads          = "threads"
	syncDeletes      = "sync-deletes"
	quiet            = "quiet"
	bundle           = "bundle"
	archiveEntries   = "archive-entries"
	detailedSummary  = "detailed-summary"
	archive          = "archive"
	syncDeletesQuiet = syncDeletes + "-" + quiet
	antFlag          = "ant"
	fromRt           = "from-rt"
	transitive       = "transitive"

	// Config flags
	interactive   = "interactive"
	encPassword   = "enc-password"
	basicAuthOnly = "basic-auth-only"
	overwrite     = "overwrite"

	// Unique upload flags
	uploadPrefix      = "upload-"
	uploadExclusions  = uploadPrefix + exclusions
	uploadRecursive   = uploadPrefix + recursive
	uploadFlat        = uploadPrefix + flat
	uploadRegexp      = uploadPrefix + regexpFlag
	uploadExplode     = uploadPrefix + explode
	uploadTargetProps = uploadPrefix + targetProps
	uploadSyncDeletes = uploadPrefix + syncDeletes
	uploadArchive     = uploadPrefix + archive
	deb               = "deb"
	symlinks          = "symlinks"
	uploadAnt         = uploadPrefix + antFlag

	// Unique download flags
	downloadPrefix       = "download-"
	downloadRecursive    = downloadPrefix + recursive
	downloadFlat         = downloadPrefix + flat
	downloadExplode      = downloadPrefix + explode
	downloadProps        = downloadPrefix + props
	downloadExcludeProps = downloadPrefix + excludeProps
	downloadSyncDeletes  = downloadPrefix + syncDeletes
	minSplit             = "min-split"
	splitCount           = "split-count"
	validateSymlinks     = "validate-symlinks"

	// Unique move flags
	movePrefix       = "move-"
	moveRecursive    = movePrefix + recursive
	moveFlat         = movePrefix + flat
	moveProps        = movePrefix + props
	moveExcludeProps = movePrefix + excludeProps

	// Unique copy flags
	copyPrefix       = "copy-"
	copyRecursive    = copyPrefix + recursive
	copyFlat         = copyPrefix + flat
	copyProps        = copyPrefix + props
	copyExcludeProps = copyPrefix + excludeProps

	// Unique delete flags
	deletePrefix       = "delete-"
	deleteRecursive    = deletePrefix + recursive
	deleteProps        = deletePrefix + props
	deleteExcludeProps = deletePrefix + excludeProps
	deleteQuiet        = deletePrefix + quiet

	// Unique search flags
	searchPrefix       = "search-"
	searchRecursive    = searchPrefix + recursive
	searchProps        = searchPrefix + props
	searchExcludeProps = searchPrefix + excludeProps
	count              = "count"
	searchTransitive   = searchPrefix + transitive

	// Unique properties flags
	propertiesPrefix  = "props-"
	propsRecursive    = propertiesPrefix + recursive
	propsProps        = propertiesPrefix + props
	propsExcludeProps = propertiesPrefix + excludeProps

	// Unique build-publish flags
	buildPublishPrefix = "bp-"
	bpDryRun           = buildPublishPrefix + dryRun
	bpDetailedSummary  = buildPublishPrefix + detailedSummary
	envInclude         = "env-include"
	envExclude         = "env-exclude"
	buildUrl           = "build-url"
	project            = "project"

	// Unique build-add-dependencies flags
	badPrefix    = "bad-"
	badDryRun    = badPrefix + dryRun
	badRecursive = badPrefix + recursive
	badRegexp    = badPrefix + regexpFlag
	badFromRt    = badPrefix + fromRt

	// Unique build-add-git flags
	configFlag = "config"

	// Unique build-scan flags
	fail = "fail"

	// Unique build-promote flags
	buildPromotePrefix  = "bpr-"
	bprDryRun           = buildPromotePrefix + dryRun
	bprProps            = buildPromotePrefix + props
	status              = "status"
	comment             = "comment"
	sourceRepo          = "source-repo"
	includeDependencies = "include-dependencies"
	copyFlag            = "copy"
	failFast            = "fail-fast"

	async = "async"

	// Unique build-discard flags
	buildDiscardPrefix = "bdi-"
	bdiAsync           = buildDiscardPrefix + async
	maxDays            = "max-days"
	maxBuilds          = "max-builds"
	excludeBuilds      = "exclude-builds"
	deleteArtifacts    = "delete-artifacts"

	repo = "repo"

	// Unique git-lfs-clean flags
	glcPrefix = "glc-"
	glcDryRun = glcPrefix + dryRun
	glcQuiet  = glcPrefix + quiet
	glcRepo   = glcPrefix + repo
	refs      = "refs"

	// Build tool config flags
	global          = "global"
	serverIdResolve = "server-id-resolve"
	serverIdDeploy  = "server-id-deploy"
	repoResolve     = "repo-resolve"
	repoDeploy      = "repo-deploy"

	// Unique maven-config flags
	repoResolveReleases  = "repo-resolve-releases"
	repoResolveSnapshots = "repo-resolve-snapshots"
	repoDeployReleases   = "repo-deploy-releases"
	repoDeploySnapshots  = "repo-deploy-snapshots"

	// Unique gradle-config flags
	usesPlugin          = "uses-plugin"
	UseWrapper          = "use-wrapper"
	deployMavenDesc     = "deploy-maven-desc"
	deployIvyDesc       = "deploy-ivy-desc"
	ivyDescPattern      = "ivy-desc-pattern"
	ivyArtifactsPattern = "ivy-artifacts-pattern"

	// Build tool flags
	deploymentThreads = "deployment-threads"
	skipLogin         = "skip-login"

	// Unique docker promote flags
	dockerPromotePrefix = "docker-promote-"
	targetDockerImage   = "target-docker-image"
	sourceTag           = "source-tag"
	targetTag           = "target-tag"
	dockerPromoteCopy   = dockerPromotePrefix + Copy

	// Unique build docker create
	imageFile = "image-file"

	// Unique npm flags
	npmPrefix          = "npm-"
	npmThreads         = npmPrefix + threads
	npmDetailedSummary = npmPrefix + detailedSummary

	// Unique nuget/dotnet config flags
	nugetV2 = "nuget-v2"

	// Unique release-bundle flags
	releaseBundlePrefix = "rb-"
	rbDryRun            = releaseBundlePrefix + dryRun
	rbRepo              = releaseBundlePrefix + repo
	rbPassphrase        = releaseBundlePrefix + passphrase
	distTarget          = releaseBundlePrefix + target
	rbDetailedSummary   = releaseBundlePrefix + detailedSummary
	sign                = "sign"
	desc                = "desc"
	releaseNotesPath    = "release-notes-path"
	releaseNotesSyntax  = "release-notes-syntax"
	distRules           = "dist-rules"
	site                = "site"
	city                = "city"
	countryCodes        = "country-codes"
	sync                = "sync"
	maxWaitMinutes      = "max-wait-minutes"
	deleteFromDist      = "delete-from-dist"

	// Template user flags
	vars = "vars"

	// User Management flags
	csv            = "csv"
	usersCreateCsv = "users-create-csv"
	usersDeleteCsv = "users-delete-csv"
	UsersGroups    = "users-groups"
	Replace        = "replace"
	Admin          = "admin"

	// Unique access-token-create flags
	groups      = "groups"
	grantAdmin  = "grant-admin"
	expiry      = "expiry"
	refreshable = "refreshable"
	audience    = "audience"

	// Unique Xray Flags for upload/publish commands
	xrayScan = "scan"

	// *** Xray Commands' flags ***
	// Unique offline-update flags
	licenseId = "license-id"
	from      = "from"
	to        = "to"
	version   = "version"
	target    = "target"

	// Audit commands
	ExcludeTestDeps = "exclude-test-deps"
	depType         = "dep-type"
	watches         = "watches"
	repoPath        = "repo-path"
	licenses        = "licenses"

	// *** Mission Control Commands' flags ***
	missionControlPrefix = "mc-"

	// Authentication flags
	mcUrl         = missionControlPrefix + url
	mcAccessToken = missionControlPrefix + accessToken

	// Unique config flags
	mcInteractive = missionControlPrefix + interactive

	// Unique license-deploy flags
	licenseCount = "license-count"

	// *** Config Commands' flags ***
	configPrefix      = "config-"
	configPlatformUrl = configPrefix + url
	configRtUrl       = "artifactory-url"
	configXrUrl       = "xray-url"
	configMcUrl       = "mission-control-url"
	configPlUrl       = "pipelines-url"
	configAccessToken = configPrefix + accessToken
	configUser        = configPrefix + user
	configPassword    = configPrefix + password
	configInsecureTls = configPrefix + InsecureTls
)

var flagsMap = map[string]cli.Flag{
	// Artifactory's commands Flags
	url: cli.StringFlag{
		Name:  url,
		Usage: "[Optional] Artifactory URL.` `",
	},
	distUrl: cli.StringFlag{
		Name:  distUrl,
		Usage: "[Optional] Distribution URL.` `",
	},
	user: cli.StringFlag{
		Name:  user,
		Usage: "[Optional] Artifactory username.` `",
	},
	password: cli.StringFlag{
		Name:  password,
		Usage: "[Optional] Artifactory password.` `",
	},
	accessToken: cli.StringFlag{
		Name:  accessToken,
		Usage: "[Optional] Artifactory access token.` `",
	},
	serverId: cli.StringFlag{
		Name:  serverId,
		Usage: "[Optional] Server ID configured using the config command.` `",
	},
	sshKeyPath: cli.StringFlag{
		Name:  sshKeyPath,
		Usage: "[Optional] SSH key file path.` `",
	},
	sshPassPhrase: cli.StringFlag{
		Name:  sshPassPhrase,
		Usage: "[Optional] SSH key passphrase.` `",
	},
	clientCertPath: cli.StringFlag{
		Name:  clientCertPath,
		Usage: "[Optional] Client certificate file in PEM format.` `",
	},
	clientCertKeyPath: cli.StringFlag{
		Name:  clientCertKeyPath,
		Usage: "[Optional] Private key file for the client certificate in PEM format.` `",
	},
	sortBy: cli.StringFlag{
		Name:  sortBy,
		Usage: "[Optional] A list of semicolon-separated fields to sort by. The fields must be part of the 'items' AQL domain. For more information, see https://www.jfrog.com/confluence/display/RTF/Artifactory+Query+Language#ArtifactoryQueryLanguage-EntitiesandFields` `",
	},
	sortOrder: cli.StringFlag{
		Name:  sortOrder,
		Usage: "[Default: asc] The order by which fields in the 'sort-by' option should be sorted. Accepts 'asc' or 'desc'.` `",
	},
	limit: cli.StringFlag{
		Name:  limit,
		Usage: "[Optional] The maximum number of items to fetch. Usually used with the 'sort-by' option.` `",
	},
	offset: cli.StringFlag{
		Name:  offset,
		Usage: "[Optional] The offset from which to fetch items (i.e. how many items should be skipped). Usually used with the 'sort-by' option.` `",
	},
	specFlag: cli.StringFlag{
		Name:  specFlag,
		Usage: "[Optional] Path to a File Spec.` `",
	},
	specVars: cli.StringFlag{
		Name:  specVars,
		Usage: "[Optional] List of variables in the form of \"key1=value1;key2=value2;...\" to be replaced in the File Spec. In the File Spec, the variables should be used as follows: ${key1}.` `",
	},
	buildName: cli.StringFlag{
		Name:  buildName,
		Usage: "[Optional] Providing this option will collect and record build info for this build name. Build number option is mandatory when this option is provided.` `",
	},
	buildNumber: cli.StringFlag{
		Name:  buildNumber,
		Usage: "[Optional] Providing this option will collect and record build info for this build number. Build name option is mandatory when this option is provided.` `",
	},
	module: cli.StringFlag{
		Name:  module,
		Usage: "[Optional] Optional module name for the build-info. Build name and number options are mandatory when this option is provided.` `",
	},
	exclusions: cli.StringFlag{
		Name:  exclusions,
		Usage: "[Optional] Semicolon-separated list of exclusions. Exclusions can include the * and the ? wildcards.` `",
	},
	uploadExclusions: cli.StringFlag{
		Name:  exclusions,
		Usage: "[Optional] Semicolon-separated list of exclude patterns. Exclude patterns may contain the * and the ? wildcards or a regex pattern, according to the value of the 'regexp' option.` `",
	},
	build: cli.StringFlag{
		Name:  build,
		Usage: "[Optional] If specified, only artifacts of the specified build are matched. The property format is build-name/build-number. If you do not specify the build number, the artifacts are filtered by the latest build number.` `",
	},
	excludeArtifacts: cli.StringFlag{
		Name:  excludeArtifacts,
		Usage: "[Default: false] If specified, build artifacts are not matched. Used together with the --build flag.` `",
	},
	includeDeps: cli.StringFlag{
		Name:  includeDeps,
		Usage: "[Default: false] If specified, also dependencies of the specified build are matched. Used together with the --build flag.` `",
	},
	includeDirs: cli.BoolFlag{
		Name:  includeDirs,
		Usage: "[Default: false] Set to true if you'd like to also apply the source path pattern for directories and not just for files.` `",
	},
	failNoOp: cli.BoolFlag{
		Name:  failNoOp,
		Usage: "[Default: false] Set to true if you'd like the command to return exit code 2 in case of no files are affected.` `",
	},
	threads: cli.StringFlag{
		Name:  threads,
		Value: "",
		Usage: "[Default: " + strconv.Itoa(Threads) + "] Number of working threads.` `",
	},
	retries: cli.StringFlag{
		Name:  retries,
		Usage: "[Default: " + strconv.Itoa(Retries) + "] Number of HTTP retries.` `",
	},
	InsecureTls: cli.BoolFlag{
		Name:  InsecureTls,
		Usage: "[Default: false] Set to true to skip TLS certificates verification.` `",
	},
	bundle: cli.StringFlag{
		Name:  bundle,
		Usage: "[Optional] If specified, only artifacts of the specified bundle are matched. The value format is bundle-name/bundle-version.` `",
	},
	archiveEntries: cli.StringFlag{
		Name:  archiveEntries,
		Usage: "[Optional] If specified, only archive artifacts containing entries matching this pattern are matched. You can use wildcards to specify multiple artifacts.` `",
	},
	detailedSummary: cli.BoolFlag{
		Name:  detailedSummary,
		Usage: "[Default: false] Set to true to include a list of the affected files in the command summary.` `",
	},
	interactive: cli.BoolTFlag{
		Name:  interactive,
		Usage: "[Default: true, unless $CI is true] Set to false if you do not want the config command to be interactive. If true, the --url option becomes optional.` `",
	},
	encPassword: cli.BoolTFlag{
		Name:  encPassword,
		Usage: "[Default: true] If set to false then the configured password will not be encrypted using Artifactory's encryption API.` `",
	},
	overwrite: cli.BoolFlag{
		Name:  overwrite,
		Usage: "[Default: false] Overwrites the instance configuration if an instance with the same ID already exists.` `",
	},
	basicAuthOnly: cli.BoolFlag{
		Name: basicAuthOnly,
		Usage: "[Default: false] Set to true to disable replacing username and password/API key with automatically created access token that's refreshed hourly. " +
			"Username and password/API key will still be used with commands which use external tools or the JFrog Distribution service. " +
			"Can only be passed along with username and password/API key options.` `",
	},
	deb: cli.StringFlag{
		Name:  deb,
		Usage: "[Optional] Used for Debian packages in the form of distribution/component/architecture. If the value for distribution, component or architecture includes a slash, the slash should be escaped with a back-slash.` `",
	},
	uploadRecursive: cli.BoolTFlag{
		Name:  recursive,
		Usage: "[Default: true] Set to false if you do not wish to collect artifacts in sub-folders to be uploaded to Artifactory.` `",
	},
	uploadFlat: cli.BoolFlag{
		Name:  flat,
		Usage: "[Default: false] If set to false, files are uploaded according to their file system hierarchy.` `",
	},
	uploadRegexp: cli.BoolFlag{
		Name:  regexpFlag,
		Usage: "[Default: false] Set to true to use a regular expression instead of wildcards expression to collect files to upload.` `",
	},
	uploadAnt: cli.BoolFlag{
		Name:  antFlag,
		Usage: "[Default: false] Set to true to use an ant pattern instead of wildcards expression to collect files to upload.` `",
	},
	dryRun: cli.BoolFlag{
		Name:  dryRun,
		Usage: "[Default: false] Set to true to disable communication with Artifactory.` `",
	},
	uploadExplode: cli.BoolFlag{
		Name:  explode,
		Usage: "[Default: false] Set to true to extract an archive after it is deployed to Artifactory.` `",
	},
	symlinks: cli.BoolFlag{
		Name:  symlinks,
		Usage: "[Default: false] Set to true to preserve symbolic links structure in Artifactory.` `",
	},
	uploadTargetProps: cli.StringFlag{
		Name:  targetProps,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Those properties will be attached to the uploaded artifacts.` `",
	},
	uploadSyncDeletes: cli.StringFlag{
		Name:  syncDeletes,
		Usage: "[Optional] Specific path in Artifactory, under which to sync artifacts after the upload. After the upload, this path will include only the artifacts uploaded during this upload operation. The other files under this path will be deleted.` `",
	},
	uploadArchive: cli.StringFlag{
		Name:  archive,
		Usage: "[Optional] Set to \"zip\" to deploy the files to Artifactory in a ZIP archive.` `",
	},
	syncDeletesQuiet: cli.BoolFlag{
		Name:  quiet,
		Usage: "[Default: $CI] Set to true to skip the sync-deletes confirmation message.` `",
	},
	downloadRecursive: cli.BoolTFlag{
		Name:  recursive,
		Usage: "[Default: true] Set to false if you do not wish to include the download of artifacts inside sub-folders in Artifactory.` `",
	},
	downloadFlat: cli.BoolFlag{
		Name:  flat,
		Usage: "[Default: false] Set to true if you do not wish to have the Artifactory repository path structure created locally for your downloaded files.` `",
	},
	minSplit: cli.StringFlag{
		Name:  minSplit,
		Value: "",
		Usage: "[Default: " + strconv.Itoa(DownloadMinSplitKb) + "] Minimum file size in KB to split into ranges when downloading. Set to -1 for no splits.` `",
	},
	splitCount: cli.StringFlag{
		Name:  splitCount,
		Value: "",
		Usage: "[Default: " + strconv.Itoa(DownloadSplitCount) + "] Number of parts to split a file when downloading. Set to 0 for no splits.` `",
	},
	downloadExplode: cli.BoolFlag{
		Name:  explode,
		Usage: "[Default: false] Set to true to extract an archive after it is downloaded from Artifactory.` `",
	},
	validateSymlinks: cli.BoolFlag{
		Name:  validateSymlinks,
		Usage: "[Default: false] Set to true to perform a checksum validation when downloading symbolic links.` `",
	},
	downloadProps: cli.StringFlag{
		Name:  props,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts with these properties will be downloaded.` `",
	},
	downloadExcludeProps: cli.StringFlag{
		Name:  excludeProps,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts without the specified properties will be downloaded.` `",
	},
	downloadSyncDeletes: cli.StringFlag{
		Name:  syncDeletes,
		Usage: "[Optional] Specific path in the local file system, under which to sync dependencies after the download. After the download, this path will include only the dependencies downloaded during this download operation. The other files under this path will be deleted.` `",
	},
	moveRecursive: cli.BoolTFlag{
		Name:  recursive,
		Usage: "[Default: true] Set to false if you do not wish to move artifacts inside sub-folders in Artifactory.` `",
	},
	moveFlat: cli.BoolFlag{
		Name:  flat,
		Usage: "[Default: false] If set to false, files are moved according to their file system hierarchy.` `",
	},
	moveProps: cli.StringFlag{
		Name:  props,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts with these properties will be moved.` `",
	},
	moveExcludeProps: cli.StringFlag{
		Name:  excludeProps,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts without the specified properties will be moved.` `",
	},
	copyRecursive: cli.BoolTFlag{
		Name:  recursive,
		Usage: "[Default: true] Set to false if you do not wish to copy artifacts inside sub-folders in Artifactory.` `",
	},
	copyFlat: cli.BoolFlag{
		Name:  flat,
		Usage: "[Default: false] If set to false, files are copied according to their file system hierarchy.` `",
	},
	copyProps: cli.StringFlag{
		Name:  props,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts with these properties will be copied.` `",
	},
	copyExcludeProps: cli.StringFlag{
		Name:  excludeProps,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts without the specified properties will be copied.` `",
	},
	deleteRecursive: cli.BoolTFlag{
		Name:  recursive,
		Usage: "[Default: true] Set to false if you do not wish to delete artifacts inside sub-folders in Artifactory.` `",
	},
	deleteProps: cli.StringFlag{
		Name:  props,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts with these properties will be deleted.` `",
	},
	deleteExcludeProps: cli.StringFlag{
		Name:  excludeProps,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts without the specified properties will be deleted.` `",
	},
	deleteQuiet: cli.BoolFlag{
		Name:  quiet,
		Usage: "[Default: $CI] Set to true to skip the delete confirmation message.` `",
	},
	searchRecursive: cli.BoolTFlag{
		Name:  recursive,
		Usage: "[Default: true] Set to false if you do not wish to search artifacts inside sub-folders in Artifactory.` `",
	},
	count: cli.BoolFlag{
		Name:  count,
		Usage: "[Optional] Set to true to display only the total of files or folders found.` `",
	},
	searchProps: cli.StringFlag{
		Name:  props,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts with these properties will be returned.` `",
	},
	searchExcludeProps: cli.StringFlag{
		Name:  excludeProps,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts without the specified properties will be returned` `",
	},
	searchTransitive: cli.BoolFlag{
		Name:  transitive,
		Usage: "[Default: false] Set to true to look for artifacts also in remote repositories. Available on Artifactory version 7.17.0 or higher.` `",
	},
	propsRecursive: cli.BoolTFlag{
		Name:  recursive,
		Usage: "[Default: true] When false, artifacts inside sub-folders in Artifactory will not be affected.` `",
	},
	propsProps: cli.StringFlag{
		Name:  props,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts with these properties are affected.` `",
	},
	propsExcludeProps: cli.StringFlag{
		Name:  excludeProps,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". Only artifacts without the specified properties are affected` `",
	},
	buildUrl: cli.StringFlag{
		Name:  buildUrl,
		Usage: "[Optional] Can be used for setting the CI server build URL in the build-info.` `",
	},
	project: cli.StringFlag{
		Name:  project,
		Usage: "[Optional] Artifactory project key.` `",
	},
	bpDryRun: cli.BoolFlag{
		Name:  dryRun,
		Usage: "[Default: false] Set to true to get a preview of the recorded build info, without publishing it to Artifactory.` `",
	},
	bpDetailedSummary: cli.BoolFlag{
		Name:  detailedSummary,
		Usage: "[Default: false] Set to true to get a command summary with details about the build info artifact.` `",
	},
	envInclude: cli.StringFlag{
		Name:  envInclude,
		Usage: "[Default: *] List of patterns in the form of \"value1;value2;...\" Only environment variables match those patterns will be included.` `",
	},
	envExclude: cli.StringFlag{
		Name:  envExclude,
		Usage: "[Default: *password*;*psw*;*secret*;*key*;*token*] List of case insensitive patterns in the form of \"value1;value2;...\". Environment variables match those patterns will be excluded.` `",
	},
	badRecursive: cli.BoolTFlag{
		Name:  recursive,
		Usage: "[Default: true] Set to false if you do not wish to collect artifacts in sub-folders to be added to the build info.` `",
	},
	badRegexp: cli.BoolFlag{
		Name:  regexpFlag,
		Usage: "[Default: false] Set to true to use a regular expression instead of wildcards expression to collect files to be added to the build info.` `",
	},
	badDryRun: cli.BoolFlag{
		Name:  dryRun,
		Usage: "[Default: false] Set to true to only get a summery of the dependencies that will be added to the build info.` `",
	},
	badFromRt: cli.BoolFlag{
		Name:  fromRt,
		Usage: "[Default: false] Set true to search the files in Artifactory, rather than on the local file system. The --regexp option is not supported when --from-rt is set to true.` `",
	},
	configFlag: cli.StringFlag{
		Name:  configFlag,
		Usage: "[Optional] Path to a configuration file.` `",
	},
	fail: cli.BoolTFlag{
		Name:  fail,
		Usage: "[Default: true] Set to false if you do not wish the command to return exit code 3, even if the 'Fail Build' rule is matched by Xray.` `",
	},
	status: cli.StringFlag{
		Name:  status,
		Usage: "[Optional] Build promotion status.` `",
	},
	comment: cli.StringFlag{
		Name:  comment,
		Usage: "[Optional] Build promotion comment.` `",
	},
	sourceRepo: cli.StringFlag{
		Name:  sourceRepo,
		Usage: "[Optional] Build promotion source repository.` `",
	},
	includeDependencies: cli.BoolFlag{
		Name:  includeDependencies,
		Usage: "[Default: false] If set to true, the build dependencies are also promoted.` `",
	},
	copyFlag: cli.BoolFlag{
		Name:  copyFlag,
		Usage: "[Default: false] If set true, the build artifacts and dependencies are copied to the target repository, otherwise they are moved.` `",
	},
	failFast: cli.BoolTFlag{
		Name:  failFast,
		Usage: "[Default: true] If set true, fail and abort the operation upon receiving an error.` `",
	},
	bprDryRun: cli.BoolFlag{
		Name:  dryRun,
		Usage: "[Default: false] If true, promotion is only simulated. The build is not promoted.` `",
	},
	bprProps: cli.StringFlag{
		Name:  props,
		Usage: "[Optional] List of properties in the form of \"key1=value1;key2=value2,...\". A list of properties to attach to the build artifacts.` `",
	},
	targetDockerImage: cli.StringFlag{
		Name:  "target-docker-image",
		Usage: "[Optional] Docker target image name.` `",
	},
	sourceTag: cli.StringFlag{
		Name:  "source-tag",
		Usage: "[Optional] The tag name to promote.` `",
	},
	targetTag: cli.StringFlag{
		Name:  "target-tag",
		Usage: "[Optional] The target tag to assign the image after promotion.` `",
	},
	dockerPromoteCopy: cli.BoolFlag{
		Name:  "copy",
		Usage: "[Default: false] If set true, the Docker image is copied to the target repository, otherwise it is moved.` `",
	},
	maxDays: cli.StringFlag{
		Name:  maxDays,
		Usage: "[Optional] The maximum number of days to keep builds in Artifactory.` `",
	},
	maxBuilds: cli.StringFlag{
		Name:  maxBuilds,
		Usage: "[Optional] The maximum number of builds to store in Artifactory.` `",
	},
	excludeBuilds: cli.StringFlag{
		Name:  excludeBuilds,
		Usage: "[Optional] List of build numbers in the form of \"value1,value2,...\", that should not be removed from Artifactory.` `",
	},
	deleteArtifacts: cli.BoolFlag{
		Name:  deleteArtifacts,
		Usage: "[Default: false] If set to true, automatically removes build artifacts stored in Artifactory.` `",
	},
	bdiAsync: cli.BoolFlag{
		Name:  async,
		Usage: "[Default: false] If set to true, build discard will run asynchronously and will not wait for response.` `",
	},
	refs: cli.StringFlag{
		Name:  refs,
		Usage: "[Default: refs/remotes/*] List of Git references in the form of \"ref1,ref2,...\" which should be preserved.` `",
	},
	glcRepo: cli.StringFlag{
		Name:  repo,
		Usage: "[Optional] Local Git LFS repository which should be cleaned. If omitted, this is detected from the Git repository.` `",
	},
	glcDryRun: cli.BoolFlag{
		Name:  dryRun,
		Usage: "[Default: false] If true, cleanup is only simulated. No files are actually deleted.` `",
	},
	glcQuiet: cli.BoolFlag{
		Name:  quiet,
		Usage: "[Default: $CI] Set to true to skip the delete confirmation message.` `",
	},
	global: cli.BoolFlag{
		Name:  global,
		Usage: "[Default: false] Set to true if you'd like the configuration to be global (for all projects). Specific projects can override the global configuration.` `",
	},
	serverIdResolve: cli.StringFlag{
		Name:  serverIdResolve,
		Usage: "[Optional] Artifactory server ID for resolution. The server should configured using the 'jfrog c add' command.` `",
	},
	serverIdDeploy: cli.StringFlag{
		Name:  serverIdDeploy,
		Usage: "[Optional] Artifactory server ID for deployment. The server should configured using the 'jfrog c add' command.` `",
	},
	repoResolveReleases: cli.StringFlag{
		Name:  repoResolveReleases,
		Usage: "[Optional] Resolution repository for release dependencies.` `",
	},
	repoResolveSnapshots: cli.StringFlag{
		Name:  repoResolveSnapshots,
		Usage: "[Optional] Resolution repository for snapshot dependencies.` `",
	},
	repoDeployReleases: cli.StringFlag{
		Name:  repoDeployReleases,
		Usage: "[Optional] Deployment repository for release artifacts.` `",
	},
	repoDeploySnapshots: cli.StringFlag{
		Name:  repoDeploySnapshots,
		Usage: "[Optional] Deployment repository for snapshot artifacts.` `",
	},
	repoResolve: cli.StringFlag{
		Name:  repoResolve,
		Usage: "[Optional] Repository for dependencies resolution.` `",
	},
	repoDeploy: cli.StringFlag{
		Name:  repoDeploy,
		Usage: "[Optional] Repository for artifacts deployment.` `",
	},
	usesPlugin: cli.BoolFlag{
		Name:  usesPlugin,
		Usage: "[Default: false] Set to true if the Gradle Artifactory Plugin is already applied in the build script.` `",
	},
	UseWrapper: cli.BoolFlag{
		Name:  UseWrapper,
		Usage: "[Default: false] Set to true if you'd like to use the Gradle wrapper.` `",
	},
	deployMavenDesc: cli.BoolTFlag{
		Name:  deployMavenDesc,
		Usage: "[Default: true] Set to false if you do not wish to deploy Maven descriptors.` `",
	},
	deployIvyDesc: cli.BoolTFlag{
		Name:  deployIvyDesc,
		Usage: "[Default: true] Set to false if you do not wish to deploy Ivy descriptors.` `",
	},
	ivyDescPattern: cli.StringFlag{
		Name:  ivyDescPattern,
		Usage: "[Default: '[organization]/[module]/ivy-[revision].xml' Set the deployed Ivy descriptor pattern.` `",
	},
	ivyArtifactsPattern: cli.StringFlag{
		Name:  ivyArtifactsPattern,
		Usage: "[Default: '[organization]/[module]/[revision]/[artifact]-[revision](-[classifier]).[ext]' Set the deployed Ivy artifacts pattern.` `",
	},
	deploymentThreads: cli.StringFlag{
		Name:  threads,
		Value: "",
		Usage: "[Default: " + strconv.Itoa(Threads) + "] Number of threads for uploading build artifacts.` `",
	},
	skipLogin: cli.BoolFlag{
		Name:  skipLogin,
		Usage: "[Default: false] Set to true if you'd like the command to skip performing docker login.` `",
	},
	npmThreads: cli.StringFlag{
		Name:  threads,
		Value: "",
		Usage: "[Default: 3] Number of working threads for build-info collection.` `",
	},
	npmDetailedSummary: cli.BoolFlag{
		Name:  detailedSummary,
		Usage: "[Default: false] Set to true to include a list of the affected files in the command summary.` `",
	},
	nugetV2: cli.BoolFlag{
		Name:  nugetV2,
		Usage: "[Default: false] Set to true if you'd like to use the NuGet V2 protocol when restoring packages from Artifactory.` `",
	},
	rbDryRun: cli.BoolFlag{
		Name:  dryRun,
		Usage: "[Default: false] Set to true to disable communication with JFrog Distribution.` `",
	},
	rbDetailedSummary: cli.BoolFlag{
		Name:  detailedSummary,
		Usage: "[Default: false] Set to true to get a command summary with details about the release bundle artifact.` `",
	},
	sign: cli.BoolFlag{
		Name:  sign,
		Usage: "[Default: false] If set to true, automatically signs the release bundle version.` `",
	},
	desc: cli.StringFlag{
		Name:  desc,
		Usage: "[Optional] Description of the release bundle.` `",
	},
	releaseNotesPath: cli.StringFlag{
		Name:  releaseNotesPath,
		Usage: "[Optional] Path to a file describes the release notes for the release bundle version.` `",
	},
	releaseNotesSyntax: cli.StringFlag{
		Name:  "release-notes-syntax",
		Usage: "[Default: plain_text] The syntax for the release notes. Can be one of 'markdown', 'asciidoc', or 'plain_text` `",
	},
	rbPassphrase: cli.StringFlag{
		Name:  passphrase,
		Usage: "[Optional] The passphrase for the signing key. ` `",
	},
	distTarget: cli.StringFlag{
		Name: target,
		Usage: "[Optional] The target path for distributed artifacts on the edge node. If not specified, the artifacts will have the same path and name on the edge node, as on the source Artifactory server. " +
			"For flexibility in specifying the distribution path, you can include placeholders in the form of {1}, {2} which are replaced by corresponding tokens in the pattern path that are enclosed in parenthesis. ` `",
	},
	rbRepo: cli.StringFlag{
		Name:  repo,
		Usage: "[Optional] A repository name at source Artifactory to store release bundle artifacts in. If not provided, Artifactory will use the default one.` `",
	},
	distRules: cli.StringFlag{
		Name:  distRules,
		Usage: "Path to distribution rules.` `",
	},
	site: cli.StringFlag{
		Name:  site,
		Usage: "[Default: '*'] Wildcard filter for site name. ` `",
	},
	city: cli.StringFlag{
		Name:  city,
		Usage: "[Default: '*'] Wildcard filter for site city name. ` `",
	},
	countryCodes: cli.StringFlag{
		Name:  countryCodes,
		Usage: "[Default: '*'] Semicolon-separated list of wildcard filters for site country codes. ` `",
	},
	sync: cli.BoolFlag{
		Name:  sync,
		Usage: "[Default: false] Set to true to enable sync distribution (the command execution will end when the distribution process ends).` `",
	},
	maxWaitMinutes: cli.StringFlag{
		Name:  maxWaitMinutes,
		Usage: "[Default: 60] Max minutes to wait for sync distribution. ` `",
	},
	deleteFromDist: cli.BoolFlag{
		Name:  deleteFromDist,
		Usage: "[Default: false] Set to true to delete release bundle version in JFrog Distribution itself after deletion is complete in the specified Edge node/s.` `",
	},
	targetProps: cli.StringFlag{
		Name:  targetProps,
		Usage: "[Optional] The list of properties, in the form of key1=value1;key2=value2,..., to be added to the artifacts after distribution of the release bundle.` `",
	},
	vars: cli.StringFlag{
		Name:  vars,
		Usage: "[Optional] List of variables in the form of \"key1=value1;key2=value2;...\" to be replaced in the template. In the template, the variables should be used as follows: ${key1}.` `",
	},
	groups: cli.StringFlag{
		Name: groups,
		Usage: "[Default: *] A list of comma-separated groups for the access token to be associated with. " +
			"Specify * to indicate that this is a 'user-scoped token', i.e., the token provides the same access privileges that the current subject has, and is therefore evaluated dynamically. " +
			"A non-admin user can only provide a scope that is a subset of the groups to which he belongs` `",
	},
	grantAdmin: cli.BoolFlag{
		Name:  grantAdmin,
		Usage: "[Default: false] Set to true to provides admin privileges to the access token. This is only available for administrators.` `",
	},
	expiry: cli.StringFlag{
		Name:  expiry,
		Usage: "[Default: " + strconv.Itoa(TokenExpiry) + "] The time in seconds for which the token will be valid. To specify a token that never expires, set to zero. Non-admin can only set a value that is equal to or less than the default 3600.` `",
	},
	refreshable: cli.BoolFlag{
		Name:  refreshable,
		Usage: "[Default: false] Set to true if you'd like the the token to be refreshable. A refresh token will also be returned in order to be used to generate a new token once it expires.` `",
	},
	audience: cli.StringFlag{
		Name:  audience,
		Usage: "[Optional] A space-separate list of the other Artifactory instances or services that should accept this token identified by their Artifactory Service IDs, as obtained by the 'jfrog rt curl api/system/service_id' command.` `",
	},
	usersCreateCsv: cli.StringFlag{
		Name:  csv,
		Usage: "[Mandatory] Path to a csv file with the users' details. The first row of the file is reserved for the cells' headers. It must include \"username\",\"password\",\"email\"` `",
	},
	usersDeleteCsv: cli.StringFlag{
		Name:  csv,
		Usage: "[Optional] Path to a csv file with the users' details. The first row of the file is reserved for the cells' headers. It must include \"username\"` `",
	},
	UsersGroups: cli.StringFlag{
		Name:  UsersGroups,
		Usage: "[Optional] A list of comma-separated groups for the new users to be associated with.` `",
	},
	Replace: cli.BoolFlag{
		Name:  Replace,
		Usage: "[Default: false] Set to true if you'd like existing users or groups to be replaced.` `",
	},
	Admin: cli.BoolFlag{
		Name:  Admin,
		Usage: "[Default: false] Set to true if you'd like to create an admin user.` `",
	},
	xrayScan: cli.StringFlag{
		Name:  xrayScan,
		Usage: "[Default: false] Set if you'd like all files to be scanned by Xray on the local file system prior to the upload, and skip the upload if any of the files are found vulnerable.` `",
	},
	// Xray's commands Flags
	licenseId: cli.StringFlag{
		Name:  licenseId,
		Usage: "[Mandatory] Xray license ID` `",
	},
	from: cli.StringFlag{
		Name:  from,
		Usage: "[Optional] From update date in YYYY-MM-DD format.` `",
	},
	to: cli.StringFlag{
		Name:  to,
		Usage: "[Optional] To update date in YYYY-MM-DD format.` `",
	},
	version: cli.StringFlag{
		Name:  version,
		Usage: "[Optional] Xray API version.` `",
	},
	target: cli.StringFlag{
		Name:  target,
		Usage: "[Default: ./] Path for downloaded update files.` `",
	},
	ExcludeTestDeps: cli.BoolFlag{
		Name:  ExcludeTestDeps,
		Usage: "[Default: false] Set to true if you'd like to exclude test dependencies from Xray scanning.` `",
	},
	depType: cli.StringFlag{
		Name:  depType,
		Usage: "[Default: all] Defines npm dependencies type. Possible values are: all, devOnly and prodOnly` `",
	},
	watches: cli.StringFlag{
		Name:  watches,
		Usage: "[Optional] A comma separated list of Xray watches, to determine Xray's violations creation. ` `",
	},
	licenses: cli.BoolFlag{
		Name:  licenses,
		Usage: "[Optional] Set to true if you'd like to receive licenses from Xray scanning. ` `",
	},
	repoPath: cli.StringFlag{
		Name:  repoPath,
		Usage: "[Optional] Target repo path, to enable Xray to determine watches accordingly. ` `",
	},
	// Mission Control's commands Flags
	mcUrl: cli.StringFlag{
		Name:  url,
		Usage: "[Optional] Mission Control URL.` `",
	},
	mcAccessToken: cli.StringFlag{
		Name:  accessToken,
		Usage: "[Optional] Mission Control Admin token.` `",
	},
	mcInteractive: cli.BoolTFlag{
		Name:  interactive,
		Usage: "[Default: true] Set to false if you do not want the config command to be interactive. If true, the other command options become optional.",
	},
	licenseCount: cli.StringFlag{
		Name:  licenseCount,
		Value: "",
		Usage: "[Default: " + strconv.Itoa(DefaultLicenseCount) + "] The number of licenses to deploy. Minimum value is 1.` `",
	},
	imageFile: cli.StringFlag{
		Name:  imageFile,
		Usage: "[Mandatory] Path to a file which includes one line in the following format: <IMAGE-TAG>@sha256:<MANIFEST-SHA256>.` `",
	},
	// Config commands Flags
	configPlatformUrl: cli.StringFlag{
		Name:  url,
		Usage: "[Optional] JFrog platform URL.` `",
	},
	configRtUrl: cli.StringFlag{
		Name:  configRtUrl,
		Usage: "[Optional] Artifactory URL.` `",
	},
	configXrUrl: cli.StringFlag{
		Name:  configXrUrl,
		Usage: "[Optional] Xray URL.` `",
	},
	configMcUrl: cli.StringFlag{
		Name:  configMcUrl,
		Usage: "[Optional] Mission Control URL.` `",
	},
	configPlUrl: cli.StringFlag{
		Name:  configPlUrl,
		Usage: "[Optional] Pipelines URL.` `",
	},
	configUser: cli.StringFlag{
		Name:  user,
		Usage: "[Optional] JFrog Platform username. ` `",
	},
	configPassword: cli.StringFlag{
		Name:  password,
		Usage: "[Optional] JFrog Platform password or API key. ` `",
	},
	configAccessToken: cli.StringFlag{
		Name:  accessToken,
		Usage: "[Optional] JFrog Platform access token. ` `",
	},
	configInsecureTls: cli.StringFlag{
		Name:  InsecureTls,
		Usage: "[Default: false] Set to true to skip TLS certificates verification, while encrypting the Artifactory password during the config process.` `",
	},
}

var commandFlags = map[string][]string{
	AddConfig: {
		interactive, encPassword, configPlatformUrl, configRtUrl, distUrl, configXrUrl, configMcUrl, configPlUrl, configUser, configPassword, configAccessToken, sshKeyPath, clientCertPath,
		clientCertKeyPath, basicAuthOnly, configInsecureTls, overwrite,
	},
	EditConfig: {
		interactive, encPassword, configPlatformUrl, configRtUrl, distUrl, configXrUrl, configMcUrl, configPlUrl, configUser, configPassword, configAccessToken, sshKeyPath, clientCertPath,
		clientCertKeyPath, basicAuthOnly, configInsecureTls,
	},
	DeleteConfig: {
		deleteQuiet,
	},
	Upload: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath, uploadTargetProps,
		clientCertKeyPath, specFlag, specVars, buildName, buildNumber, module, uploadExclusions, deb,
		uploadRecursive, uploadFlat, uploadRegexp, retries, dryRun, uploadExplode, symlinks, includeDirs,
		failNoOp, threads, uploadSyncDeletes, syncDeletesQuiet, InsecureTls, detailedSummary, project,
		uploadAnt, uploadArchive,
	},
	Download: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, specFlag, specVars, buildName, buildNumber, module, exclusions, sortBy,
		sortOrder, limit, offset, downloadRecursive, downloadFlat, build, includeDeps, excludeArtifacts, minSplit, splitCount,
		retries, dryRun, downloadExplode, validateSymlinks, bundle, includeDirs, downloadProps, downloadExcludeProps,
		failNoOp, threads, archiveEntries, downloadSyncDeletes, syncDeletesQuiet, InsecureTls, detailedSummary, project,
	},
	Move: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, specFlag, specVars, exclusions, sortBy, sortOrder, limit, offset, moveRecursive,
		moveFlat, dryRun, build, includeDeps, excludeArtifacts, moveProps, moveExcludeProps, failNoOp, threads, archiveEntries,
		InsecureTls, retries,
	},
	Copy: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, specFlag, specVars, exclusions, sortBy, sortOrder, limit, offset, copyRecursive,
		copyFlat, dryRun, build, includeDeps, excludeArtifacts, bundle, copyProps, copyExcludeProps, failNoOp, threads,
		archiveEntries, InsecureTls, retries,
	},
	Delete: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, specFlag, specVars, exclusions, sortBy, sortOrder, limit, offset,
		deleteRecursive, dryRun, build, includeDeps, excludeArtifacts, deleteQuiet, deleteProps, deleteExcludeProps, failNoOp, threads, archiveEntries,
		InsecureTls, retries,
	},
	Search: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, specFlag, specVars, exclusions, sortBy, sortOrder, limit, offset,
		searchRecursive, build, includeDeps, excludeArtifacts, count, bundle, includeDirs, searchProps, searchExcludeProps, failNoOp, archiveEntries,
		InsecureTls, searchTransitive, retries,
	},
	Properties: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, specFlag, specVars, exclusions, sortBy, sortOrder, limit, offset,
		propsRecursive, build, includeDeps, excludeArtifacts, bundle, includeDirs, failNoOp, threads, archiveEntries, propsProps, propsExcludeProps,
		InsecureTls, retries,
	},
	BuildPublish: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, buildUrl, bpDryRun,
		envInclude, envExclude, InsecureTls, project, bpDetailedSummary,
	},
	BuildAppend: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, buildUrl, bpDryRun,
		envInclude, envExclude, InsecureTls, project,
	},
	BuildAddDependencies: {
		specFlag, specVars, uploadExclusions, badRecursive, badRegexp, badDryRun, project, badFromRt, serverId,
	},
	BuildAddGit: {
		configFlag, serverId, project,
	},
	BuildCollectEnv: {
		project,
	},
	BuildDockerCreate: {
		buildName, buildNumber, module, url, user, password, accessToken, sshPassPhrase, sshKeyPath,
		serverId, imageFile, project,
	},
	BuildScan: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, fail, InsecureTls,
		project,
	},
	BuildPromote: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, status, comment,
		sourceRepo, includeDependencies, copyFlag, failFast, bprDryRun, bprProps, InsecureTls, project,
	},
	BuildDiscard: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, maxDays, maxBuilds,
		excludeBuilds, deleteArtifacts, bdiAsync, InsecureTls, project,
	},
	GitLfsClean: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, refs, glcRepo, glcDryRun,
		glcQuiet, InsecureTls, retries,
	},
	MvnConfig: {
		global, serverIdResolve, serverIdDeploy, repoResolveReleases, repoResolveSnapshots, repoDeployReleases, repoDeploySnapshots,
	},
	GradleConfig: {
		global, serverIdResolve, serverIdDeploy, repoResolve, repoDeploy, usesPlugin, UseWrapper, deployMavenDesc,
		deployIvyDesc, ivyDescPattern, ivyArtifactsPattern,
	},
	Mvn: {
		buildName, buildNumber, deploymentThreads, InsecureTls, project, detailedSummary, xrayScan,
	},
	Gradle: {
		buildName, buildNumber, deploymentThreads, project, detailedSummary, xrayScan,
	},
	DockerPromote: {
		targetDockerImage, sourceTag, targetTag, dockerPromoteCopy, url, user, password, accessToken, sshPassPhrase, sshKeyPath,
		serverId,
	},
	ContainerPush: {
		buildName, buildNumber, module, url, user, password, accessToken, sshPassPhrase, sshKeyPath,
		serverId, skipLogin, threads, project, detailedSummary,
	},
	ContainerPull: {
		buildName, buildNumber, module, url, user, password, accessToken, sshPassPhrase, sshKeyPath,
		serverId, skipLogin, project,
	},
	NpmConfig: {
		global, serverIdResolve, serverIdDeploy, repoResolve, repoDeploy,
	},
	Npm: {
		buildName, buildNumber, module, npmThreads, project,
	},
	NpmPublish: {
		buildName, buildNumber, module, project, npmDetailedSummary, xrayScan,
	},
	YarnConfig: {
		global, serverIdResolve, repoResolve,
	},
	Yarn: {
		buildName, buildNumber, module, project,
	},
	NugetConfig: {
		global, serverIdResolve, repoResolve, nugetV2,
	},
	Nuget: {
		buildName, buildNumber, module, project,
	},
	DotnetConfig: {
		global, serverIdResolve, repoResolve, nugetV2,
	},
	Dotnet: {
		buildName, buildNumber, module, project,
	},
	GoConfig: {
		global, serverIdResolve, serverIdDeploy, repoResolve, repoDeploy,
	},
	GoPublish: {
		url, user, password, accessToken, buildName, buildNumber, module, project, detailedSummary,
	},
	Go: {
		buildName, buildNumber, module, project,
	},
	Ping: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, InsecureTls,
	},
	RtCurl: {
		serverId,
	},
	PipConfig: {
		global, serverIdResolve, repoResolve,
	},
	PipInstall: {
		buildName, buildNumber, module, project,
	},
	ReleaseBundleCreate: {
		url, distUrl, user, password, accessToken, sshKeyPath, sshPassPhrase, serverId, specFlag, specVars, targetProps,
		rbDryRun, sign, desc, exclusions, releaseNotesPath, releaseNotesSyntax, rbPassphrase, rbRepo, InsecureTls, distTarget, rbDetailedSummary,
	},
	ReleaseBundleUpdate: {
		url, distUrl, user, password, accessToken, sshKeyPath, sshPassPhrase, serverId, specFlag, specVars, targetProps,
		rbDryRun, sign, desc, exclusions, releaseNotesPath, releaseNotesSyntax, rbPassphrase, rbRepo, InsecureTls, distTarget, rbDetailedSummary,
	},
	ReleaseBundleSign: {
		url, distUrl, user, password, accessToken, sshKeyPath, sshPassPhrase, serverId, rbPassphrase, rbRepo,
		InsecureTls, rbDetailedSummary,
	},
	ReleaseBundleDistribute: {
		url, distUrl, user, password, accessToken, sshKeyPath, sshPassPhrase, serverId, rbDryRun, distRules,
		site, city, countryCodes, sync, maxWaitMinutes, InsecureTls,
	},
	ReleaseBundleDelete: {
		url, distUrl, user, password, accessToken, sshKeyPath, sshPassPhrase, serverId, rbDryRun, distRules,
		site, city, countryCodes, sync, maxWaitMinutes, InsecureTls, deleteFromDist, deleteQuiet,
	},
	TemplateConsumer: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, vars,
	},
	RepoDelete: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, deleteQuiet,
	},
	ReplicationDelete: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, deleteQuiet,
	},
	PermissionTargetDelete: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, deleteQuiet,
	},
	AccessTokenCreate: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, clientCertPath,
		clientCertKeyPath, groups, grantAdmin, expiry, refreshable, audience,
	},
	UserCreate: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId,
		UsersGroups, Replace, Admin,
	},
	UsersCreate: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId,
		usersCreateCsv, UsersGroups, Replace,
	},
	UsersDelete: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId,
		usersDeleteCsv, deleteQuiet,
	},
	GroupCreate: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId,
		Replace,
	},
	GroupAddUsers: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId,
	},
	GroupDelete: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, deleteQuiet,
	},
	// Xray's commands
	OfflineUpdate: {
		licenseId, from, to, version, target,
	},
	XrCurl: {
		serverId,
	},
	AuditMvn: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, ExcludeTestDeps, InsecureTls, project, watches, repoPath, licenses,
	},
	AuditGradle: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, ExcludeTestDeps, UseWrapper, project, watches, repoPath, licenses,
	},
	AuditNpm: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, depType, project, watches, repoPath, licenses,
	},
	XrScan: {
		url, user, password, accessToken, sshPassPhrase, sshKeyPath, serverId, specFlag, threads, project, watches, repoPath, licenses,
	},
	// Mission Control's commands
	McConfig: {
		mcUrl, mcAccessToken, mcInteractive,
	},
	LicenseAcquire: {
		mcUrl, mcAccessToken,
	},
	LicenseDeploy: {
		mcUrl, mcAccessToken, licenseCount,
	},
	LicenseRelease: {
		mcUrl, mcAccessToken,
	},
	JpdAdd: {
		mcUrl, mcAccessToken,
	},
	JpdDelete: {
		mcUrl, mcAccessToken,
	},
}

func GetCommandFlags(cmd string) []cli.Flag {
	flagList, ok := commandFlags[cmd]
	if !ok {
		log.Error("The command \"", cmd, "\" is not found in commands flags map.")
		return nil
	}
	return buildAndSortFlags(flagList)
}

func buildAndSortFlags(keys []string) (flags []cli.Flag) {
	for _, flag := range keys {
		flags = append(flags, flagsMap[flag])
	}
	sort.Slice(flags, func(i, j int) bool { return flags[i].GetName() < flags[j].GetName() })
	return
}
