FROM ubuntu:wily
MAINTAINER caktux

ENV DEBIAN_FRONTEND noninteractive

# Usual update / upgrade
RUN apt-get update
RUN apt-get upgrade -q -y
RUN apt-get dist-upgrade -q -y

# Install Vector
RUN apt-get install -q -y software-properties-common
RUN add-apt-repository ppa:vector/vector
RUN add-apt-repository ppa:vector/vector-dev
RUN apt-get update
RUN apt-get install -q -y gvec

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/usr/bin/gvec"]
