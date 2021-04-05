FROM gcr.io/cloud-builders/bazel AS builder
COPY . .
RUN bazel build --stamp --workspace_status_command="./scripts/workspace-status.sh" //mpdev:mpdev

# Mpdev container requires docker, gsutil, and zip
FROM gcr.io/google.com/cloudsdktool/cloud-sdk:slim
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
RUN echo \
      "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian \
      $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update -y
RUN apt-get install -y zip docker-ce docker-ce-cli containerd.io
COPY --from=builder bazel-bin/mpdev/mpdev_/mpdev /usr/bin/mpdev
ENTRYPOINT ["/usr/bin/mpdev"]
