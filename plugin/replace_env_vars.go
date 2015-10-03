package plugin

import (
	"os"
	"strings"
)

func ReplaceEnvVars(args []string) (res []string) {
	for _, v := range args {

		if strings.HasPrefix(v, Prefix) && strings.HasSuffix(v, Suffix) {
			varname := strings.TrimSuffix(strings.TrimPrefix(v, Prefix), Suffix)
			v = os.Getenv(varname)
		}
		res = append(res, v)
	}
	return
}
