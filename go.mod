module github.com/leychan/cell

go 1.14

require (
	github.com/leychan/cell/craw v0.0.0-00010101000000-000000000000
	github.com/leychan/cell/downloader v0.0.0-00010101000000-000000000000 // indirect
)

replace (
    github.com/leychan/cell/downloader => ./downloader
    github.com/leychan/cell/craw => ./craw
)

