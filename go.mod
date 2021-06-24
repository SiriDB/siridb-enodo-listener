module github.com/SiriDB/EnodoListener

go 1.16

replace github.com/SiriDB/siridb-enodo-go-lib => ../siridb-enodo-go-lib

require (
	github.com/SiriDB/siridb-enodo-go-lib v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.2.0 // indirect
	github.com/spf13/viper v1.8.0 // indirect
	// github.com/SiriDB/siridb-enodo-go-lib v0.0.0-20210310193853-3724ed83a0cc // indirect
	github.com/transceptor-technology/go-qpack v0.0.0-20190116123619-49a14b216a45
)
