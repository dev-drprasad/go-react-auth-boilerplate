# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

WORKDIR /go/src/repoboost
ADD go.mod .
RUN go mod download
# Copy the local package files to the container's workspace.
COPY cmd cmd
COPY internal internal

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install cmd/main.go




# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/main

# Document that the service listens on port 8080.
EXPOSE 8000
