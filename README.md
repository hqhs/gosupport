Employee & Customer support dashboard
===

![](.examples/better_example.gif)

# Setup

``` sh
# complete list of commands used in gif above
$ git clone https://github.com/hqhs/gosupport.git
$ echo "POSTGRES_PASSWORD=s3cr3tpa$$sw0rd" > .env
$ make setuppostgres # docker required
$ make postgres
$ go build .
$ ./gosupport serve \
    --dbpassword "3cr3tpa$$sw0d" \
    --tgtokens "707774881:AAHYC-8LDL20Xns1dabLzXCUXdt4YsvYLJs" # Already turned
    # off, don't panic!
```

# Features

# Project status

# Motivation

## Technology stack

# Project documentation

# On the road to beta release

# Contributing
