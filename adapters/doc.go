// Package adapters installs all adapters from its subpackages into the objconv
// package.
//
// This package exposes no functions or types and is solely useful for the side
// effect of setting up extra adapters on the objconv package on initialization.
package adapters

import (
	_ "github.com/dolab/objconv/adapters/net"
	_ "github.com/dolab/objconv/adapters/net/mail"
	_ "github.com/dolab/objconv/adapters/net/url"
)
