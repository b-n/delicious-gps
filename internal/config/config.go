package config

import "flag"

type Options struct {
	ShowDebug bool
	Database  string
}

func Init(args []string) Options {
	opts := Options{}

	flag.StringVar(&opts.Database, "database", "data.db", "the name of the database file to output to")
	flag.BoolVar(&opts.ShowDebug, "debug", false, "if true, output debug logging")

	flag.Parse()

	return opts
}
