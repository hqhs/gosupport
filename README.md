Employee & Customer support dashboard
===

![](.examples/better_example.gif)

# Setup

``` sh
# complete list of commands used in gif above
$ git clone https://github.com/hqhs/gosupport.git
$ echo "POSTGRES_PASSWORD=s3cr3tpa$$sw0rd" > .env
$ make setuppostgres # docker required
$ make postgres # run database in docker container 
$ # run database migrations with alembic
$ cd scripts/ && virtualenv -p python3 .env
$ pip install -r requirements.txt
$ alembic upgrade head
$ go build .
$ ./gosupport serve \
    --dbpassword "3cr3tpa$$sw0d" \
    --tgtokens "707774881:AAHYC-8LDL20Xns1dabLzXCUXdt4YsvYLJs" # Already revoked, don't panic!
```

# Features

- Authorization sequence
- Media support (photos & files)
- Messaging part made as SPA, pure sockets & rest api. 
- Migrations with alembic & sqlalchemy for faster development
- Small dependencies 

# Project status

This was one of the first real-world applications I made in go months ago. Since
then I've learned a lot and done some refactoring before open sourcing it. It's 
not ready for use in production yet, but I'm looking forward to making it better.

# Motivation

After hours spent on github I haven't found useful open source techsupport (also
"helpdesk") dashboards with messangers/social networks integration as primary
feature. There are some like [zammad](https://github.com/zammad/zammad), but it
requires 5 docker containers to run including Elasticsearch (sic!) and chats
work only with webhooks (which in turn requires seperate static IP even for
development). Time passed and I decided to write it by myself with Go as primary 
backend language.

## Technology stack

- Postgres as primary database choice, without any ORM, pure sql requests (feels
  "Go way").
- Pure Go for backend with [Chi router](https://github.com/go-chi/chi), [gorilla
  websockets](https://github.com/gorilla/websocket), and various social
  network/messangers api wrappers. 
- Preact for frontend with some templates rendering.

Those are the tools I used in active demo development stage:
- Alembic & Sqlalchemy for migrations generation.
- Bootstrap templates.

As you can see, then choosing between some nice and large library (such as
[gorm](https://github.com/jinzhu/gorm)), and Go STD I'm almost always choose
later for guaranteed backward compatibility, extensive documentation and stable
performance. Code sometimes seems too verbose, but in return I always know what 
goes on "under the hood".

# On the road to beta release

My primary pain is frontend part, all preact components currenty live in single
main.js file, and it helps to make it useful (that is essential). Those are the features 
I plan to do before spending time on advertasing (writing on reddit/go mailing list):
- Basic admin management (deleting & blocking & creating new admins)
- Message & User search purely with postgres.
- Test coverage.
- Basic notification management (notificate on new message/then message is
  unanswered N hours/then user writes first time in a while)
- (Optional) WeChat & Viber support

# Contributing

If you feel somewhat interested, feel free to open issue/write to me on
[twitter](https://twitter.com/hqhqhs) or [telegram](http://t.me/hqhqhs). I'm
going to make pause in active development after beta and try to bring some
interest in project, and if I fail I'll just move on to some other stuff and
project will become "yet another codebase for the sake of codebase" :) 
