package api

type Instance struct {
	Id           string        `json:"id"`
	Title        string        `json:"title"`
	Image        *Image        `json:"image"`
	InstanceType *InstanceType `json:"instance_type"`
	Disk         int64         `json:"disk"`
}

type Image struct {
	Id string `json:"id"`
}

type InstanceType struct {
	Id string `json:"id"`
}
