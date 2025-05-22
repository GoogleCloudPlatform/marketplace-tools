# The marketplace_test.tfvars file is used to validate the Terraform template.
# Marketplace will validate your product with this file as its `--var-file`
# argument.
#
# Do not include the following variables in marketplace_test.tfvars, as they
# will be provided by Marketplace:
#
# - project_id
# - helm_chart_repo
# - helm_chart_name
# - helm_chart_version
# - Any variables declared in schema.yaml

create_cluster   = true
cluster_name     = "marketplace-test"
cluster_location = "us-central1"

# TODO: Add values used for testing.
