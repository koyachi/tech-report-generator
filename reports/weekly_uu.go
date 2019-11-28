package reports

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/li-go/tech-report-generator/config"
)

type WeeklyUU struct {
	Date string
	UU   int
}

func newSheetsService(ctx context.Context) (*sheets.Service, error) {
	jsonCredentials, err := config.GetGoogleServiceAccountCredentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	credentials, err := google.CredentialsFromJSON(ctx, jsonCredentials, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sheets.NewService(ctx, option.WithCredentials(credentials))
}

func GetAppWeeklyUU(ctx context.Context) (*WeeklyUU, error) {
	srv, err := newSheetsService(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	spreadsheetID, err := config.GetSpreadsheetID()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	spreadsheetTabName, err := config.GetSpreadsheetTabName()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, spreadsheetTabName).Do()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(resp.Values) == 0 {
		return nil, errors.New("no data found")
	}

	lastRow := resp.Values[len(resp.Values)-1]

	date, ok := lastRow[0].(string)
	if !ok {
		return nil, errors.New("no date found")
	}
	uu, err := strconv.Atoi(strings.ReplaceAll(lastRow[1].(string), ",", ""))
	if err != nil {
		return nil, errors.New("no UU found")
	}
	return &WeeklyUU{Date: date, UU: uu}, nil
}
