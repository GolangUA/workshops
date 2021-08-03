To generate a pb file use (protoc bin must be installed):

    protoc --go_out=. filename.proto

To generate grpc service file

    protoc --go_out=plugins=grpc:package filename.proto    

or

    protoc --proto_path=path_to_proto_files/*.proto --go_out=plugins=grpc:pb
    protoc --proto_path=models/*.proto --go_out=plugins=grpc:generated_protoc
Concrete case

    protoc --go_out=plugins=grpc:chat chat.proto


export PATH="$PATH:$(go env GOPATH)/bin"


protoc --go_out=../pb/ --go_opt=paths=source_relative \
    --go-grpc_out=../pb/ --go-grpc_opt=paths=source_relative \
    *.proto

###Some info on gRPC

gRPC utilizes HTTP/2 whereas REST utilizes HTTP 1.1

gRPC utilizes the protocol buffer data format as opposed to the standard JSON data format that is typically used within REST APIs

With gRPC you can utilize HTTP/2 capabilities such as server-side streaming, client-side streaming or even bidirectional-streaming should you wish.