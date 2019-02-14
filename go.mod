module github.com/hqhs/gosupport

require (
	github.com/hqhs/gosupport/cmd v0.0.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
)

replace github.com/hqhs/gosupport/cmd => ./cmd

replace github.com/hqhs/gosupport/internal/app => ./internal/app

replace github.com/hqhs/gosupport/pkg/templator => ./pkg/templator
