package reports

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"github.com/liyy7/tech-report-generator/config"
)

type WeeklyUU struct {
	Date string
	UU   int
}

func newSheetsService() (*sheets.Service, error) {
	credentials, err := config.GetGoogleServiceAccountCredentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jwtConfig, err := google.JWTConfigFromJSON(credentials, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	client := jwtConfig.Client(context.Background())

	return sheets.New(client)
}

func GetAppWeeklyUU() (*WeeklyUU, error) {
	srv, err := newSheetsService()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// https://docs.google.com/spreadsheets/d/1CC74lKi09Eo3cxOX2lHKxbSR1wOoOzaiBCADiwIZ7bs/edit#gid=1233256154
	const spreadsheetID = "1CC74lKi09Eo3cxOX2lHKxbSR1wOoOzaiBCADiwIZ7bs"
	const readRange = "APP_1WEEK!A2:B"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
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
