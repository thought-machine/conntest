subinclude("///third_party/subrepos/pleasings//docker")

subinclude("///third_party/subrepos/pleasings//k8s")


go_binary(
    name = "conntest",
    srcs = ["main.go"],
    static = False,
    deps = [
        "//src/srvendpoints:srvendpoints",
        "//src/tcpconn:tcpconn",
        "//third_party/go:logrus",
        "//third_party/go:prometheus",
        "//third_party/go:go-flags",
    ],
)

docker_image(
    name = "conntest_alpine",
    srcs = [
        ":conntest",
    ],
    dockerfile = "Dockerfile-conntest",
    image = "conntest",
)

k8s_config(
    name = "k8s",
    srcs = [
        "src/k8s/conntest.yaml",
        "src/k8s/conntest-svc.yaml",
    ],
    containers = [
        ":conntest_alpine",
    ],
)
