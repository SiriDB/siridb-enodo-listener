package main

type SeriesConfig struct {
	Name       string
	IsRealtime bool
}

type SeriesToWatch interface {
	Get(key string) SeriesConfig
	Set(key string, series SeriesConfig)
}

type SeriesCount interface {
	Get(key string) int
	Set(key string, count int)
}
