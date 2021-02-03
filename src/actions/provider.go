package actions

type ActionProvider struct {
	FilenameKey       string
	FormFilesKey      string
	FormFileKey       string
	RootFileDirectory string
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
