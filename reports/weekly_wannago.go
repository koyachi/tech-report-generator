package reports

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"

	"github.com/liyy7/tech-report-generator/config"
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
	var replyScore = struct {
		ReplyScoreType   int
		ReplyScoreStatus int
	}{2, 1}
	if err := db.Table("reply_score").
		Where(replyScore).Where("create_datetime BETWEEN ? -INTERVAL 7 DAY AND ?", untilDate, untilDate).
		Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}
	return count, nil
}
