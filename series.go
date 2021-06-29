package main

import "regexp"

type SeriesConfig struct {
	Name       string
	IsRealtime bool
}

type GroupConfig struct {
	Name  string
	Regex *regexp.Regexp
}

// type SeriesToWatch interface {
// 	Get(key string) SeriesConfig
// 	Set(key string, series SeriesConfig)
// }

// type GroupsToWatch interface {
// 	Get(key string) GroupConfig
// 	Set(key string, group GroupConfig)
// }

type SeriesCount interface {
	Get(key string) int
	Set(key string, count int)
}
