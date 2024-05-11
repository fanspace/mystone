package viewmodel

type Apis struct {
	Apis      []*Apitem `json:"apis"`
	IsLeaf    bool      `json:"isLeaf"`
	GroupName string    `json:"groupName"`
}

type Apitem struct {
}
