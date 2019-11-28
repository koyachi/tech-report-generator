package config

import (
	"encoding/base64"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func getConfig(name string) (string, error) {
	v, ok := viper.Get(name).(string)
	if !ok {
		return "", errors.Errorf("Not found config: %v in %v", name, viper.AllKeys())
	}
	return v, nil
}

func GetDataSourceName() (string, error) {
	return getConfig("dataSourceName")
}

func GetCountQuery() (string, error) {
	return getConfig("countQuery")
}

func GetGoogleServiceAccountCredentials() ([]byte, error) {
	credentials, err := getConfig("googleServiceAccountCredentials")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	bs, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return bs, nil
}

func GetSpreadsheetID() (string, error) {
	return getConfig("spreadsheetID")
}

func GetSpreadsheetTabName() (string, error) {
	return getConfig("spreadsheetTabName")
}

func getAccount(name string) (string, string, error) {
	account, err := getConfig(name + "Account")
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	ss := strings.Split(account, ":")
	if len(ss) != 2 {
		return "", "", errors.Errorf("Invalid account info: %v", account)
	}
	return ss[0], ss[1], nil
}

func GetFabricAccount() (string, string, error) {
	return getAccount("fabric")
}

func GetFabricOrganization() (string, error) {
	return getConfig("fabricOrganization")
}

func GetIOSAppScheme() (string, error) {
	return getConfig("iosAppScheme")
}

func GetAndroidAppScheme() (string, error) {
	return getConfig("androidAppScheme")
}

func GetNewrelicAccount() (string, string, error) {
	return getAccount("newrelic")
}

func GetNewrelicTransactionID() (string, error) {
	return getConfig("newrelicTransactionID")
}

func GetPagerdutyAccount() (string, string, error) {
	return getAccount("pagerduty")
}

func GetPagerdutyOrganization() (string, error) {
	return getConfig("pagerdutyOrganization")
}
