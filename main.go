// Copyright (c) 2015 Joseph Naegele. See LICENSE file.
package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/BurntSushi/toml"
)

type Path struct {
	Path     string
	DirFacl  string `toml:"dir_facl"`
	FileFacl string `toml:"file_facl"`
}

type Config struct {
	Title string
	Paths map[string]Path
}

func main() {
	var config map[string]Path
	if _, err := toml.DecodeFile("example.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	section := config["series"]

	cmd := exec.Command(setFacl, "-b", "-M", "-", section.Path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, section.DirFacl)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)
}
