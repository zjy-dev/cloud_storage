module user

go 1.16

require (
	github.com/J-Y-Zhang/cloud-storage/common v0.0.0-20220222005224-a03003cba602
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20220113124808-70ae35bab23f // indirect
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/asim/go-micro/plugins/config/source/consul/v4 v4.0.0-20220216084336-d47c2d984b92
	github.com/asim/go-micro/plugins/registry/consul/v4 v4.0.0-20220216084336-d47c2d984b92
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/consul/api v1.12.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.9.7 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/miekg/dns v1.1.46 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	go-micro.dev/v4 v4.6.0
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/mysql v1.3.2
	gorm.io/gorm v1.23.1
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
