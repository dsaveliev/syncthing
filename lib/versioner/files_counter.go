// Copyright (C) 2017 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package versioner

import "path/filepath"

// filesCounter is intended to track empty directories
type filesCounter struct {
	filesPerDir map[string]int
	bypassOrder []string
}

func newFilesCounter() *filesCounter {
	return &filesCounter{
		filesPerDir: make(map[string]int),
		bypassOrder: []string{},
	}
}

func (c *filesCounter) addDir(path string) {
	c.bypassOrder = append(c.bypassOrder, path)
	c.filesPerDir[path] = 0
	if path != "." {
		dir := filepath.Dir(path)
		c.filesPerDir[dir]++
	}
}

func (c *filesCounter) addFile(path string) {
	dir := filepath.Dir(path)
	c.filesPerDir[dir]++
}

func (c *filesCounter) emptyDirs() []string {
	empty := []string{}

	// We need to keep the bypass order of the Filesystem.Walk method
	// to start from the nested dirs up to the root.
	for i := len(c.bypassOrder) - 1; i >= 0; i-- {
		path := c.bypassOrder[i]
		if c.filesPerDir[path] > 0 {
			continue
		}
		empty = append(empty, path)
		dir := filepath.Dir(path)
		c.filesPerDir[dir]--
	}

	return empty
}
