package webpack

import (
	"fmt"
	"os"
	"path/filepath"
)

const webpackCfgJS string = `
{
	"compilerOptions": {
		"module": "commonjs",
		"moduleResolution": "node",
		"target": "ESNext",
		"removeComments": true,
		"allowSyntheticDefaultImports": true,
		"jsx": "react",
		"allowJs": true,
		"baseUrl": "ui",
		"esModuleInterop": true,
		"resolveJsonModule": true,
		"downlevelIteration": true,
		"paths": {
			"*": [
				"%s/*"
			]
		},
		"lib": [
			"esnext",
			"dom",
			"dom.iterable"
		],
		"sourceMap": true,
		"noImplicitAny": false,
	},
	"exclude": [
		"node_modules",
		"%s"
	],
	"include": [
		"./%s",
		"webpack.config.ts"
	]
}
`

const webpackCfgFile = "webpack.config.ts"

func writeWebpackCfg(opts *Options) error {
	w, err := os.Create(filepath.Join(opts.Root, webpackCfgFile))
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := fmt.Fprintf(w, webpackCfgJS, opts.Src, opts.Dst, opts.Src); err != nil {
		return err
	}

	return nil
}

type Options struct {
	Name string
	Root string
	Src  string
	Dst  string
}

func CreateConfig(opts *Options) error {
	if err := writeWebpackCfg(opts); err != nil {
		return err
	}
	return nil
}
