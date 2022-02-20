FROM golang:1.18-rc-bullseye AS builder
WORKDIR /app
COPY . .
RUN go build -o keys github.com/toshokan/sshaft/cmd/keys
RUN go build -o login github.com/toshokan/sshaft/cmd/login

FROM debian:bullseye-slim

RUN useradd -m sshaft -p '*'
RUN useradd -m mfa -p '*'
RUN apt-get update && apt-get install --no-install-recommends -y ssh
WORKDIR /sshaft
COPY --from=builder /app/keys keys
COPY --from=builder /app/login login
COPY ./container/sshd_config /etc/ssh/sshd_config
RUN rm /etc/ssh/ssh_host_*_key
RUN mkdir /run/sshd
CMD ["/usr/sbin/sshd", "-D", "-e"]
