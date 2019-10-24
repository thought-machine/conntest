package(default_visibility = ["PUBLIC"])

subinclude("///third_party/subrepos/pleasings//remote")

go_get(
    name = "x_crypto",
    get = "golang.org/x/crypto/...",
    revision = "e84da0312774c21d64ee2317962ef669b27ffb41",
    deps = [":x_sys"],
)

go_get(
    name = "x_sys",
    get = "golang.org/x/sys",
    install = [
        "unix",
        "cpu",
    ],
    revision = "95b1ffbd15a57cc5abb3f04402b9e8ec0016a52c",
)


go_get(
    name = "logrus",
    get = "github.com/sirupsen/logrus",
    licences = ["mit"],
    revision = "v1.3.0",
    deps = [
        ":x_crypto",
        ":x_sys",
    ],
)

go_get(
    name = "difflib",
    get = "github.com/pmezard/go-difflib/...",
    revision = "792786c7400a136282c1664665ae0a8db921c6c2",
)

go_get(
    name = "go-spew",
    get = "github.com/davecgh/go-spew/spew",
    licences = ["isc"],
    revision = "ecdeabc65495df2dec95d7c4a4c3e021903035e5",
)

go_get(
    name = "objx",
    get = "github.com/stretchr/objx",
    licences = ["mit"],
    revision = "1a9d0bb9f541897e62256577b352fdbc1fb4fd94",
)

go_get(
    name = "yaml.v2",
    get = "gopkg.in/yaml.v2",
    revision = "f7716cbe52baa25d2e9b0d0da546fcf909fc16b4",
)

go_get(
    name = "testify",
    get = "github.com/stretchr/testify",
    install = [
        "assert",
        "require",
        "mock",
        "suite",
        "vendor/...",
    ],
    revision = "v1.4.0",
    deps = [
        ":difflib",
        ":go-spew",
        ":objx",
        ":yaml.v2",
    ],
)

git_repo(
    name = "tcpinfo-repo",
    repo = "https://github.com/brucespang/go-tcpinfo",
    revision = "e6cc7410d081080ec38065c00bb3b72f3b144913",
)
go_get(
    name = "tcpinfo",
    get = "github.com/brucespang/go-tcpinfo",
    licences = ["mit"],
    revision = "e6cc7410d081080ec38065c00bb3b72f3b144913",
)

go_get(
    name = "concurrent",
    get = "github.com/modern-go/concurrent/...",
    revision = "1.0.3",
)

go_get(
    name = "reflect2",
    get = "github.com/modern-go/reflect2/...",
    revision = "v1.0.1",
    deps = [
        ":concurrent",
    ],
)

go_get(
    name = "json-iterator",
    get = "github.com/json-iterator/go/...",
    revision = "v1.1.6",
    deps = [
        ":concurrent",
        ":reflect2",
    ],
)

go_get(
    name = "protobuf",
    get = "github.com/golang/protobuf/...",
    licences = ["bsd-3-clause"],
    revision = "v1.1.0",
    strip = [
        "proto/proto3_proto",
        "conformance",
    ],
)

go_get(
    name = "client_model",
    get = "github.com/prometheus/client_model/...",
    licences = ["apache-2.0"],
    revision = "fd36f4220a901265f90734c3183c5f0c91daa0b8",
    deps = [
        ":protobuf",
    ],
)

go_get(
    name = "perks",
    get = "github.com/beorn7/perks/...",
    licences = ["mit"],
    revision = "4c0e84591b9aa9e6dcfdf3e020114cd81f89d5f9",
)

go_get(
    name = "x_text",
    get = "golang.org/x/text/...",
    revision = "v0.3.0",
    strip = [
        "cmd",
        "message/pipeline",
    ],
)

go_get(
    name = "x_net",
    get = "golang.org/x/net/...",
    revision = "65e2d4e15006aab9813ff8769e768bbf4bb667a0",
    deps = [
        ":x_crypto",
        ":x_text",
    ],
)

go_get(
    name = "x_sync",
    get = "golang.org/x/sync/...",
    revision = "fd80eb99c8f653c847d294a001bdf2a3a6f768f5",
    deps = [
        ":x_net",
    ],
)

go_get(
    name = "procfs",
    get = "github.com/prometheus/procfs/...",
    licences = ["apache-2.0"],
    revision = "v0.0.3",
    deps = [
        ":x_sync",
    ],
)

go_get(
    name = "golang_protobuf_extensions",
    get = "github.com/matttproud/golang_protobuf_extensions/...",
    licences = ["apache-2.0"],
    revision = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
    deps = [
        ":protobuf",
    ],
)

go_get(
    name = "units",
    get = "github.com/alecthomas/units",
    licences = ["mit"],
    revision = "2efee857e7cfd4f3d0138cc3cbb1b4966962b93a",
)

go_get(
    name = "template",
    get = "github.com/alecthomas/template",
    licences = ["bsd-3-clause"],
    revision = "a0175ee3bccc567396460bf5acd36800cb10c49c",
)

go_get(
    name = "kingpin",
    get = "gopkg.in/alecthomas/kingpin.v2",
    licences = ["mit"],
    revision = "v2.2.6",
    deps = [
        ":template",
        ":units",
    ],
)

go_get(
    name = "prometheus_common",
    get = "github.com/prometheus/common",
    install = [
        "expfmt",
        "model",
        "internal/...",
        "log",
    ],
    licences = ["apache-2.0"],
    revision = "v0.6.0",
    deps = [
        ":client_model",
        ":golang_protobuf_extensions",
        ":kingpin",
        ":logrus",
        ":protobuf",
    ],
)

go_get(
    name = "prometheus",
    get = "github.com/prometheus/client_golang/...",
    licences = ["apache-2.0"],
    revision = "v1.0.0",
    deps = [
        ":client_model",
        ":json-iterator",
        ":perks",
        ":procfs",
        ":prometheus_common",
        ":protobuf",
        ":x_net",
    ],
)

go_get(
    name = "go-flags",
    get = "github.com/jessevdk/go-flags",
    licences = ["bsd-3-clause"],
    revision = "c0795c8afcf41dd1d786bebce68636c199b3bb45",
)
