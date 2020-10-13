module github.com/sschwartz96/syncapod

go 1.15

//replace github.com/sschwartz96/stockpile => /home/sam/go/src/github.com/sschwartz96/stockpile
replace github.com/sschwartz96/stockpile => C:/users/sam/go/src/github.com/sschwartz96/stockpile

require (
	github.com/golang/protobuf v1.4.2
	github.com/sschwartz96/minimongo v0.2.3 // indirect
	github.com/sschwartz96/stockpile v0.2.3
	github.com/tcolgate/mp3 v0.0.0-20170426193717-e79c5a46d300
	go.mongodb.org/mongo-driver v1.4.2
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	golang.org/x/sys v0.0.0-20201008064518-c1f3e3309c71 // indirect
	google.golang.org/grpc v1.33.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
