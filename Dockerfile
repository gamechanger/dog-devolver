FROM docker.gamechanger.io/go1.3.1
MAINTAINER Travis Thieman <travis@gc.io>

ADD . /gc/dog-devolver
WORKDIR /gc/dog-devolver
# Drone does this GOPATH magic for us since it's a Go tool, apparently
RUN mkdir -p $GOPATH/src/github.com/gamechanger
RUN ln -s /gc/dog-devolver $GOPATH/src/github.com/gamechanger/dog-devolver

RUN godep restore
RUN go build -v

RUN mkdir -p /var/log/dog-devolver

CMD ./docker-init
