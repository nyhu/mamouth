FROM golang:1.10 AS mamouth-base

WORKDIR /go/src/github.com/nyhu/mamouth

RUN \
    go get \
        github.com/golang/dep/cmd/dep \
        github.com/golang/lint/golint \
        github.com/vektra/mockery/.../ \
        github.com/jstemmer/go-junit-report \
        honnef.co/go/tools/cmd/megacheck \
        github.com/client9/misspell/cmd/misspell \
        github.com/canthefason/go-watcher/cmd/watcher \
        github.com/devimteam/microgen/cmd/microgen

#########################################################

FROM mamouth-base as mamouth-build


#WORKDIR /root
#RUN git clone https://github.com/edenhill/librdkafka.git
#WORKDIR /root/librdkafka
#RUN ./configure --prefix /usr
#RUN make
#RUN make install


WORKDIR /go/src/github.com/nyhu/mamouth

COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure -vendor-only

COPY . .

RUN \
    mkdir -p build && \
    go build -tags=jsoniter -ldflags "-X main.VERSION=$VERSION" -o build/mamouth github.com/nyhu/mamouth/cmd/mamouth/.

ENTRYPOINT ["./build/mamouth"]
