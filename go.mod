module github.com/leychan/cell

go 1.14

require github.com/leychan/cell/craw v0.0.0-00010101000000-000000000000

replace (
	github.com/leychan/cell/craw => ./craw
	github.com/leychan/cell/downloader => ./downloader
)
