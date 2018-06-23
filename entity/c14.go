package entity

import "time"

type CreateArchive struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Protocols     []string  `json:"protocols"`
	Platforms     []string  `json:"platforms"`
}

type Archive struct {
	Ref string `json:"$ref"`
	Uuid_ref string `json:"uuid_ref"`
	Name string `json:"name"`
	Status string `json:"status"`
}

type UnArchive struct {
	Platform string `json:"location_id"`
	Protocols []string 	`json:"protocols"`
}

type CreateSafe struct {
	Name string `json:"name"`
}

func NewCredentials(credentials map[string]interface{}) Credentials {
	return Credentials{
		Protocol: credentials["protocol"].(string),
		Uri:      credentials["uri"].(string),
		Login:    credentials["login"].(string),
		Password: credentials["password"].(string),
	}
}

type Credentials struct {
	Protocol      string    `json:"protocol"`
	Uri           string    `json:"uri"`
	Login         string    `json:"login"`
	Password      string    `json:"password"`
}

type Bucket struct {
	Ref           string    `json:"$ref"`
	Uuid_ref      string    `json:"uuid_ref"`
	Status        string    `json:"status"`
	Archival_date time.Time `json:"archival_date"`
	Credentials   []interface{} `json:"credentials"`
}
