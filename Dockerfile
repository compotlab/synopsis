FROM golang:1.5.3

# Copy the local package files to the container workspace.
ADD . /go/src/github.com/compotlab/synopsis

# Change workdir in container
WORKDIR /go/src/github.com/compotlab/synopsis

# Get app dependency and install app
RUN go get
RUN go install

# Export some needed env vars
ENV HOST ""
ENV PORT 8080
ENV FILE "/data/config.json"
ENV THREAD 50
ENV OUTPUT "/data/output"

# Create volume directory
RUN mkdir /data
RUN mkdir /root/.ssh

# Add volume directory
VOLUME ["/data", "/root/.ssh"]

# Set container entrypoint
ENTRYPOINT /go/bin/synopsis

# Set port that the container will listen
EXPOSE 8080