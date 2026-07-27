// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/synapseq-foundation/synapseq/v4/sbg"
)

const exampleContent = "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"

func ExampleLoadContent() {
	loaded, err := sbg.LoadContent(exampleContent)
	if err != nil {
		panic(err)
	}

	fmt.Println(loaded.Duration())
	// Output: 1m0s
}

func ExampleLoadFile() {
	dir, err := os.MkdirTemp("", "synapseq-sbg-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "session.sbg")
	if err := os.WriteFile(path, []byte(exampleContent), 0o600); err != nil {
		panic(err)
	}
	loaded, err := sbg.LoadFile(path)
	if err != nil {
		panic(err)
	}

	fmt.Println(loaded.Duration())
	// Output: 1m0s
}
