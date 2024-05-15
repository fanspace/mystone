package viewmodel

type Apitem struct {
}

type ApiTree struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	GroupName string `json:"groupName"`
	Domain    string `json:"domain"`
	Pid       int64  `json:"pid"`
	NameCn    string `json:"nameCn"`
	Url       string `json:"url"`
}
