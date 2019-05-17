package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/li-go/chromedp-samples/samples/fabric"
	"github.com/li-go/chromedp-samples/samples/newrelic"
	"github.com/li-go/chromedp-samples/samples/pagerduty"
	"github.com/spf13/viper"

	"github.com/li-go/tech-report-generator/config"
	"github.com/li-go/tech-report-generator/reports"
)

func main() {
	logger := log.New(os.Stdout, "[TRGen] ", log.LstdFlags|log.Lshortfile)

	viper.SetConfigName(".config") // name of config file (without extension)
	viper.AddConfigPath(".")       // optionally look for config in the working directory
	err := viper.ReadInConfig()    // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		logger.Printf("need valid config file: %+v", err)
		return
	}

	chromedp.DefaultExecAllocatorOptions = []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	}
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(logger.Printf))
	defer cancel()

	printAppWeeklyUU(logger)
	printAppWeeklyWannago(logger)
	printFreeUsersRate(ctx, logger)
	printErrorRate(ctx, logger)
	printOncallCount(ctx, logger)
}

func printAppWeeklyUU(logger *log.Logger) {
	appWeeklyUU, err := reports.GetAppWeeklyUU()
	if err != nil {
		logger.Printf("printAppWeeklyUU: %+v", err)
		return
	}
	fmt.Printf("%s週UU: %d\n", appWeeklyUU.Date, appWeeklyUU.UU)
}

func printAppWeeklyWannago(logger *log.Logger) {
	year, month, day := time.Now().Date()
	untilDate := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	wannago, err := reports.GetAppWeeklyWannago(untilDate)
	if err != nil {
		logger.Printf("printAppWeeklyWannago: %+v", err)
		return
	}
	fmt.Printf("%d-%d-%d週行きたい: %d\n", year, month, day, wannago)
}

func printFreeUsersRate(ctx context.Context, logger *log.Logger) {
	f := fabric.New(logger)
	user, pass, err := config.GetFabricAccount()
	if err != nil {
		logger.Printf("printFreeUsersRate: %+v", err)
		return
	}
	if err := f.Login(ctx, user, pass); err != nil {
		logger.Printf("printFreeUsersRate: %+v", err)
		return
	}

	organization, _ := config.GetOrganization()
	iosApp, _ := config.GetIOSApp()
	androidApp, _ := config.GetAndroidApp()

	iosCrashFreeUsers, err := f.GetIOSCrashFreeUsers(ctx, organization, iosApp)
	if err != nil {
		logger.Printf("printFreeUsersRate: %+v", err)
		return
	}
	androidCrashFreeUsers, err := f.GetAndroidCrashFreeUsers(ctx, organization, androidApp)
	if err != nil {
		logger.Printf("printFreeUsersRate: %+v", err)
		return
	}
	fmt.Printf("FreeUsersRate: %.2f%%\t%.2f%%\n", iosCrashFreeUsers, androidCrashFreeUsers)
}

func printErrorRate(ctx context.Context, logger *log.Logger) {
	n := newrelic.New(logger)

	user, pass, err := config.GetNewrelicAccount()
	if err != nil {
		logger.Printf("printErrorRate: %+v", err)
		return
	}

	if err := n.Login(ctx, user, pass); err != nil {
		logger.Printf("printErrorRate: %+v", err)
		return
	}

	errorRate, err := n.GetErrorRate(ctx)
	if err != nil {
		logger.Printf("printErrorRate: %+v", err)
		return
	}
	fmt.Printf("ErrorRate: %f%%\n", errorRate)

	appPerformance, err := n.GetAppPerformance(ctx, "5b225765625472616e73616374696f6e2f526571756573744174747269627574652f76352e342f6170702f6d652f66656564202847455429202d206a6170616e222c22225d")
	if err != nil {
		logger.Printf("printErrorRate: %+v", err)
		return
	}
	fmt.Printf("App response: %d ms\n", appPerformance.AppResponse)
	fmt.Printf("App histogram: %d ms\n", appPerformance.AppHistogram)
	fmt.Printf("App percentile: %d ms\n", appPerformance.AppPercentile)
}

func printOncallCount(ctx context.Context, logger *log.Logger) {
	p := pagerduty.New(logger)

	user, pass, err := config.GetPagerdutyAccount()
	if err != nil {
		logger.Printf("printOncallCount: %+v", err)
		return
	}

	organization, _ := config.GetOrganization()

	if err := p.Login(ctx, organization, user, pass); err != nil {
		logger.Printf("printOncallCount: %+v", err)
		return
	}

	count, err := p.GetOncallCount(ctx, organization)
	if err != nil {
		logger.Printf("printOncallCount: %+v", err)
		return
	}
	fmt.Printf("OncallCount: %d\n", count)
}
