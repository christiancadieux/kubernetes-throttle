FROM fedora:25

ENV SERVER_PORT 9191
ADD ./throttle_server /throttle_server

ENTRYPOINT ["/throttle_server"]

