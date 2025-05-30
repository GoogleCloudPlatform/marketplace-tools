# Preparing a Terraform K8s App Package


## Prerequisites

You have an existing K8s App packaged as a Helm Chart. You should have a values.yaml file that parametrizes the deployment.

Some requirements regarding to Charts and Docker Images:



*   Only one chart is allowed, all dependency charts need to be downloaded under the same folder.  Marketplace does not download charts from remote locations.
*   All Docker Images need to be included in the solution.  Marketplace does not pull Docker Images outside your GCP projects. 


## Step 1: Push Your Chart and Images to AR

If you don’t have one yet, create an Artifact Registry Repo in your Producer project. You can create a single repo for all of your Marketplace Terraform K8s App products, or have individual repositories.  \
 \
If you choose to place your images across multiple products in a single repository, it is recommended that you place all images and charts for each product under its own “folder”. \
 \
Your Artifact Registry repo may look like this: 


```
us-docker.pkg.dev/[partner-project]/[repo-name]
# example:
us-docker.pkg.dev/acme-ai-studio/gcloud-marketplace-repo/
```

Assuming you have a product to provide booking service as TravelAgent, and you the following assets locally

*   a docker image called **_travel-agent-image_**
*   a helm chart under folder **_~/dist/travel-agent-chart_**

Assuming you already tested out deploying this to a GKE cluster, everything works great.  Now you want to publish the chart and image to AR.  As a convention, it is recommended to put the chart and images under the same folder, say **_travel-agent-app_** in this example.

**Push the Docker Image to AR:**

**NOTE: **Make sure you put the Marketplace annotation into your image: [Migrate to annotating container images with their service name | Google Cloud Marketplace Partners](https://cloud.google.com/marketplace/docs/partners/migrations/container-image-annotations) 

example:

```sh
gcrane mutate -a "com.googleapis.cloudmarketplace.product.service.name=services/travel-agent-app.endpoints.your-company-name.cloud.goog" us-docker.pkg.dev/your-project/your-repo/travel-agent@sha256:667842b0a4a75983e56bcf78a305ffa4603814c8dc533e6f8f92779dce4b0015
```

Please note, this command will generate a new version with new sha, you may need to update the tag.


```sh
# tag local with AR destination repo path
docker tag travel-agent us-docker.pkg.dev/[your-gcp-project-id]/[your-repo]/travel-agent:1.0

# to verify, this will show a new entry you just tagged
docker image list

# push to AR
docker push us-docker.pkg.dev/[you-gcp-project-id]/[your-repo]/travel-agent:1.0

# then list the new image to verify
gcloud artifacts docker images list
```


At the minimum, your image should be tagged with a semantic version of MAJOR.MINOR. You can add other tags (MAJOR.MINOR.PATCH) as you like, but the MAJOR.MINOR version must exist and be the same across all images and Helm Charts used in this version of your app.

For details about how to push Docker Images to AR, reference public doc: [Push and pull images | Artifact Registry documentation | Google Cloud](https://cloud.google.com/artifact-registry/docs/docker/pushing-and-pulling#pushing). 

NOTE: you can add multiple images.

**Push the Helm Chart to AR:**

Assume your helm chart is under folder **_~/dist/travel-agent-chart,_** and the version is 1.0 in Chart.yaml,


```sh
# this will create a travel-agent-chart-1.0.tgz file
cd ~/dist
helm package travel-agent-chart --version 1.0 --app-version 1.0

# push to AR, under [your-repo]
helm push travel-agent-chart-1.0.tgz oci://us-docker.pkg.dev/[your-gcp-project-id]/[your-repo]
# Example:
# helm push travel-agent-chart-1.0.tgz oci://us-docker.pkg.dev/acme-ai-studio/gcloud-marketplace-repo/

# this will push the travel-agent-chart to the repo, also add a tag 1.0
# the 1.0 tag is from the version in Chart.yaml
# output will look like:
# Pushed: us-docker.pkg.dev/acme-ai-studio/gcloud-marketplace-repo/travel-agent-chart:1.0
# Digest: ...
```

Now you have your Docker Images and Chart on AR.

NOTE: for simplicity, Marketplace requires the tag for the Chart and Docker Images it refers to to have the same MAJOR.MINOR version tag.  In the example, all Charts and container images should have the tag “1.0”.  Later, if you have a major upgrade, you can bump the Chart version to _2.0_, then all Docker Images it refers to also need to have tags as _2.0_.  For minor changes, update the MINOR part.  Besides the MAJOR.MINOR tag, you can have other tags (e.g., 1.0.2, 1.0.2+patch1) as you like for version management.


## Step 2: Create a Terraform Wrapper Module for your Helm Chart.

We will use Terraform to deploy our solution, and in Terraform, we use Helm Provider to deploy Helm Charts. To be specific for Google Cloud, we use InfraManager, a hosted Terraform infrastructure for deployment management.

Download the starter module from [GitHub](https://github.com/GoogleCloudPlatform/marketplace-tools/tree/master/docs/terraform-k8s-app/starter-terraform-module).  

Mainly you want to change the **helm.tf** to customize input variables. 

In most cases your Helm Chart will have some parameters as inputs, for example, port number, kubernetes service account name, etc.  Our goal is that you don’t need to modify your Helm Chart files to adapt to Terraform.  Instead, for each parameter in Helm Chart, you can simply create a Terraform variable as a proxy.  Your customers will be able to provide Terraform variables later either in the command line or UI configuration page for deployment.

Terraform takes input from consumers through Terraform Variables. It then passes the values of Terraform Variables to the Helm Chart (to _override_ values in values.yaml) through the `helm_release` resource.

You can use `set` blocks to override individual values in the Chart’s values.yaml. This is equivalent to running `helm install --set name1=value1 --set name2=value2`...

Let’s say you have a values.yaml in your Helm Chart like this:


```yaml
# Default values for travel-agent-chart.
# Declare variables to be passed into your K8s templates.
replicaCount: 1

image:
  repository: us-docker.pkg.dev/your-proj-id/your-repo/travel-agent
  tag: "1.0"

service:
  type: LoadBalancer
  port: 80
  targetPort: 8251
```


You can define corresponding variables in Terraform for those your customers may want to customize, like this in **variables.tf**:


```
# variables.tf
variable "replica_count" {
  type = number
  default = 1
}

variable "image_repo" {
  type = string
  default = "us-docker.pkg.dev/your-proj-id/your-repo/travel-agent"
}

variable "image_tag" {
  type = string
  default = "1.0"
}

variable "service_port" {
  type = number
  default = 8080
}
```


Then in **helm.tf**, you can assign the values of Terraform variables to those defined in Helm Chart values.yaml.


```
# helm.tf
resource "helm_release" "my_app" {
  # ...

  set { 
    name = "replicaCount"
    value = jsencode(var.replica_count)
  }

  set {
    name = "image.repository"
    value = jsencode(var.image_repo)
  }

  set {
    name = "image.tag"
    value = jsencode(var.image_tag)
  }

  set {
    name = "service.port"
    value = jsencode(var.service_port)
  }

  # ...
}
```


You can customize other parts of the Terraform module as you like.

NOTE:


```
  repository = var.helm_chart_repo
  chart      = var.helm_chart_name
  version    = var.helm_chart_version
```



## Step 3: Create a schema YAML file for your package

When you publish your solution, Marketplace will copy your Helm Chart and Docker Images to the Google owned repository for availability and security reasons.  Customers will deploy your solution based on the copy in the Google owned repository.

During publishing, the helm chart location in **_helm.tf_ **will be automatically replaced with that to the Google owned location.

For Docker Images, our goal is that you don’t need to change your existing Helm Chart files. You will need to create a mapping file called **_schema.yaml_** to define Terraform variables for each Docker Image you have, and those Terraform variables will be used to override your Helm Chart variables.  The values of those Terraform variables you defined in **_schema.yaml_** are populated by our system during the publishing process.

A prerequisite is that you need to have Docker Image URLs in your Helm values.yaml file in some way.  The whole URL needs to be represented in values.yaml, otherwise if part of the URL is hard-coded in your k8s template files, it won’t be correctly overridden during the publishing process.

In our example above, we have one image defined in values.yaml


```yaml
# values.yaml
image:
  repository: us-docker.pkg.dev/your-proj-id/your-repo/travel-agent
  tag: "1.0"
```

In deployment.yaml, you could use those variables like this: 

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}-deployment
spec:
    ...
    template:
    ...
    spec:
      containers:
      - name: travel-agent-container
        # NOTE: make sure no hard-coded strings here, all variables
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
```


 \
Another popular way to define the Docker Image related variables could be:


```yaml
# break the URL into 3 parts
image:
  registry: docker.io
  repository: bitnami/airflow
  tag: 2.10.5-debian-12-r6

# similar to our example, 2 parts
image:
  # image.repository -- Repository to use for Datadog Operator image
  repository: gcr.io/datadoghq/operator
  # image.tag -- Define the Datadog Operator version to use
  tag: 1.13.0
```

Instructions to create the **_schema.yaml_**: 


1. Create a YAML file named schema.yaml in the Terraform module folder.  The starter zip has one similar to this:

```yaml
images:
  # this is the image key
  travel-agent:
    variables:
      image_repo:
        type: REPO_WITH_REGISTRY_WITH_NAME
      image_tag:
        type: TAG
```

2. For each Docker Image, add one entry under the **_images_**, the value is the part after the AR repo name.  In this example, for the AR URL: <code>us-docker.pkg.dev/your-proj-id/your-repo/<strong>travel-agent</strong>, </code>the image key is <code>travel-agent</code>.  If you have a Docker Image URL as <code>us-docker.pkg.dev/your-proj-id/your-repo/<strong>my-product-A/webserver-image</strong>, </code>then the image key is <code>my-product-A/webserver-image</code>.
3. Note the URIs to your Docker Images, which might resemble us-docker.pkg.dev/project-id/repo-name/path/to/chart:tag. We will call
    1. the “us-docker.pkg.dev” part as the <strong>registry</strong>
    2. the “project-id/repo-name” part as the <strong>repo</strong>
    3. the “path/to/chart” part as the <strong>name</strong>. 
    4. the tag part as the <strong>tag</strong>.
4. For each Terraform Variable that references a segment of this URI, declare it under <code>variables</code> with a <code>type</code> sub-property.  This will tell our system to set the corresponding value to the Terraform var.
5. Figure out the type for each variable you created. The type indicates which part of the Docker Image URI that your variable is supposed to represent. Take <code>us-docker.pkg.dev/project-id/repo-name/path/to/chart:tag</code> as an example.

<table>
  <tr>
   <td>
<strong>Type</strong>
   </td>
   <td><strong>Corresponding Value</strong>
   </td>
  </tr>
  <tr>
   <td>REGISTRY
   </td>
   <td>us-docker.pkg.dev
   </td>
  </tr>
  <tr>
   <td>REPO_WITHOUT_REGISTRY_WITHOUT_NAME
   </td>
   <td>project-id/repo-name
   </td>
  </tr>
  <tr>
   <td>REPO_WITHOUT_REGISTRY_WITH_NAME
   </td>
   <td>project-id/repo-name/path/to/chart
   </td>
  </tr>
  <tr>
   <td>REPO_WITH_REGISTRY_WITHOUT_NAME
   </td>
   <td>us-docker.pkg.dev/project-id/repo-name
   </td>
  </tr>
  <tr>
   <td>REPO_WITH_REGISTRY_WITH_NAME
   </td>
   <td>us-docker.pkg.dev/project-id/repo-name/path/to/chart
   </td>
  </tr>
  <tr>
   <td>NAME
   </td>
   <td>path/to/chart
   </td>
  </tr>
  <tr>
   <td>TAG
   </td>
   <td>tag
   </td>
  </tr>
</table>


So let’s revisit the sample **_schema.yaml_**:


```yaml
images:
  travel-agent:
    variables:
      image_repo:
        type: REPO_WITH_REGISTRY_WITH_NAME
      image_tag:
        type: TAG
```

This will tell the Marketplace system to populate corresponding values into the two tf variables: **image_repo** and **image_tag**.

The image_repo and image_tag also need to be defined in the variables.tf as mentioned before: 

```hcl
variable "image_repo" {
  type = string
  default = "us-docker.pkg.dev/your-proj-id/your-repo/travel-agent"
}

variable "image_tag" {
  type = string
  default = "1.0"
}
```


In the helm.tf, we assign the tf variable values to Helm variables:


```
  set {
    name = "image.repository"
    value = jsencode(var.image_repo)
  }

  set {
    name = "image.tag"
    value = jsencode(var.image_tag)
  }
```


As the final result, the value in K8s template file, say deployment.yaml, which uses the variables, is replaced with the URL of the Docker Image in Google-owned AR:


```yaml
...
spec:
...
  template:
   ...
    spec:
      containers:
      - name: travel-agent-container
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
```


 \
NOTE: you can have multiple image records in the **_schema.yaml_**.


## Step 4: Test Deployment Locally

With the terraform module ready, please try to deploy your solution into your GCP project, to check if there are any issues.


```
# Initialize Terraform
terraform init

# Review the changes before applying, depending on your variables.tf, set the corresponding values in the CLI
terraform plan \
  -var="project_id=your-gcp-project" \
  -var="cluster_name=your-cluster-name" \
  -var="cluster_location=us-central1-a"

# Apply the configuration
terraform apply \
  -var="project_id=your-gcp-project" \
  -var="cluster_name=your-cluster-name" \
  -var="cluster_location=us-central1-a" \
  -auto-approve
```


 \
Marketplace will do a verification based on the Terraform modules you provide, and with a set of test variables.  Please provide test variables in **marketplace_test.tfvars**, and put default values there for verification purposes:


```
cluster_location = us-central-1
cluster_name     = test_cluster
create_cluster   = true
# make sure no project_id, helm_chart_{repo,name,version} or ANY variable declared in schema.yaml in this file
```


NOTE: for the variables helm_chart_repo, helm_chart_name, helm_chart_version, please provide the values for local testing either in cli or as default values in variables.tf.


## Step 5: Upload Terraform Module to GCS

Package your Terraform module along with the schema.yaml file into a single zip file. Upload it to a GCS versioned bucket ([how to create a versioned bucket on GCS](https://cloud.google.com/storage/docs/using-object-versioning#set-console)) under your partner project in Google Cloud.  Make sure the bucket supports versioning when it is created.

Assume you have a folder **_travel-agent-terraform, _**and files under the folder, e.g. main.tf, schema.yaml, etc., please make sure when you create the zip file, DO NOT have the _travel-agent-terraform _folder included in the zip file.


```sh
cd travel-agent-terraform
zip -r ../travel-agent-tf.zip *
gsutil cp ../travel-agent-tf.zip gs://your-bucket/
```

