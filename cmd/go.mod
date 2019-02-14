module github.com/hqhs/gosupport/cmd

require (
	github.com/go-kit/kit v0.8.0
	github.com/hqhs/gosupport/internal/app v0.0.0
	github.com/hqhs/gosupport/pkg/templator v0.0.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	gopkg.in/jinzhu/gorm.v1 v1.9.2
)

replace github.com/hqhs/gosupport/internal/app => ../internal/app

replace github.com/hqhs/gosupport/pkg/templator => ../pkg/templator
