FROM scratch
COPY goblin /usr/bin/goblin
ENTRYPOINT ["/usr/bin/goblin"]
WORKDIR /goblin