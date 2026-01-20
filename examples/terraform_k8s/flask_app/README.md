# How to Use the Terraform K8s App Example

This guide walks through how to use the example Kubernetes app package in this directory. To use this example, first create a Kubernetes App (using Terraform) Listing in Producer Portal.

Clone this repository and cd into this directory (examples/terraform\_k8s/flask\_app).

## Preparing Artifact Registry Assets

If you haven't yet, [enable the Artifact Registry API](https://cloud.google.com/artifact-registry/docs/enable-service) for your project, and create
a **multi-region** repository in the ``us'' multi-region. The repo must be under the us-docker.pkg.dev host.

### Build and upload container image

Use docker (or podman), build the Flask App container image and push it to the AR repo:

```shell
AR_REPO="us-docker.pkg.dev/your-project/your-repo"
AR_PATH_PREFIX="k8s-app-example"
AR_IMAGE_NAME="app"
TAG="1.0"
APP_IMAGE_FULL_URI="${AR_REPO}/${AR_PATH_PREFIX}/${AR_IMAGE_NAME}:${TAG}"

docker build . -t "$APP_IMAGE_FULL_URI"
docker push "$APP_IMAGE_FULL_URI"
```

Set the required [service name annotation](https://cloud.google.com/marketplace/docs/partners/migrations/container-image-annotations) using crane or gcrane:

```shell
SERVICE_NAME="your-listing.endpoints.your-project.cloud.goog"

crane mutate -a "com.googleapis.cloudmarketplace.product.service.name=${SERVICE_NAME}" "$APP_IMAGE_FULL_URI"
```

### Build and upload Helm chart

```shell
AR_CHART_NAME="chart"
helm package ./chart --version "$TAG"

helm push "chart-${TAG}.tgz" "oci://${AR_REPO}/${AR_PATH_PREFIX}"
```

You should see two items in your AR repo:

* us-docker.pkg.dev/your-project/your-repo/k8s-app-example/app:1.0  
* us-docker.pkg.dev/your-project/your-repo/k8s-app-example/chart:1.0

## Prepare the Terraform Package

```shell
cd tf
```

You will need to modify schema.yaml to reference the image URI you chose from the previous step. For example, the image path (relative to repo root) is gh-ui-example/app. Therefore, the schema.yaml needs to contain:

```
images:
  k8s-app-example/app:
    variables:
      image_repo:
        type: REPO_WITH_REGISTRY_WITH_NAME
      image_tag:
        type: TAG
```

Now you can zip the module and upload it to GCS with gsutil:

```shell
zip -r ../tf.zip .
gcloud storage cp ../tf.zip gs://your-bucket/tf.zip
```

## Proceed to set up in Producer Portal

In Producer Portal, enter your chart’s URI (without the tag): us-docker.pkg.dev/your-project/your-repo/k8s-app-example/chart, and follow the guide to set up your releases.
