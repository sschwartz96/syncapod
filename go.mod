module github.com/sschwartz96/syncapod

go 1.15

replace github.com/sschwartz96/minimongo => /home/sam/go/src/github.com/sschwartz96/minimongo

require (
	github.com/golang/protobuf v1.4.2
	github.com/schollz/closestmatch v2.1.0+incompatible
	github.com/sschwartz96/minimongo v0.1.2
	github.com/tcolgate/mp3 v0.0.0-20170426193717-e79c5a46d300
	go.mongodb.org/mongo-driver v1.4.1
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
)
