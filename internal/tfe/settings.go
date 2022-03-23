package tfe

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	token  string
	domain string
)

func Configure() error {
	pflag.String("token", "", "TFE token to use in API requests")
	pflag.String("domain", "terraform.management.form3.tech", "domain name of the TFE instance to connect to")
	pflag.Parse()

	if err := viper.BindEnv("TFE_TOKEN", "TFE_DOMAIN"); err != nil {
		return err
	}

	var err error
	pflag.CommandLine.VisitAll(func(flag *pflag.Flag) {
		if err == nil {
			envName := fmt.Sprintf("TFE_%s", strings.ToUpper(strings.ReplaceAll(flag.Name, "-", "_")))
			err = viper.BindPFlag(envName, flag)
		}
	})
	if err != nil {
		return err
	}

	token = viper.GetString("TFE_TOKEN")
	domain = viper.GetString("TFE_DOMAIN")

	return nil
}
