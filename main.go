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

	viper.SetConfigName(".config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatalf("need valid config file: %+v", err)
	}

	ctx := context.Background()
	printAppWeeklyUU(ctx, logger)
	printAppWeeklyWannago(logger)

	ctx, cancel := newChromedp(ctx, false)
	defer cancel()
	printFreeUsersRate(ctx, logger)
	printErrorRate(ctx, logger)
	printOncallCount(ctx, logger)
}

func newChromedp(ctx context.Context, headless bool) (context.Context, context.CancelFunc) {
	var opts []chromedp.ExecAllocatorOption
	for _, opt := range chromedp.DefaultExecAllocatorOptions {
		opts = append(opts, opt)
	}
	if !headless {
		opts = append(opts,
			chromedp.Flag("headless", false),
			chromedp.Flag("hide-scrollbars", false),
			chromedp.Flag("mute-audio", false),
		)
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	return ctx, func() {
		cancel()
		allocCancel()
	}
}

func printAppWeeklyUU(ctx context.Context, logger *log.Logger) {
	appWeeklyUU, err := reports.GetAppWeeklyUU(ctx)
	if err != nil {
		logger.Printf("printAppWeeklyUU: %+v", err)
		return
	}
	fmt.Printf("%s週UU: %d\n", appWeeklyUU.Date, appWeeklyUU.UU)
}

func printAppWeeklyWannago(logger *log.Logger) {
	untilDate := time.Now().Truncate(24 * time.Hour)
	wannago, err := reports.GetAppWeeklyWannago(untilDate)
	if err != nil {
		logger.Printf("printAppWeeklyWannago: %+v", err)
		return
	}
	fmt.Printf("%s週行きたい: %d\n", untilDate.Format("2006-01-02"), wannago)
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

	organization, err := config.GetFabricOrganization()
	if err != nil {
		logger.Printf("printFreeUsersRate: %+v", err)
		return
	}

	iosAppScheme, err := config.GetIOSAppScheme()
	if err != nil {
		logger.Printf("printFreeUsersRate: %+v", err)
		return
	}

	androidAppScheme, err := config.GetAndroidAppScheme()
	if err != nil {
		logger.Printf("printFreeUsersRate: %+v", err)
		return
	}

	iosCrashFreeUsers, err := f.GetIOSCrashFreeUsers(ctx, organization, iosAppScheme)
	if err != nil {
		logger.Printf("printFreeUsersRate: %+v", err)
		return
	}
	androidCrashFreeUsers, err := f.GetAndroidCrashFreeUsers(ctx, organization, androidAppScheme)
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

	transactionID, err := config.GetNewrelicTransactionID()
	if err != nil {
		logger.Printf("printErrorRate: %+v", err)
		return
	}

	appPerformance, err := n.GetAppPerformance(ctx, transactionID)
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

	organization, err := config.GetPagerdutyOrganization()
	if err != nil {
		logger.Printf("printOncallCount: %+v", err)
		return
	}

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
