load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "67b4d1f517ba73e0a92eb2f57d821f2ddc21f5bc2bd7a231573f11bd8758192e",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.50.0/rules_go-v0.50.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.50.0/rules_go-v0.50.0.zip",
    ],
)

git_repository(
    name = "deploymentmanager-autogen",
    remote = "https://github.com/GoogleCloudPlatform/deploymentmanager-autogen",
    tag = "0.1.0",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.23.0")

http_archive(
    name = "bazel_gazelle",
    sha256 = "b760f7fe75173886007f7c2e616a21241208f3d90e8657dc65d36a771e916b6a",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.39.1/bazel-gazelle-v0.39.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.39.1/bazel-gazelle-v0.39.1.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("//:repos.bzl", "go_repositories")

# gazelle:repository_macro repos.bzl%go_repositories
go_repositories()

gazelle_dependencies()


http_archive(
    name = "com_google_protobuf",
    sha256 = "cf6b0458b280081a07b236278abf48321bac0a268bc14ef3510da624ea243993",
    strip_prefix = "protobuf-21.8",
    urls = ["https://github.com/protocolbuffers/protobuf/releases/download/v21.8/protobuf-all-21.8.zip"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()
