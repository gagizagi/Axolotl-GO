FROM golang:1.8-jessie

# install glide
RUN go get github.com/Masterminds/glide

# create a working directory
WORKDIR /go/src/app

# add glide.yaml and glide.lock
ADD glide.yaml glide.yaml
ADD glide.lock glide.lock

# install packages
RUN glide install

# add source files
ADD src src

# build app binary
RUN go install app/src

CMD [ "src" ]