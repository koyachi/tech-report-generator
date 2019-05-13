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

func GetNewrelicAccount() (string, string, error) {
	return getAccount("newrelic")
}

func GetPagerdutyAccount() (string, string, error) {
	return getAccount("pagerduty")
}

func GetOrganization() (string, error) {
	return getConfig("organization")
}

func GetIOSApp() (string, error) {
	return getConfig("iosApp")
}

func GetAndroidApp() (string, error) {
	return getConfig("androidApp")
}
