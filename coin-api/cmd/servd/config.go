package servd

type Config struct {
	Server  Server
	Bitcoin Bitcoin
	SkyCoin SkycoinConfig
}

type Server struct {
	ListenStr    string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type Bitcoin struct {
	NodeAddress   string
	User          string
	Password      string
	CertFile      string
	TLS           bool
	BlockExplorer string
}

// TODO(stgleb): Implement this
type SkycoinConfig struct {
	Host string
}
