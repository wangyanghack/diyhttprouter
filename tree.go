package diyhttprouter

import (
	"fmt"

	"github.com/pkg/errors"
)

type node struct {
	path     string
	indices  string
	children []*node
	handle   Handle
}

func (n *node) insertNode(fullPath string, handle Handle) error {
	if n.path == fullPath {
		// （这里检查fullPath和n.path的长度不一致，一致的话会有冲突  ？？？)
		if n.handle != nil {
			errors.New(fmt.Sprintf("handle of path %s has existed", fullPath))
		} else {
			n.handle = handle
			return nil
		}

	}
	// 1.新路由path完全属于前缀（当前节点path是abcdef,新路由path是abc），第一种情况是属于当前节点的前缀，那需要新建root节点，原root节点下移成为新root节点的子节点，
	// 第二种情况是属于当前节点的子节点node的前缀，那需要新建root的子节点，原子节点node下移成为新子节点的子节点
	// 2.当前节点path完全属于前缀（当前节点path是abc,新路由path是abcdef),
	// 3.有共同前缀（当前节点path是abcde,新路由path是abcfg),
	l := commonPrefixLength(n.path, fullPath)
	if l == 0 {
		// 没有共同前缀(没有这种可能，肯定至少会有共同前缀/)
		return nil
	} else if l == len(fullPath) {
		// fullPath完全属于前缀，可以说明当前节点是root节点；因为如果当前节点不是root节点，应该是在当前节点的父节点就可以发现fullpath完全属于前缀（这里默认fullPath和n.path的长度不一致，一致的话会有冲突）
		child := &node{
			path:     n.path[l:],
			indices:  n.indices,
			children: n.children,
			handle:   n.handle,
		}

		n.indices = string([]byte{child.path[0]})
		n.path = fullPath
		n.children = []*node{child}
		n.handle = handle

	} else if l == len(n.path) {
		//当前节点path完全属于前缀（当前节点path是abc,新路由path是abcdef)
		fullPath = fullPath[l:]
		index := getFromIndices(n.indices, fullPath[0])
		if index != -1 {
			// // 这个子节点有共同前缀
			// // 判断是否属于当前节点的子节点node的前缀，这里判断了就可以保证上面（fullPath完全属于前缀，可以说明当前节点是root节点）
			// child := n.children[index]
			// ll := commonPrefixLength(child.path, fullPath)
			// // 肯定有共同前缀
			// if ll == len(fullPath) {
			// 	child.path = child.path[ll:]
			// 	newNode := &node{
			// 		path:     fullPath[:ll],
			// 		indices:  string(child.path[0]),
			// 		children: []*node{child},
			// 	}
			// 	n.children[index] = newNode
			// } else {
			// 	return child.insertNode(fullPath, handle)
			// }

			child := n.children[index]
			return child.insertNode(fullPath, handle)
		} else {
			// 没有子节点有共同前缀
			newNode := &node{
				path:   fullPath,
				handle: handle,
			}
			n.children = append(n.children, newNode)
			n.indices += string([]byte{fullPath[0]})
		}
	} else {
		// 有共同前缀
		index := getFromIndices(n.indices, fullPath[l])
		if index != -1 {
			old := &node{
				path:     n.path[l:],
				indices:  n.indices,
				children: n.children,
				handle:   n.handle,
			}
			child := &node{
				path:   fullPath[l:],
				handle: handle,
			}
			n.path = fullPath[:l]
			n.indices = string([]byte{fullPath[l], n.path[l]})
			n.children = []*node{child, old}
		}
	}
	return nil
}

func (n *node) getValue(fullPath string) (Handle, error) {
	for {
		if n.path == fullPath {
			if n.handle == nil {
				return nil, errors.New("no handle registed")
			}
			return n.handle, nil
		}
		l := commonPrefixLength(n.path, fullPath)
		if l < len(n.path) {
			return nil, errors.New("no route registed")
		} else if l == len(n.path) {
			flag := false
			for i := 0; i < len(n.indices); i++ {
				if n.indices[i] == fullPath[l] {
					fullPath = fullPath[:l]
					n = n.children[i]
					flag = true
					break
				}
			}
			if !flag {
				return nil, errors.New("no route registed")
			}
		}
	}
}

func getFromIndices(indices string, b byte) int {
	for i := range indices {
		if indices[i] == b {
			return i
		}
	}
	return -1
}

func updatePath(path string, l int) string {
	return path[l:]
}

func commonPrefixLength(srcPath, dstPath string) int {
	l := 0
	for i := 0; i < min(len(srcPath), len(dstPath)); i++ {
		if srcPath[i] == dstPath[i] {
			l++
		} else {
			break
		}
	}
	return l
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
