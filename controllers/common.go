package controllers

import "os"

var (
	userCacheDir, _        = os.UserCacheDir()
	DefaultFrameworkFolder = userCacheDir + "/ct.go.processor/"
)
