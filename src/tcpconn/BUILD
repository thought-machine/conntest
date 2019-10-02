go_library(
    name = "tcpconn",
    srcs = ["tcpconn.go"],
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go:logrus",
        "//third_party/go:prometheus",
        "//third_party/go:tcpinfo",
    ],
)

go_test(
    name = "tcpconn_test",
    srcs = ["tcpconn_test.go"],
    # visibility = ["//conntest/..."],
    deps = [
        ":tcpconn",
        "//third_party/go:logrus",
        "//third_party/go:testify",
    ],
)
