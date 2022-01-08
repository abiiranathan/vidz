package constants

var formats = []string{"mp4", "webm", "3gp"}

var ignoreDirs = []string{
	"AndroidStudioProjects",
	"Android",
	"NetBeansProjects",
	"node_modules",
	"Qt",
	"VirtualBoxVMs",
	"vmime",
	"venv",
	"env",
}

// Returns the list of supported formats on the web
func GetSupportedFormats() []string {
	return formats
}

// Add a new format to the list of supported formats
func AddFormat(format string) {
	formats = append(formats, format)
}

// Remove a format from the list of supported formats
func RemoveFormat(format string) {
	// Remove format from slice
	for i, f := range formats {
		if f == format {
			formats = append(formats[:i], formats[i+1:]...)
		}
	}

}

// Returns the list of ignored directories
func GetIgnoreDirs() []string {
	return ignoreDirs
}

// Add a new ignored directory
func AddIgnoreDir(dir string) {
	ignoreDirs = append(ignoreDirs, dir)
}

// Remove an ignored directory
func RemoveIgnoreDir(dir string) {
	for i, f := range ignoreDirs {
		if f == dir {
			ignoreDirs = append(ignoreDirs[:i], ignoreDirs[i+1:]...)
		}
	}
}
