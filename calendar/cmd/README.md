This package represents starting points of your application.

If you need different entrypoints to serve multiple servers as different processes or even cli
entrypoint - then you should have something simmilar to this:

    cmd/cli/main.go
    cmd/grpc-server/main.go
    cmd/http-server/main.go