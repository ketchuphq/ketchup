package version

var ketchupVersion = ""

func Set(version string) {
	ketchupVersion = version
}

func Get() string {
	return ketchupVersion
}
