package npm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const packageJSONTpl string = `
{
	"name": "%s",
	"version": "1.0.0",
	"description": "",
	"main": "index.js",
	"scripts": {
	  "watch-dev": "webpack --watch --mode=development --env development",
	  "build-dev": "webpack --mode=development --env development",
	  "watch-prd": "webpack --watch --mode=production --env production",
	  "build-prd": "webpack --mode=production --env production"
	},
	"keywords": [],
	"author": "Kelly Norton <kellegous@gmail.com>",
	"license": "ISC",
	"devDependencies": {
	}
  }
`

const packageJSONFile string = "package.json"

var BasePackages []string = []string{
	"@types/node",
	"@types/webpack",
	"@types/webpack-dev-server",
	"css-loader",
	"node-sass",
	"sass-loader",
	"style-loader",
	"ts-loader",
	"ts-node",
	"tsconfig-paths-webpack-plugin",
	"typescript",
	"webpack",
	"webpack-cli",
	"webpack-dev-server",
}

var ReactPackages []string = []string{
	"@types/react",
	"@types/react-dom",
}

var Base DepSet = DepSet{
	Dev: []string{
		"@types/node",
		"@types/webpack",
		"@types/webpack-dev-server",
		"css-loader",
		"node-sass",
		"sass-loader",
		"style-loader",
		"ts-loader",
		"ts-node",
		"tsconfig-paths-webpack-plugin",
		"typescript",
		"webpack",
		"webpack-cli",
		"webpack-dev-server",
	},
}

var React DepSet = DepSet{
	Dev: []string{
		"@types/react",
		"@types/react-dom",
	},
	Runtime: []string{
		"react",
		"react-dom",
	},
}

type DepSet struct {
	Runtime []string
	Dev     []string
}

func (s *DepSet) Combine(o *DepSet) *DepSet {
	return &DepSet{
		Runtime: append(s.Runtime, o.Runtime...),
		Dev:     append(s.Dev, o.Dev...),
	}
}

func writePackageJSON(root string, name string) error {
	w, err := os.Create(filepath.Join(root, packageJSONFile))
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := fmt.Fprintf(w, packageJSONTpl, name); err != nil {
		return err
	}
	return nil
}

func installDeps(root string, dev bool, deps []string) error {
	args := []string{"install"}
	if dev {
		args = append(args, "--save-dev")
	} else {
		args = append(args, "--save")
	}

	c := exec.Command("npm", append(args, deps...)...)
	c.Dir = root
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func InitPackage(root string, name string, deps *DepSet) error {
	if err := writePackageJSON(root, name); err != nil {
		return err
	}

	if len(deps.Dev) > 0 {
		if err := installDeps(root, true, deps.Dev); err != nil {
			return err
		}
	}

	if len(deps.Runtime) > 0 {
		if err := installDeps(root, false, deps.Runtime); err != nil {
			return err
		}
	}

	return nil
}
