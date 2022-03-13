package webpack

import (
	"html/template"
	"os"
	"path/filepath"
)

const tsconfigJSON string = `
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
				"{{.Src}}/*"
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
		"{{.Dst}}"
	],
	"include": [
		"./{{.Src}}",
		"webpack.config.ts"
	]
}
`

const webpackCfgJS string = `
import path from "path";
import { Configuration, DefinePlugin } from "webpack";
import TsconfigPathsPlugin from "tsconfig-paths-webpack-plugin";

const webpackConfig = (): Configuration => ({
	entry: {
		"index": "./{{.Src}}/index.ts"
	},
	...(process.env.production || !process.env.development
		? {}
		: { devtool: "eval-source-map" }),
	resolve: {
		extensions: [".ts", ".tsx", ".js"],
		plugins: [new TsconfigPathsPlugin({ configFile: "./tsconfig.json" })],
	},
	output: {
		path: path.join(__dirname, "/{{.Dst}}"),
		filename: "[name].js",
	},
	module: {
		rules: [
			{
				test: /\.tsx?$/,
				loader: "ts-loader",
				options: {
					transpileOnly: true,
				},
				exclude: /{{.Dst}}/,
			},
			{
				test: /\.s?css$/,
				use: ["style-loader", "css-loader", "sass-loader"],
			},
		],
	},
	// devServer: {
	// 	port: 3000,
	// 	open: true,
	// 	historyApiFallback: true,
	// },
	plugins: [
		// DefinePlugin allows you to create global constants which can be configured at compile time
		new DefinePlugin({
			"process.env": process.env.production || !process.env.development,
		}),
	],
});

export default webpackConfig;
`

const (
	tsconfigFile      = "tsconfig.json"
	webpackConfigFile = "webpack.config.ts"
)

func writeTSConfig(opts *Options) error {
	t, err := template.New("text").Parse(tsconfigJSON)
	if err != nil {
		return err
	}

	w, err := os.Create(filepath.Join(opts.Root, tsconfigFile))
	if err != nil {
		return err
	}
	defer w.Close()

	if err := t.Execute(w, opts); err != nil {
		return err
	}

	return nil
}

func writeWebpackConfig(opts *Options) error {
	t, err := template.New("text").Parse(webpackCfgJS)
	if err != nil {
		return err
	}

	w, err := os.Create(filepath.Join(opts.Root, webpackConfigFile))
	if err != nil {
		return err
	}
	defer w.Close()

	if err := t.Execute(w, opts); err != nil {
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
	if err := writeWebpackConfig(opts); err != nil {
		return err
	}

	if err := writeTSConfig(opts); err != nil {
		return err
	}

	return nil
}
