module github.com/sdslabs/katana

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gliderlabs/ssh v0.2.2
	github.com/go-git/go-git/v5 v5.4.2
	github.com/gofiber/fiber/v2 v2.1.0
	github.com/golang/protobuf v1.5.2
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/hashicorp/terraform-exec v0.15.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	go.mongodb.org/mongo-driver v1.5.3
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9 // indirect
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1 // indirect
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.2
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace google.golang.org/grpc => github.com/rohithvarma3000/grpc-go v1.44.0
