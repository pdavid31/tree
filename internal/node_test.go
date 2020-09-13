package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode_Equals(t *testing.T) {
	tests := []struct {
		name     string
		arg      []Node
		expected bool
	}{
		{
			name: "equal nodes",
			arg: []Node{
				{Path: ".idea"},
				{Path: ".idea"},
			},
			expected: true,
		},
		{
			name: "not equal nodes",
			arg: []Node{
				{Path: ".idea"},
				{Path: ".gitignore"},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := test.arg[0].Equals(test.arg[1])
			assert.Equal(t, test.expected, output)
		})
	}
}

func TestNode_GetConfig(t *testing.T) {
	tests := []struct {
		name            string
		arg             *Node
		prepareFunction func(*Node)
		accessFunction  func(*Node) *TreeConfig
		expected        *TreeConfig
	}{
		{
			name: "simple",
			arg:  &Node{Config: &TreeConfig{DirectoriesOnly: true}},
			prepareFunction: func(node *Node) {
				return
			},
			accessFunction: func(node *Node) *TreeConfig {
				return node.GetConfig()
			},
			expected: &TreeConfig{DirectoriesOnly: true},
		},
		{
			name: "nested",
			arg: &Node{Config: &TreeConfig{DirectoriesOnly: true}, Children: []*Node{
				{Path: ".idea"},
			}},
			prepareFunction: func(node *Node) {
				node.Children[0].Parent = node
			},
			accessFunction: func(node *Node) *TreeConfig {
				return node.Children[0].GetConfig()
			},
			expected: &TreeConfig{DirectoriesOnly: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// add the parent node to the child
			// to enable use of get root function
			test.prepareFunction(test.arg)

			output := test.accessFunction(test.arg)
			assert.Equal(t, test.expected, output)
		})
	}
}

func TestNode_GetRoot(t *testing.T) {
	tests := []struct {
		name            string
		arg             *Node
		prepareFunction func(*Node)
		accessFunction  func(*Node) Node
		expected        Node
	}{
		{
			name: "simple",
			arg:  &Node{},
			prepareFunction: func(node *Node) {
				return
			},
			accessFunction: func(node *Node) Node {
				return node.GetRoot()
			},
			expected: Node{},
		},
		{
			name: "nested",
			arg: &Node{Children: []*Node{
				{Path: ".idea"},
			}},
			prepareFunction: func(node *Node) {
				node.Children[0].Parent = node
			},
			accessFunction: func(node *Node) Node {
				return node.Children[0].GetRoot()
			},
			expected: Node{Children: []*Node{
				{Path: ".idea"},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// add the parent node to the child
			test.prepareFunction(test.arg)
			// add the parent node to the child
			// of the expected object as well
			test.prepareFunction(&test.expected)

			output := test.accessFunction(test.arg)
			assert.Equal(t, test.expected, output)
		})
	}
}

func TestNode_String(t *testing.T) {
	tests := []struct {
		name            string
		arg             *Node
		prepareFunction func(*Node)
		expected        string
	}{
		{
			name: "simple",
			arg:  &Node{Config: &TreeConfig{}, Path: ".", Info: fakeFile{name: "."}},
			prepareFunction: func(node *Node) {
				return
			},
			expected: ".\n",
		},
		{
			name: "nested",
			arg: &Node{Config: &TreeConfig{}, Path: ".", Info: fakeFile{name: "."}, Children: []*Node{
				{Path: "./.idea", Info: fakeFile{name: ".idea"}, Children: []*Node{
					{Path: "./.idea/run.xml", Info: fakeFile{name: "run.xml"}},
				}},
				{Path: "./.gitignore", Info: fakeFile{name: ".gitignore"}},
			}},
			prepareFunction: func(node *Node) {
				node.Children[0].Parent = node
				node.Children[1].Parent = node
				node.Children[0].Children[0].Parent = node.Children[0]
			},
			expected: ".\n├── .idea\n|\t└── run.xml\n└── .gitignore\n",
		},
		{
			name: "without indentation",
			arg: &Node{Config: &TreeConfig{DisableIndentation: true}, Path: ".", Info: fakeFile{name: "."}, Children: []*Node{
				{Path: "./.idea", Info: fakeFile{name: ".idea"}, Children: []*Node{
					{Path: "./.idea/run.xml", Info: fakeFile{name: "run.xml"}},
				}},
				{Path: "./.gitignore", Info: fakeFile{name: ".gitignore"}},
			}},
			prepareFunction: func(node *Node) {
				node.Children[0].Parent = node
				node.Children[1].Parent = node
				node.Children[0].Children[0].Parent = node.Children[0]
			},
			expected: ".\n.idea\nrun.xml\n.gitignore\n",
		},
		{
			name: "without indentation, with full paths",
			arg: &Node{Config: &TreeConfig{DisableIndentation: true, FullPaths: true}, Path: ".", Info: fakeFile{name: "."}, Children: []*Node{
				{Path: "./.idea", Info: fakeFile{name: ".idea"}, Children: []*Node{
					{Path: "./.idea/run.xml", Info: fakeFile{name: "run.xml"}},
				}},
				{Path: "./.gitignore", Info: fakeFile{name: ".gitignore"}},
			}},
			prepareFunction: func(node *Node) {
				node.Children[0].Parent = node
				node.Children[1].Parent = node
				node.Children[0].Children[0].Parent = node.Children[0]
			},
			expected: ".\n./.idea\n./.idea/run.xml\n./.gitignore\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// add the parent node to the child
			test.prepareFunction(test.arg)

			output := test.arg.String()
			assert.Equal(t, test.expected, output)
		})
	}
}

func TestNode_isLastElement(t *testing.T) {
	tests := []struct {
		name            string
		arg             *Node
		prepareFunction func(*Node)
		accessFunction  func(*Node) bool
		expected        bool
	}{
		{
			name: "root",
			arg:  &Node{},
			prepareFunction: func(node *Node) {
				return
			},
			accessFunction: func(node *Node) bool {
				return node.isLastElement()
			},
			expected: true,
		},
		{
			name: "nested - false",
			arg: &Node{Children: []*Node{
				{Path: ".idea"},
				{Path: ".gitignore"},
			}},
			prepareFunction: func(node *Node) {
				node.Children[0].Parent = node
				node.Children[1].Parent = node
			},
			accessFunction: func(node *Node) bool {
				return node.Children[0].isLastElement()
			},
			expected: false,
		},
		{
			name: "nested - true",
			arg: &Node{Children: []*Node{
				{Path: ".idea"},
				{Path: ".gitignore"},
			}},
			prepareFunction: func(node *Node) {
				node.Children[0].Parent = node
				node.Children[1].Parent = node
			},
			accessFunction: func(node *Node) bool {
				return node.Children[1].isLastElement()
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// add the parent node to the child
			test.prepareFunction(test.arg)

			output := test.accessFunction(test.arg)
			assert.Equal(t, test.expected, output)
		})
	}
}
