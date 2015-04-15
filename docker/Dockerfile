FROM google/golang:1.4

MAINTAINER Bryan-Kirk Reinhardt <composer22@users.noreply.github.com>

RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/composer22/clidemo
COPY Dockerfile.final /gopath/bin/Dockerfile

CMD docker build -t composer22/clidemo /gopath/bin
