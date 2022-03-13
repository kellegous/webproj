package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/kellegous/webproj/pkg/npm"
)

type Flags struct {
	Name      string
	WithReact bool
	SrcDir    string
	DstDir    string
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.Name,
		"name",
		"",
		"the name of the project, inferred from directory by default")
	fs.BoolVar(
		&f.WithReact,
		"with-react",
		false,
		"whether to include react in deps")
	fs.StringVar(
		&f.SrcDir,
		"src-dir",
		"src",
		"the source directory to be used")
	fs.StringVar(
		&f.DstDir,
		"dst-dir",
		"dist",
		"the destination directory")
}

func inferName(root string, name string) string {
	if name != "" {
		return name
	}
	return filepath.Base(root)
}

func main() {
	var flags Flags
	flags.Register(flag.CommandLine)
	flag.Parse()

	root, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	name := inferName(root, flags.Name)
	deps := &npm.Base
	if flags.WithReact {
		deps = deps.Combine(&npm.React)
	}

	if err := npm.InitPackage(root, name, deps); err != nil {
		log.Panic(err)
	}

	// Things left to do.
	// 1. Write webpack config
	// 2. Write tsc config
	// 3. Write empty ts file in the source dir
}
