# SaaS Test Tool
This page describes how to set up a billing account, project, a service account, and use BigQuery billing export, on Google Cloud to run the SaaS Test Tool using `mpdev`.

## Setup


### Configure your test billing account
To run the SaaS Test Tool, you must create a test billing account. Usage with this billing account is 100% discounted to prevent unexpected charges for testing.


#### Create your billing account
To create a new billing account, follow the instructions at [Create a new Cloud Billing account](https://cloud.google.com/billing/docs/how-to/manage-billing-account#create_a_new_billing_account).


#### Configure a 100% Discount for your billing account
> **🛈 NOTE:** The background process takes 24 hours to complete.


> **🛈 NOTE:** You must use the same project that created the SaaS product you’re testing.


1. Navigate to [Billing](https://console.cloud.google.com/billing) and select the billing account you created previously.
1. In the left navigation panel, click **Account management**.
1. Copy the **Billing account ID**.
1. Visit the [Configure Reports and Testing Billing Accounts](https://console.cloud.google.com/producer-portal/reports) page.
1. Click **ADD TEST BILLING ACCOUNT**.
1. In the **Test Account** field, paste the Billing account ID you copied previously.
1. Click **DONE**.
1. Click **SAVE**.


### Procurement Account
> **🛈 NOTE:** You must be a **Billing Account Administrator** on the test billing account you created.


To connect a SaaS product to a customer, you must use a procurement account. When you use the SaaS Test Tool, the SaaS Test Tool takes on the role of the partner and customer, requiring a procurement account.


#### Create a procurement account
To create an **unapproved** procurement account, place an order for the SaaS product using [Google Cloud Marketplace](https://console.cloud.google.com/marketplace).


#### Approve the procurement account
Please refer to the [backend integration](https://cloud.google.com/marketplace/docs/partners/integrated-saas/backend-integration?_ga=2.266876492.-359088086.1642023532#create-account) documentation for more information on approving procurement accounts.


### Configuring your Google Cloud project
This page details how to create a new Google Cloud project solely for running the SaaS Test Tool.


#### Create a new Google Cloud project
> **🛈 NOTE:** You must link this new project with the test billing account you created previously.


To create a new Google Cloud project, follow the instructions at [Creating and managing projects](https://cloud.google.com/resource-manager/docs/creating-managing-projects).


#### Enable the Cloud Commerce Consumer Procurement API
1. Go to the [Cloud Console](https://console.cloud.google.com/).
1. Select the project you created previously.
1. Go to the [Cloud Commerce Consumer Procurement API](https://console.developers.google.com/apis/api/cloudcommerceprocurement.googleapis.com/overview).
1. Click **ENABLE**.


#### Enable the Cloud Commerce Partner Procurement API
1. Go to the [Cloud Console](https://console.cloud.google.com/).
1. Select the project you created previously.
1. Select the [Cloud Commerce Partner Procurement API](https://console.developers.google.com/apis/api/cloudcommerceprocurement.googleapis.com/overview).
1. Click **ENABLE**.


### Configuring your service account
To run the SaaS Test Tool, you must use a service account.


#### Create your service account
> **🛈 NOTE:** You must use the same project you created previously.

To create a new service account, follow the steps at [Creating and managing service accounts](https://cloud.google.com/iam/docs/creating-managing-service-accounts).


#### Grant Billing Administrator permissions to your service account
1. From the [Billing](https://console.cloud.google.com/billing) page, select the billing account you created previously.
1. In the left navigation panel, click **Account management**.
1. In the right pane, in the **INFO PANEL**, click **+ADD PRINCIPAL**.
1. Set the new principals to be the service account and email address you created previously.
1. Select the `Billing Account Administrator` IAM role


#### Grant your service account access to Billing metering results
1. Go to the [Service accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) page.
1. Select the project you created previously.
1. Select the service account you created previously.
1. Copy the **Email** field.
1. Go to the [IAM](https://console.cloud.google.com/iam-admin/iam) page.
1. Click **+ADD**.
1. In the **New principals** text field, paste the service account email address you copied previously.
1. Select the `BigQuery Admin` for the **role**.
1. Click **SAVE**.


#### Link your service account to call the Procurement API and report usage
1. Go to the [Service accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) page.
1. Select the project you created previously.
1. Select the service account you created previously.
1. Copy the **Email** field.
1. Go to the **Technical integration** section of your test product.
1. Under **Link a service account to call Procurement API**, click **+ADD SERVICE ACCOUNT**.
1. Paste the service account email that you copied previously.
1. At the end of the input field, click **LINK**.
1. Under **Add 'roles/servicemanagement.serviceController' to a service account**, click **+ADD SERVICE ACCOUNT**.
1. Paste the service account email you copied previously.
1. At the end of the input field, click **LINK**.


### Setup BigQuery billing export
> **🛈 NOTE:** You must use the same test billing account that you configured previously


To set up billing export, follow the steps at [Set up Cloud Billing data export to BigQuery](https://cloud.google.com/billing/docs/how-to/export-data-bigquery#setup).


---


## Using `mpdev` to run the SaaS Test Tool


### Before you begin
To use `mpdev` to run the SaaS Test Tool, you must install both `Docker` and `mpdev`.

1. To install Docker for Linux, follow the [Docker installation guide](https://docs.docker.com/engine/install/).
1. To install mpdev, follow the installation instructions for [marketplace-tools](https://github.com/GoogleCloudPlatform/marketplace-tools#installation).


### Create your local sandbox
> **🛈 NOTE:** These instructions have only been tested on Linux.


After you install all the required dependencies, you must create a working directory that you use to store service account credentials and configuration files.

```bash
mkdir -p ~/saas-test-tool-sandbox
cd ~/saas-test-tool-sandbox

# Pull a sample SaaS Test Tool configuration and place it in your sandbox directory.
# Include the dot (.) in the end of the command.
mpdev pkg get https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/saas-template/partner-integration .

cd partner-integration
```


### Download your service account token
The SaaS Test Tool must authenticate with Google Cloud to create/delete orders and query usage data. You can use a JSON Web Token (JWT) to do this.

> **🛈 NOTE:** You must use the same project you created previously.


1. Go to the [Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) page.
1. Select the service account you created previously.
1. Go to the **KEYS** tab.
1. Select ADD **KEY > Create new key**.
1. Select the **JSON** key type.
1. Click **CREATE**.
1. Save the credential file as `cred.json`.
1. Copy the credential file to the `~/saas-test-tool-sandbox/partner-integration` sandbox directory you created previously.


### Update your SaaS Test Tool configuration file
The configuration file contains the following information:
* The SaaS product that you want to test.
* The plans you want to test orders/upgrades on.
* The 100% discounted test billing account to use.
* A docker run command that runs a driver.
* Expected usage billing after running the driver.


The driver is any Docker image that emulates usage of the SaaS product that you’re testing. That usage is reported to Google Cloud, and the SaaS Test Tool checks to see if the usage is within the expected range.


#### Here’s an example of a working configuration file:
```yaml
apiVersion: dev.marketplace.cloud.google.com/v1alpha1
kind: SaasListingTemplate
credFilePath: cred.json
integrationTestConfig:
  provider: "providers/e2e-testing"
  productExternalName: "procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog"
  billingAccount: "billingAccounts/017778-0B6CC1-FB92E9"
  plans: ["plan-a", "plan-d"]
  approveEntitlementTimeoutSeconds: 600
  approvePlanChangeTimeoutSeconds: 600
  billingMeteringTestConfig:
    - driver:
        driverCommand: "docker run -v /var/run/docker.sock:/var/run/docker.sock -v ~/saas-test-tool-sandbox/partner-integration:/keys -e GOOGLE_APPLICATION_CREDENTIALS=/keys/cred.json -e GOOGLE_CLOUD_PROJECT=testing-producer-saas-322600 gcr.io/marketplace-saas-tools/billing-metering-driver --service=procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog --metric=procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog/plan_a_metric_1 --consumer_id=project_number:{@USAGE_REPORTING_ID}"
        planId: "plan-a"
        connectionInfo:
          project: "projects/testing-producer-saas-322600"
          tableName: "testing-producer-saas-322600.testing_producer_saas_billing_export.gcp_billing_export_resource_v1_017778_0B6CC1_FB92E9"
        expectation:
          skuId: "0D8F-CA36-DD7A"
          usageExpectation:
            min: 150
            max: 200
            baseUnits: "requests"
          costExpectation:
            min: 800
            max: 1000
            currency: "USD"
```


### Run the SaaS Test Tool test
> **🛈 NOTE:** It takes approximately 8 hours for the billing export to complete, so it takes approximately 8 hours for the SaaS Test Tool to test billing usage.


> **🛈 NOTE:** Logs are written to `output.log` in the same directory as the configuration file

```bash
cd ~/saas-test-tool-sandbox/partner-integration

mpdev test -f config.yaml 2>&1 | tee output.log
```


---


## Troubleshooting known error messages


### Requested entity already exists
If you encounter the following error message, please run the SaaS Test Tool again.


This message states that an order already exists (bad state). This usually happens when the SaaS Test Tool doesn’t correctly finish running. When you run the SaaS Test Tool again, the state is reset.


```
Caused by: com.google.api.client.googleapis.json.GoogleJsonResponseException: 409 Conflict
POST https://cloudcommerceconsumerprocurement.googleapis.com/v1alpha1/billingAccounts/010079-E1D1FB-A9454A/orders:place
{
  "code": 409,
  "errors": [
    {
      "domain": "global",
      "message": "Requested entity already exists",
      "reason": "alreadyExists"
    }
  ],
  "message": "Requested entity already exists",
  "status": "ALREADY_EXISTS"
}
```


### Billing Account Permission Error
This error means that the service account hasn’t been granted sufficient permissions on the billing account. To resolve this error, refer to [Grant Billing Administrator permissions to your service account](#Grant-Billing-Administrator-permissions-to-your-service-account).

```
Caused by: com.google.cloud.marketplace.saas.partnerintegration.test.IntegrationTestException: There's a permission issue. Please make sure the service account you're using has billing account admin permissions on your billing account, and provider admin permissions on the provider.
        at com.google.cloud.marketplace.saas.partnerintegration.test.TestUtil.castHttpException(TestUtil.java:136)
        at com.google.cloud.marketplace.saas.partnerintegration.test.PlaceOrderIntegrationTest.execute(PlaceOrderIntegrationTest.java:58)
        at com.google.cloud.marketplace.saas.partnerintegration.testrunner.TestRunner.runPlaceOrderTest(TestRunner.java:137)
        ... 10 more
Caused by: com.google.api.client.googleapis.json.GoogleJsonResponseException: 403 Forbidden
POST https://cloudcommerceprocurement.googleapis.com/v1/providers/e2e-testing/accounts/E-B8D8-0415-157D-E362:reset
{
  "code": 403,
  "errors": [
    {
      "domain": "global",
      "message": "The caller does not have permission",
      "reason": "forbidden"
    }
  ],
  "message": "The caller does not have permission",
  "status": "PERMISSION_DENIED"
}
```


### User does not have BigQuery Permission
This error occurs when the service account user hasn’t been granted permission to create BigQuery jobs. To resolve this issue, refer to [Grant your service account access to Billing metering results](#Grant-your-service-account-access-to-Billing-metering-results).

```
Caused by: com.google.api.client.googleapis.json.GoogleJsonResponseException: 403 Forbidden
POST https://www.googleapis.com/bigquery/v2/projects/saas-bb-tf-project-01/queries
{
  "code" : 403,
  "errors" : [ {
    "domain" : "global",
    "message" : "Access Denied: Project saas-bb-tf-project-01: User does not have bigquery.jobs.create permission in project saas-bb-tf-project-01.",
    "reason" : "accessDenied"
  } ],
  "message" : "Access Denied: Project saas-bb-tf-project-01: User does not have bigquery.jobs.create permission in project saas-bb-tf-project-01.",
  "status" : "PERMISSION_DENIED"
}
```


### Caller does not have permissions
This error occurs when the service account hasn’t been granted sufficient permissions to call the Partner Procurement API. To resolve this issue, refer to [Link your service account to call the Procurement API and report usage](#Link-your-service-account-to-call-the-Procurement-API-and-report-usage).

```
Caused by: com.google.cloud.marketplace.saas.partnerintegration.test.IntegrationTestException: There's a permission issue. Please make sure the service account you're using has billing account admin permissions on your billing account, and provider admin permissions on the provider.
        at com.google.cloud.marketplace.saas.partnerintegration.test.TestUtil.castHttpException(TestUtil.java:136)
        at com.google.cloud.marketplace.saas.partnerintegration.test.PlaceOrderIntegrationTest.execute(PlaceOrderIntegrationTest.java:58)
        at com.google.cloud.marketplace.saas.partnerintegration.testrunner.TestRunner.runPlaceOrderTest(TestRunner.java:137)
        ... 10 more
Caused by: com.google.api.client.googleapis.json.GoogleJsonResponseException: 403 Forbidden
POST https://cloudcommerceprocurement.googleapis.com/v1/providers/e2e-testing/accounts/E-B8D8-0415-157D-E362:approve
{
  "code": 403,
  "errors": [
    {
      "domain": "global",
      "message": "The caller does not have permission",
      "reason": "forbidden"
    }
  ],
  "message": "The caller does not have permission",
  "status": "PERMISSION_DENIED"
}
```


### Account must be in approved state
This error occurs when the procurement account hasn’t been approved by the partner. To resolve this issue, the partner must approve the procurement account request. For more information, refer to [Approve the procurement account](#Approve-the-procurement-account).

```
Exception in thread "main" java.lang.IllegalStateException: Account must be in approved state in order to proceed with test.
        at com.google.cloud.marketplace.saas.partnerintegration.client.ConsumerProcurementServiceHelperImpl.getAccountByProvider(ConsumerProcurementServiceHelperImpl.java:164)
        at com.google.cloud.marketplace.saas.partnerintegration.testrunner.TestRunner.getProcurementAccountName(TestRunner.java:114)
        at com.google.cloud.marketplace.saas.partnerintegration.testrunner.TestRunner.execute(TestRunner.java:70)
        at com.google.cloud.marketplace.saas.partnerintegration.Main.main(Main.java:30)
Error: Error in resource {Group:dev.marketplace.cloud.google.com Kind:SaasListingTemplate Name:}
```


### servicemanagement.services.check denied for the consumer project

This error occurs when the billing metering driver fails to report usage due to insufficient permissions. To resolve this error, refer to [Link your service account to call the Procurement API and report usage](#Link-your-service-account-to-call-the-Procurement-API-and-report-usage).

```
googleapiclient.errors.HttpError: <HttpError 403 when requesting https://servicecontrol.googleapis.com/v1/services/procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog:check?alt=json returned "Permission 'servicemanagement.services.check' denied for the consumer project (or it may not exist)". Details: "Permission 'servicemanagement.services.check' denied for the consumer project (or it may not exist)">
```


### Precondition check failed

This error occurs when the billing account isn’t configured with a 100% discount. To resolve this issue, refer to [Configure a 100% Discount for your billing account](#Configure-a-100-Discount-for-your-billing-account).

```
Exception in thread "main" com.google.cloud.marketplace.saas.partnerintegration.testrunner.ProductLifecycleTestException: Integration test failed to complete.
        at com.google.cloud.marketplace.saas.partnerintegration.testrunner.TestRunner.runPlaceOrderTest(TestRunner.java:139)
        at java.base/java.util.stream.ReferencePipeline$3$1.accept(ReferencePipeline.java:197)
        at java.base/java.util.Spliterators$ArraySpliterator.forEachRemaining(Spliterators.java:992)
        at java.base/java.util.stream.AbstractPipeline.copyInto(AbstractPipeline.java:509)
        at java.base/java.util.stream.AbstractPipeline.wrapAndCopyInto(AbstractPipeline.java:499)
        at java.base/java.util.stream.ReduceOps$ReduceOp.evaluateSequential(ReduceOps.java:921)
        at java.base/java.util.stream.AbstractPipeline.evaluate(AbstractPipeline.java:234)
        at java.base/java.util.stream.ReferencePipeline.collect(ReferencePipeline.java:682)
        at com.google.cloud.marketplace.saas.partnerintegration.testrunner.TestRunner.runPlaceOrderTests(TestRunner.java:132)
        at com.google.cloud.marketplace.saas.partnerintegration.testrunner.TestRunner.execute(TestRunner.java:75)
        at com.google.cloud.marketplace.saas.partnerintegration.Main.main(Main.java:30)
Caused by: com.google.cloud.marketplace.saas.partnerintegration.test.IntegrationTestException: Some parameters of the test are incorrect. Please make sure the billing account name, product external name, and plan names are correct, then retry.
        at com.google.cloud.marketplace.saas.partnerintegration.test.TestUtil.castHttpException(TestUtil.java:131)
        at com.google.cloud.marketplace.saas.partnerintegration.test.PlaceOrderIntegrationTest.execute(PlaceOrderIntegrationTest.java:58)
        at com.google.cloud.marketplace.saas.partnerintegration.testrunner.TestRunner.runPlaceOrderTest(TestRunner.java:137)
        ... 10 more
Caused by: com.google.api.client.googleapis.json.GoogleJsonResponseException: 400 Bad Request
POST https://cloudcommerceconsumerprocurement.googleapis.com/v1alpha1/billingAccounts/0128B6-4EBE44-E24ED4/orders:place
{
  "code": 400,
  "errors": [
    {
      "domain": "global",
      "message": "Precondition check failed.",
      "reason": "failedPrecondition"
    }
  ],
  "message": "Precondition check failed.",
  "status": "FAILED_PRECONDITION"
}
```
