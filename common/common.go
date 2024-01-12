package common

import "path"

func GetBasePath(str string) string {
	dir := path.Dir(str)
	return dir
}
