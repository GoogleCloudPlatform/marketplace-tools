load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "a8d6b1b354d371a646d2f7927319974e0f9e52f73a2452d2b3877118169eb6bb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
    ],
)

git_repository(
    name = "deploymentmanager-autogen",
    remote = "https://github.com/GoogleCloudPlatform/deploymentmanager-autogen",
    tag = "0.1.0",
)

http_archive(
    name = "com_google_protobuf",
    sha256 = "aed089110977f7cabda19d550b4d503eb93fa0f73c7a470c4695f27b63029544",
    strip_prefix = "protobuf-3.9.1",
    urls = ["https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protobuf-all-3.9.1.zip"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

http_archive(
    name = "bazel_gazelle",
    sha256 = "7fc87f4170011201b1690326e8c16c5d802836e3a0d617d8f75c3af2b23180c4",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/0.18.2/bazel-gazelle-0.18.2.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.18.2/bazel-gazelle-0.18.2.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("//:repos.bzl", "go_repositories")

# gazelle:repository_macro repos.bzl%go_repositories
go_repositories()

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "95d39fd84ff4474babaf190450ee034d958202043e366b9fc38f438c9e6c3334",
    strip_prefix = "rules_docker-0.16.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.16.0/rules_docker-v0.16.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

http_file(
    name = "gcloud_sdk_apt_key",
    urls = [
        "https://packages.cloud.google.com/apt/doc/apt-key.gpg",
    ],
    sha256 = "d9b7ded07f7ff37138f842e98c30f8ff54c185fc0013d47b5df749b6b4391e87",
    downloaded_file_path = "gcloud.pub.gpg",
)

http_file(
    name = "bazel_apt_key",
    urls = [
        "https://bazel.build/bazel-release.pub.gpg",
    ],
    sha256 = "547ec71b61f94b07909969649d52ee069db9b0c55763d3add366ca7a30fb3f6d",
    downloaded_file_path = "bazel-release.pub.gpg",
)

http_archive(
    name = "ubuntu1604",
    strip_prefix = "base-images-docker-36456edd3cc5a4d17852439cdcb038022cd912e5/ubuntu1604",
    urls = ["https://github.com/GoogleContainerTools/base-images-docker/archive/36456edd3cc5a4d17852439cdcb038022cd912e5.tar.gz"],
    sha256 = "eb50020790e22538676e17c0242b9272bdb5c81f7bc6b128a4abfa7ad31faf5b"

)

load("@ubuntu1604//:deps.bzl", ubuntu1604_deps = "deps")

ubuntu1604_deps()