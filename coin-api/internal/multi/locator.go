package multi

// Node locates nodes and returns their credentials
type Node struct {
	Host string
	Port int
	node string
}

// NewLocatorNode get new locator instance
func NewLocatorNode(host string, port int) *Node {
	return &Node{
		Host: host,
		Port: port,
	}
}

// SetNode set node name to locate it
func (l *Node) SetNode(node string) *Node {
	l.node = node
	return l
}

// GetNodeHost returns given node ip or address
func (l *Node) GetNodeHost() string {
	return l.Host
}

// GetNodePort return given node port
func (l *Node) GetNodePort() int {
	return l.Port
}
