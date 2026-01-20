This example meets the
[Marketplace requirements](https://cloud.google.com/marketplace/docs/partners/vm/configure-terraform-deployment#requirements_for_custom_terraform_modules)
for a Terraform **CLI Only** module

To validate your module, Marketplace will execute:

```
terraform init
terraform plan -var-file=marketplace_test.tfvars -var project_id=<test-project>
```

The `marketplace_test.tfvars` can be used to set variables to default values for
testing only purposes.

For a Marketplace Partner to reuse this module, they must:

1.  Declare "image" as a variable in Producer Portal.
1.  Replace the default value of "image" with the image declared in Producer
    Portal
1.  Zip the module and upload to GCS

```
zip terraform.zip * -r
gcloud storage cp terraform.zip gs://<YOUR-BUCKET>/<FOLDER>/terraform.zip
```

Additionally, Partners should include a README.md containing instructions on how
to deploy their product.
