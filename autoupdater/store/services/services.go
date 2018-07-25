package services

var Store Storer

type Storer interface {
	Get(string) *Service
	Store(*Service)
}

type Services []*Service

type Service struct {
	CurrentReleaseString string `json:"current_version"`
	currentRelease date.Time
}

func InitStorer(store string) {
	switch store {
	case "json":
		Store = newJsonStore("")
	}
}
