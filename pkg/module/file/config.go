package file

// module names for each sub-module
const (
	UploadModuleName   = "upload"
	DownloadModuleName = "download"
)

type operation int

const (
	download operation = iota
	upload
)

var operations map[string]operation = map[string]operation{
	UploadModuleName:   upload,
	DownloadModuleName: download,
}

// struct members must be public for serialization to work properly
type config struct {
	Name      string
	Operation operation
	Local     string
	Remote    string
	Content   []byte
}
