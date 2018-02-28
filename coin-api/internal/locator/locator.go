package locator

// Node locates nodes and returns their credentials
type Node struct {
	node string
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

// GetNode returns given node ip or address
func (l *Node) GetNode() string {
	return "127.0.0.1"
}

// GetNodePort return given node port
func (l *Node) GetNodePort() int {
	return 8080
}
