package models

// Info ...
type Info struct {
	Name         string                        `json:"name"`
	Build        InfoBuild                     `json:"build"`
	Dependencies map[string]InfoDependencyItem `json:"dependencies"`
}

// InfoBuild ...
type InfoBuild struct {
	Version   string `jsin:"version"`
	Date      string `json:"date"`
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	GoVersion string `json:"go_version"`
}

// InfoDependencyItem ...
type InfoDependencyItem struct {
	Status   int     `json:"status"`
	Latency  float64 `json:"latency"`
	Optional bool    `json:"optional"`
}
