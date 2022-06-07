package dsm

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func getConfig() (string, string, string, bool) {
	if !IsSet("SENHASEGURA_URL", "SENHASEGURA_CLIENT_ID", "SENHASEGURA_CLIENT_SECRET") {
		log.Fatalf("Authentication data not found or missing parameters\n")
	}

	return viper.GetString("SENHASEGURA_URL"),
		viper.GetString("SENHASEGURA_CLIENT_ID"),
		viper.GetString("SENHASEGURA_CLIENT_SECRET"),
		Verbose
}

func IsSet(name ...string) bool {
	for _, n := range name {
		if viper.GetString(n) == "" {
			v("The parameter %s is empty\n", n)
			return false
		}
	}
	return true
}

func v(format string, a ...interface{}) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

func replaceSpecials(value string) string {
	value = strings.Replace(value, "+", "-", -1)
	value = strings.Replace(value, "/", "_", -1)
	value = strings.Replace(value, "=", ",", -1)
	return value
}
