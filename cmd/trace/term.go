//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

var (
	reSize = regexp.MustCompilePOSIX(`^([[:digit:]]+)[[:space:]]+([[:digit:]]+)$`)
)

// Size returns the screen size in rows and columns.
func Size() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	m := reSize.FindStringSubmatch(string(output))
	if m == nil {
		return 0, 0, fmt.Errorf("could not get tty size")
	}
	h, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, 0, err
	}
	w, err := strconv.Atoi(m[2])
	if err != nil {
		return 0, 0, err
	}
	return h, w, nil
}
