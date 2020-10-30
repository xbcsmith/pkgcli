// Copyright © 2020 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package compress

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"

	"github.com/xbcsmith/pkgcli/lpak/common"
)

// Compress func takes no input and returns src string, excludes []string, writers ...io.Writer error
func Compress(src string, excludes []string, writers ...io.Writer) error {
	fmt.Printf("Creating tar xz : %s\n", src)

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	xzw, err := xz.NewWriter(mw)
	if err != nil {
		log.Fatalf("xz.NewWriter error %s", err)
	}
	defer xzw.Close()

	tw := tar.NewWriter(xzw)
	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if common.StringInSlice(fi.Name(), excludes) {
			fmt.Printf("skipping file : %s\n", fi.Name())
			return nil
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.ReplaceAll(file, src, ""), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})
}
