// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Portions Copyright 2015 Labstack.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

import (
	"strings"
)

// node - a node structure for nodes within the tree
type node struct {
	typ      ntype
	label    byte
	prefix   string
	parent   *node
	children children
	resource *resource
	pnames   pNames
}

// pNames - map of method to pnames, as different methods can have different pnames
type pNames map[string][]string

// newNode - create a new router tree node
func newNode(t ntype, pre string, p *node, c children, h *resource, pnames pNames) *node {
	n := &node{
		typ:      t,
		label:    pre[0],
		prefix:   pre,
		parent:   p,
		children: c,
		// create a resource method to handler map for this node
		resource: h,
		pnames:   pnames,
	}
	return n
}

// addChild - Add a child node to this node
func (n *node) addChild(c *node) {
	n.children = append(n.children, c)
}

// findChild - find a child node of this node
func (n *node) findChild(search string, t ntype) *node {
	for _, c := range n.children {
		if strings.HasPrefix(search, c.prefix) && c.typ == t {
			return c
		}
	}
	return nil
}

// findChildWithLabel - find a child with a matching label, label being the first byte in the prefix
func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.children {
		if c.label == l {
			return c
		}
	}
	return nil
}

// findChildWithType - find a child with a matching type
func (n *node) findChildWithType(t ntype) *node {
	for _, c := range n.children {
		if c.typ == t {
			return c
		}
	}
	return nil
}
