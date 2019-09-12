module tedis

require (
	github.com/coreos/etcd v3.3.15+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548
	github.com/eapache/channels v1.1.0
	github.com/eapache/queue v1.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/golang/snappy v0.0.1
	github.com/gorilla/mux v1.7.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.8.1 // indirect
	github.com/juju/errors v0.0.0-20190806202954-0232dcc7464d
	github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8 // indirect
	github.com/juju/testing v0.0.0-20190723135506-ce30eb24acd2 // indirect
	github.com/montanaflynn/stats v0.5.0 // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pingcap/parser v0.0.0-20190806084718-1a31cabbaef2
	github.com/pingcap/pd v0.0.0-20190711034019-ee98bf9063e9
	github.com/pingcap/tidb v0.0.0-20190828105439-836982c617fb
	github.com/prometheus/client_golang v0.9.2
	github.com/remyoudompheng/bigfft v0.0.0-20190728182440-6a916e37a237 // indirect
	github.com/sirupsen/logrus v1.3.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/twinj/uuid v1.0.0
	github.com/unrolled/render v1.0.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	google.golang.org/grpc v1.23.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v2 v2.2.2
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace github.com/Sirupsen/logrus v1.4.2 => github.com/sirupsen/logrus v1.0.6

go 1.13
