package main

import (
	"fmt"
	"os"

	"github.com/google/go-github/v54/github"
	"github.com/safedep/dry/utils"
	"github.com/safedep/vet/internal/auth"
	"github.com/safedep/vet/internal/connect"
	"github.com/safedep/vet/internal/ui"
	"github.com/safedep/vet/pkg/analyzer"
	"github.com/safedep/vet/pkg/common/logger"
	"github.com/safedep/vet/pkg/models"
	"github.com/safedep/vet/pkg/parser"
	"github.com/safedep/vet/pkg/readers"
	"github.com/safedep/vet/pkg/reporter"
	"github.com/safedep/vet/pkg/scanner"
	"github.com/spf13/cobra"
)

var (
	lockfiles                   []string
	lockfileAs                  string
	baseDirectory               string
	purlSpec                    string
	githubRepoUrls              []string
	githubOrgUrl                string
	githubOrgMaxRepositories    int
	scanExclude                 []string
	transitiveAnalysis          bool
	transitiveDepth             int
	concurrency                 int
	dumpJsonManifestDir         string
	celFilterExpression         string
	celFilterSuiteFile          string
	celFilterFailOnMatch        bool
	markdownReportPath          string
	jsonReportPath              string
	consoleReport               bool
	summaryReport               bool
	summaryReportMaxAdvice      int
	csvReportPath               string
	silentScan                  bool
	disableAuthVerifyBeforeScan bool
	syncReport                  bool
	syncReportProject           string
	syncReportStream            string
	listExperimentalParsers     bool
	failFast                    bool
)

func newScanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan and analyse package manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			startScan()
			return nil
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cmd.Flags().BoolVarP(&silentScan, "silent", "s", false,
		"Silent scan to prevent rendering UI")
	cmd.Flags().BoolVarP(&failFast, "fail-fast", "", false,
		"Fail fast when an issue is identified")
	cmd.Flags().StringVarP(&baseDirectory, "directory", "D", wd,
		"The directory to scan for lockfiles")
	cmd.Flags().StringArrayVarP(&scanExclude, "exclude", "", []string{},
		"Name patterns to ignore while scanning a directory")
	cmd.Flags().StringArrayVarP(&lockfiles, "lockfiles", "L", []string{},
		"List of lockfiles to scan")
	cmd.Flags().StringVarP(&purlSpec, "purl", "", "",
		"PURL to scan")
	cmd.Flags().StringArrayVarP(&githubRepoUrls, "github", "", []string{},
		"Github repository URL (Example: https://github.com/{org}/{repo})")
	cmd.Flags().StringVarP(&githubOrgUrl, "github-org", "", "",
		"Github organization URL (Example: https://github.com/safedep)")
	cmd.Flags().IntVarP(&githubOrgMaxRepositories, "github-org-max-repo", "", 1000,
		"Maximum number of repositories to process for the Github Org")
	cmd.Flags().StringVarP(&lockfileAs, "lockfile-as", "", "",
		"Parser to use for the lockfile (vet scan parsers to list)")
	cmd.Flags().BoolVarP(&transitiveAnalysis, "transitive", "", false,
		"Analyze transitive dependencies")
	cmd.Flags().IntVarP(&transitiveDepth, "transitive-depth", "", 2,
		"Analyze transitive dependencies till depth")
	cmd.Flags().IntVarP(&concurrency, "concurrency", "C", 5,
		"Number of concurrent analysis to run")
	cmd.Flags().StringVarP(&dumpJsonManifestDir, "json-dump-dir", "", "",
		"Dump enriched package manifests as JSON files to dir")
	cmd.Flags().StringVarP(&celFilterExpression, "filter", "", "",
		"Filter and print packages using CEL")
	cmd.Flags().StringVarP(&celFilterSuiteFile, "filter-suite", "", "",
		"Filter packages using CEL Filter Suite from file")
	cmd.Flags().BoolVarP(&celFilterFailOnMatch, "filter-fail", "", false,
		"Fail the scan if the filter match any package (security gate)")
	cmd.Flags().BoolVarP(&disableAuthVerifyBeforeScan, "no-verify-auth", "", false,
		"Do not verify auth token before starting scan")
	cmd.Flags().StringVarP(&markdownReportPath, "report-markdown", "", "",
		"Generate consolidated markdown report to file")
	cmd.Flags().BoolVarP(&consoleReport, "report-console", "", false,
		"Print a report to the console")
	cmd.Flags().BoolVarP(&summaryReport, "report-summary", "", true,
		"Print a summary report with actionable advice")
	cmd.Flags().IntVarP(&summaryReportMaxAdvice, "report-summary-max-advice", "", 5,
		"Maximum number of package risk advice to show")
	cmd.Flags().StringVarP(&csvReportPath, "report-csv", "", "",
		"Generate CSV report of filtered packages")
	cmd.Flags().StringVarP(&jsonReportPath, "report-json", "", "",
		"Generate consolidated JSON report to file (EXPERIMENTAL schema)")
	cmd.Flags().BoolVarP(&syncReport, "report-sync", "", false,
		"Enable syncing report data to cloud")
	cmd.Flags().StringVarP(&syncReportProject, "report-sync-project", "", "",
		"Project name to use in cloud")
	cmd.Flags().StringVarP(&syncReportStream, "report-sync-stream", "", "",
		"Project stream name (e.g. branch) to use in cloud")

	cmd.AddCommand(listParsersCommand())
	return cmd
}

func listParsersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parsers",
		Short: "List available lockfile parsers",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Available Lockfile Parsers\n")
			fmt.Printf("==========================\n\n")

			for idx, p := range parser.List(listExperimentalParsers) {
				fmt.Printf("[%d] %s\n", idx, p)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&listExperimentalParsers, "experimental", "", false,
		"Include experimental parsers in the list")

	return cmd
}

func startScan() {
	if !disableAuthVerifyBeforeScan {
		failOnError("auth/verify", auth.Verify(&auth.VerifyConfig{
			ControlPlaneApiUrl: auth.DefaultControlPlaneApiUrl(),
		}))
	}

	if auth.CommunityMode() {
		ui.PrintSuccess("Running in Community Mode")
	}

	failOnError("scan", internalStartScan())
}

func internalStartScan() error {
	readerList := []readers.PackageManifestReader{}
	var reader readers.PackageManifestReader
	var err error

	githubClientBuilder := func() *github.Client {
		githubClient, err := connect.GetGithubClient()
		if err != nil {
			logger.Fatalf("Failed to build Github client: %v", err)
		}

		return githubClient
	}

	// We can easily support both directory and lockfile reader. But current UX
	// contract is to support one of them at a time. Lets not break the contract
	// for now and figure out UX improvement later
	if len(lockfiles) > 0 {
		// nolint:ineffassign,staticcheck
		reader, err = readers.NewLockfileReader(lockfiles, lockfileAs)
	} else if len(githubRepoUrls) > 0 {
		githubClient := githubClientBuilder()

		// nolint:ineffassign,staticcheck
		reader, err = readers.NewGithubReader(githubClient, githubRepoUrls, lockfileAs)
	} else if len(githubOrgUrl) > 0 {
		githubClient := githubClientBuilder()

		// nolint:ineffassign,staticcheck
		reader, err = readers.NewGithubOrgReader(githubClient, &readers.GithubOrgReaderConfig{
			OrganizationURL: githubOrgUrl,
			IncludeArchived: false,
			MaxRepositories: githubOrgMaxRepositories,
		})
	} else if len(purlSpec) > 0 {
		// nolint:ineffassign,staticcheck
		reader, err = readers.NewPurlReader(purlSpec)
	} else {
		// nolint:ineffassign,staticcheck
		reader, err = readers.NewDirectoryReader(baseDirectory, scanExclude)
	}

	if err != nil {
		return err
	}

	readerList = append(readerList, reader)

	// We will always use this analyzer
	lfpAnalyzer, err := analyzer.NewLockfilePoisoningAnalyzer(analyzer.LockfilePoisoningAnalyzerConfig{
		FailFast: failFast,
	})

	if err != nil {
		return err
	}

	analyzers := []analyzer.Analyzer{lfpAnalyzer}
	if !utils.IsEmptyString(dumpJsonManifestDir) {
		task, err := analyzer.NewJsonDumperAnalyzer(dumpJsonManifestDir)
		if err != nil {
			return err
		}

		analyzers = append(analyzers, task)
	}

	if !utils.IsEmptyString(celFilterExpression) {
		task, err := analyzer.NewCelFilterAnalyzer(celFilterExpression,
			failFast || celFilterFailOnMatch)
		if err != nil {
			return err
		}

		analyzers = append(analyzers, task)
	}

	if !utils.IsEmptyString(celFilterSuiteFile) {
		task, err := analyzer.NewCelFilterSuiteAnalyzer(celFilterSuiteFile,
			failFast || celFilterFailOnMatch)
		if err != nil {
			return err
		}

		analyzers = append(analyzers, task)
	}

	reporters := []reporter.Reporter{}
	if consoleReport {
		rp, err := reporter.NewConsoleReporter()
		if err != nil {
			return err
		}

		reporters = append(reporters, rp)
	}

	if summaryReport {
		rp, err := reporter.NewSummaryReporter(reporter.SummaryReporterConfig{
			MaxAdvice: summaryReportMaxAdvice,
		})
		if err != nil {
			return err
		}

		reporters = append(reporters, rp)
	}

	if !utils.IsEmptyString(markdownReportPath) {
		rp, err := reporter.NewMarkdownReportGenerator(reporter.MarkdownReportingConfig{
			Path: markdownReportPath,
		})
		if err != nil {
			return err
		}

		reporters = append(reporters, rp)
	}

	if !utils.IsEmptyString(jsonReportPath) {
		rp, err := reporter.NewJsonReportGenerator(reporter.JsonReportingConfig{
			Path: jsonReportPath,
		})
		if err != nil {
			return err
		}

		reporters = append(reporters, rp)
	}

	if !utils.IsEmptyString(csvReportPath) {
		rp, err := reporter.NewCsvReporter(reporter.CsvReportingConfig{
			Path: csvReportPath,
		})
		if err != nil {
			return err
		}

		reporters = append(reporters, rp)
	}

	if syncReport {
		rp, err := reporter.NewSyncReporter(reporter.SyncReporterConfig{
			ProjectName: syncReportProject,
			StreamName:  syncReportStream,
		})
		if err != nil {
			return err
		}

		reporters = append(reporters, rp)
	}

	insightsEnricher, err := scanner.NewInsightBasedPackageEnricher(scanner.InsightsBasedPackageMetaEnricherConfig{
		ApiUrl:     auth.ApiUrl(),
		ApiAuthKey: auth.ApiKey(),
	})
	if err != nil {
		return err
	}

	enrichers := []scanner.PackageMetaEnricher{
		insightsEnricher,
	}

	pmScanner := scanner.NewPackageManifestScanner(scanner.Config{
		TransitiveAnalysis: transitiveAnalysis,
		TransitiveDepth:    transitiveDepth,
		ConcurrentAnalyzer: concurrency,
		ExcludePatterns:    scanExclude,
	}, readerList, enrichers, analyzers, reporters)

	// Redirect log to files to create space for UI rendering
	redirectLogToFile(logFile)

	// Trackers to handle UI
	var packageManifestTracker any
	var packageTracker any

	manifestsCount := 0
	pmScanner.WithCallbacks(scanner.ScannerCallbacks{
		OnStartEnumerateManifest: func() {
			logger.Infof("Starting to enumerate manifests")
		},
		OnEnumerateManifest: func(manifest *models.PackageManifest) {
			logger.Infof("Discovered a manifest at %s with %d packages",
				manifest.GetDisplayPath(), manifest.GetPackagesCount())

			ui.IncrementTrackerTotal(packageManifestTracker, 1)
			ui.IncrementTrackerTotal(packageTracker, int64(manifest.GetPackagesCount()))

			manifestsCount = manifestsCount + 1
			ui.SetPinnedMessageOnProgressWriter(fmt.Sprintf("Scanning %d discovered manifest(s)",
				manifestsCount))
		},
		OnStart: func() {
			if !silentScan {
				ui.StartProgressWriter()
			}

			packageManifestTracker = ui.TrackProgress("Scanning manifests", 0)
			packageTracker = ui.TrackProgress("Scanning packages", 0)
		},
		OnAddTransitivePackage: func(pkg *models.Package) {
			ui.IncrementTrackerTotal(packageTracker, 1)
		},
		OnDoneManifest: func(manifest *models.PackageManifest) {
			ui.IncrementProgress(packageManifestTracker, 1)
		},
		OnDonePackage: func(pkg *models.Package) {
			ui.IncrementProgress(packageTracker, 1)
		},
		BeforeFinish: func() {
			ui.MarkTrackerAsDone(packageManifestTracker)
			ui.MarkTrackerAsDone(packageTracker)
			ui.StopProgressWriter()
		},
	})

	return pmScanner.Start()
}
