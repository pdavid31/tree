package internal

import (
	"io/ioutil"
	"os"
)

// Node type
type Node struct {
	Path     string
	Parent   *Node
	Info     os.FileInfo
	Children []*Node
	Config   *TreeConfig
}

// NewNode is the initial tree / root node constructor
func NewNode(path string, config *TreeConfig) (*Node, error) {
	n := &Node{Path: path, Config: config}

	info, err := os.Lstat(n.Path)
	if err != nil {
		return nil, err
	}

	n.Info = info

	if err := n.Recursive(); err != nil {
		return nil, err
	}

	return n, nil
}

// GetRoot gets the trees root
func (n Node) GetRoot() Node {
	if n.Parent == nil {
		return n
	}

	return n.Parent.GetRoot()
}

// GetConfig gets the config
// object stored in the trees root
func (n Node) GetConfig() *TreeConfig {
	return n.GetRoot().Config
}

// shouldBeIncluded checks if the node
// should be included in the output
// by applying all set filters
func (n Node) shouldBeIncluded() bool {
	c := n.GetConfig()

	// check for hidden files
	if !c.AllFiles && n.Info.Name()[0] == '.' {
		return false
	}

	// check for directories only
	if c.DirectoriesOnly && !n.Info.IsDir() {
		return false
	}

	// check for pattern matching
	// include matching files and directories that contain a matching file
	if !n.matchesGlob() {
		return false
	}

	return true
}

func (n Node) matchesGlob() bool {
	c := n.GetConfig()

	// TODO: check if the parent matches the glob

	// return true if the own name matches the glob
	if c.Pattern.Match(n.Info.Name()) {
		return true
	}

	// check if children matches the glob
	childrenMatches := false
	for _, v := range n.Children {
		if v.matchesGlob() {
			childrenMatches = true
		}
	}

	return childrenMatches
}

// Recursive recursively creates the Node tree by
// calling Recursive on the nodes if they are directories
func (n *Node) Recursive() error {
	if !n.Info.IsDir() {
		return nil
	}

	files, err := ioutil.ReadDir(n.Path)
	if err != nil {
		return err
	}

	var children []*Node
	for _, v := range files {
		p := n.Path + string(os.PathSeparator) + v.Name()
		n := &Node{Path: p, Info: v, Parent: n}
		if err := n.Recursive(); err != nil {
			return err
		}

		if !n.shouldBeIncluded() {
			continue
		}

		children = append(children, n)
	}

	n.Children = children

	return nil
}

// Equals checks if two nodes are equal
// by comparing their pathes
func (n Node) Equals(node Node) bool {
	if n.Path == node.Path {
		return true
	}

	return false
}

// getDepth gets the depth of the current Node
// by calling getDepth on the Parent and adding 1
func (n Node) getDepth() int {
	// if there is no Parent element
	// we reached the top level
	if n.Parent == nil {
		return 0
	}

	return n.Parent.getDepth() + 1
}

// isLastElement checks if the given Node
// is the last sibling
func (n Node) isLastElement() bool {
	// if there is no Parent element we are in the root
	// and therefore there are no siblings
	if n.Parent == nil {
		return true
	}

	noOfSiblings := len(n.Parent.Children)
	// return true if the current element
	// equals the last sibling
	if n.Equals(*n.Parent.Children[noOfSiblings-1]) {
		return true
	}

	return false
}

// getPrefixes recursively gets the prefixes
// of the current Node and prepends a connector
// if the the Parent element is not the last sibling
func (n Node) getPrefixes() string {
	// if the depth is 0 or 1 there is no prefix
	if n.getDepth() < 2 {
		return ""
	}

	// get the parents prefixes
	pc := n.Parent.getPrefixes()
	// add the connector if the Parent element
	// is not the last sibling
	if !n.Parent.isLastElement() {
		pc += "|"
	}

	// return prefixes and add 4 spaces
	return pc + "    "
}

// generatePrefix adds the computed prefixes
// and the actual tree symbol
func (n Node) generatePrefix() string {
	indicator := "├── "
	if n.isLastElement() {
		indicator = "└── "
	}

	return n.getPrefixes() + indicator
}

// String converts the Node to string
func (n Node) String() string {
	c := n.GetConfig()

	s := n.Info.Name() + "\n"
	if n.getDepth() == 0 || c.FullPaths {
		s = n.Path + "\n"
	}

	for _, v := range n.Children {
		if !c.DisableIndentation {
			s += v.generatePrefix()
		}
		s += v.String()
	}

	return s
}
