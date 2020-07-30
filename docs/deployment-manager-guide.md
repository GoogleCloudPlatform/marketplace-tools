## Generating a deployment manager template

The mpdev tool can autogenerate deployment manager templates from an
[autogen specification](./autogen-reference.md). Follow these steps to
build create a deployment manager template.

1. Use the `mpdev pkg get` command to retrieve a preconfigured autogen
specification. See an example [here](../examples/deployment-manager/autogen/singlevm/README.md).
1. Use `mpdev cfg set` to customize values in the retrieved autogen
specification, or edit the spec manually.
3. Execute `mpdev apply` to generate the deployment manager template.

## Upload solution to partner portal

Open [Partner Portal](https://pantheon.corp.google.com/partner/solutions) and 
select your solution from the list of solutions, then follow these steps:

1. Next to **Deployment Package**, click **Edit**.
1. Select to **Upload a Package**, and then click **Continue**
1. Select the deployment package created by mpdev for **Upload a package**,
uncheck the **Metadata selection** box and then click **Continue**.
1. Click **Save** to save changes.

**Warning:** Unchecking the **Metadata selection** box is crucial, so that
mpdev does not override the solution metadata configured in earlier steps of
the partner portal wizard.