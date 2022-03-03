FROM golang:1.17-alpine

ARG USER=worker
ENV HOME /home/$USER

# install sudo as root
RUN apk add --update sudo

# add new user
RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

USER $USER
WORKDIR $HOME

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY vendor ./vendor
COPY docs ./docs

RUN go mod download

RUN go build -o ./job-manager

CMD [ "./job-manager" ]
