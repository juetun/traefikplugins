module github.com/juetun/traefikplugins

go 1.15

replace (
	github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6
	go.etcd.io/bbolt v1.3.6 => github.com/coreos/bbolt v1.3.6
)

require github.com/juetun/base-wrapper v0.0.126 // indirect
