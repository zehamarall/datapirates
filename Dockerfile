### STAGE 1: Build ###

# The builder node
FROM golang:latest as builder

# create working directory
WORKDIR /go/src/datapirates

# copy the content 
COPY . .

# install dependencies
RUN go get ./...

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o datapirateschallenge .


### STAGE 2: Setup ###

# The runner node
FROM alpine:latest as runner 

# setup env
RUN apk --no-cache add ca-certificates
WORKDIR /root/
RUN mkdir data

# copy the binary from previous stage
COPY --from=builder /go/src/datapirates/datapirateschallenge .

RUN ls -ltr

# execute
CMD ["./datapirateschallenge"]
