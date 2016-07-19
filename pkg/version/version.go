package version

var VERSION string
var GITCOMMIT string

func init() {
	// Version should be updated by hand at each release
	VERSION = "2.4.1+dev"
	// GitCommit will be overwritten automatically by the build system
	GITCOMMIT = "HEAD"
}
