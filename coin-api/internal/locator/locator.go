package locator

// Node locates nodes and returns their credentials
type Node struct {
	Host string
	node string
	Port int32
}

// NewLocatorNode get new locator instance
func NewLocatorNode() *Node {
	return &Node{}
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
func (l *Node) GetNodePort() int32 {
	return l.Port
}
