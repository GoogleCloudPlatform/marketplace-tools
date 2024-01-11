> [!IMPORTANT]
> If your module can be used directly as a root module, consider following the [simple_cli_only](../simple_cli_only) example instead.

This example meets the
[Marketplace requirements](https://cloud.google.com/marketplace/docs/partners/vm/configure-terraform-deployment#requirements_for_custom_terraform_modules)
for a Terraform **CLI Only** module

If the `examples/marketplace_test` folder exists, marketplace will validate your
module by executing:

```
terraform -chdir=examples/marketplace_test init
terraform -chdir=examples/marketplace_test plan -var project_id=<test-project>
```

The module in `examples/marketplace_test` must reference your module in the root
directory as shown below:

```
# The test module references the module in the root directory
module "test" {
  source = "../.."
}
```

This test module can be used to create resources that the module in your root
directory assumes already exists. For instance, this example creates a second
network interface that is passed as an argument to the module in the root
directory.

For a Marketplace Partner to reuse this module, they must:

1.  Declare "image" as a variable in Producer Portal.
1.  Replace the default value of "image" with the image declared in Producer
    Portal
1.  Zip the module and upload to GCS

```
zip terraform.zip * -r
gsutil cp terraform.zip gs://<YOUR-BUCKET>/<FOLDER>/terraform.zip
```

Additionally, Partners should include a README.md containing instructions on how
to deploy their product.
