FROM ubuntu:latest
LABEL authors="gilang"

ENTRYPOINT ["top", "-b"]