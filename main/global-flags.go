package main

import "flag"

type AppConfig struct {
	RootPath string
}

func AppParseFlags() (AppConfig, []string) {
	var conf AppConfig
	flag.StringVar(&conf.RootPath, "root", ".", "Root directory")
	flag.Parse()
	return conf, flag.Args()
}
