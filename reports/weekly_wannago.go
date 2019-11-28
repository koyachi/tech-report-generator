package reports

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"

	"github.com/li-go/tech-report-generator/config"
)

func GetAppWeeklyWannago(untilDate time.Time) (int, error) {
	dataSourceName, err := config.GetDataSourceName()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer db.Close()

	db.LogMode(false)

	var count int
	query, err := config.GetCountQuery()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	query = strings.ReplaceAll(query, "%DATE_PLACEHOLDER%", untilDate.String())
	row := db.Raw(query).Row()
	if err := row.Scan(&count); err != nil {
		return 0, errors.WithStack(err)
	}
	return count, nil
}
