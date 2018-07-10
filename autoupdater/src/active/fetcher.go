package active

import "time"

type Fetcher interface {
	SetInterval(duration time.Duration)
	Start()
	Stop()
}

func New(name, url string) Fetcher{
	switch name {
	case "git":
		return newGit(url)
	}
	return newGit(url)
}