package entities

type ScanResult struct {
	Stdout string `json:"stdout" bson:"stdout"`
	Stderr string `json:"stderr" bson:"stderr"`
	Error  error `json:"error" bson:"error"`
}

type ProcessScanResult struct {
	MalwareFound bool `json:"malware_found" bson:"malware_found"`
	InfectedFilesCount int `json:"infected_files_count" bson:"infected_files_count"`
}