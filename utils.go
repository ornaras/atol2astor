package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getExePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s является директорией", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s является директорией", p)
		}
	}
	return "", err
}

func copyFile(originPath, copyPath string) error {
	var data []byte
	var err error

	if data, err = os.ReadFile(originPath); err != nil {
		return err
	}

	if err = os.WriteFile(copyPath, data, 0644); err != nil {
		return err
	}

	return nil
}
