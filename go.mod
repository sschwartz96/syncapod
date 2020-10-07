module github.com/sschwartz96/syncapod

go 1.15

//replace github.com/sschwartz96/minimongo => /home/sam/go/src/github.com/sschwartz96/minimongo

require (
	github.com/golang/protobuf v1.4.2
	github.com/sschwartz96/minimongo v0.2.3
	github.com/tcolgate/mp3 v0.0.0-20170426193717-e79c5a46d300
	go.mongodb.org/mongo-driver v1.4.1
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
