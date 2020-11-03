// Copyright (c) 2015 Joseph Naegele. See LICENSE file.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Path struct {
	Path     string
	User     string
	Group    string
	MinDepth int
	DirFacl  string `toml:"dir_facl"`
	FileFacl string `toml:"file_facl"`
}

type Config map[string]Path

func loadConfig(path string, config *Config) error {
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return err
	}
	return nil
}

func setOwner(p, user, group string, mindepth int) error {
	cmd := exec.Command("find", p, "-mindepth", strconv.Itoa(mindepth), "-exec", "chown", user+":"+group, "{}", "+")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func setFaclFind(p string, value string, typ string, mindepth int) error {
	tmpfile, err := ioutil.TempFile("", typ+"-*.facl")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(value)); err != nil {
		tmpfile.Close()
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	cmd := exec.Command("find", p, "-mindepth", strconv.Itoa(mindepth), "-type", typ, "-exec", "setfacl", "-b", "-M", tmpfile.Name(), "{}", "+")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func setFacl(p Path) error {
	if err := setOwner(p.Path, p.User, p.Group, p.MinDepth); err != nil {
		return err
	}

	if err := setFaclFind(p.Path, p.DirFacl, "d", p.MinDepth); err != nil {
		return err
	}
	if err := setFaclFind(p.Path, p.FileFacl, "f", p.MinDepth); err != nil {
		return err
	}
	return nil
}

func main() {
	var config Config
	if err := loadConfig("sections.toml", &config); err != nil {
		log.Fatal(err)
	}

	section := config["series"]
	if err := setFacl(section); err != nil {
		log.Fatal(err)
	}
}
