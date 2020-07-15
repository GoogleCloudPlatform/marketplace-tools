# Autogen Reference
<a name="top"></a>

## Table of Contents

- [DeploymentPackageAutogenSpec](#DeploymentPackageAutogenSpec)
    - [AcceleratorSpec](#cloud.deploymentmanager.autogen.AcceleratorSpec)
    - [ApplicationStatusSpec](#cloud.deploymentmanager.autogen.ApplicationStatusSpec)
    - [ApplicationStatusSpec.WaiterSpec](#cloud.deploymentmanager.autogen.ApplicationStatusSpec.WaiterSpec)
    - [ApplicationStatusSpec.WaiterSpec.ScriptSpec](#cloud.deploymentmanager.autogen.ApplicationStatusSpec.WaiterSpec.ScriptSpec)
    - [BooleanExpression](#cloud.deploymentmanager.autogen.BooleanExpression)
    - [BooleanExpression.BooleanDeployInputField](#cloud.deploymentmanager.autogen.BooleanExpression.BooleanDeployInputField)
    - [BooleanExpression.ExternalIpAvailability](#cloud.deploymentmanager.autogen.BooleanExpression.ExternalIpAvailability)
    - [DeployInputField](#cloud.deploymentmanager.autogen.DeployInputField)
    - [DeployInputField.BooleanCheckbox](#cloud.deploymentmanager.autogen.DeployInputField.BooleanCheckbox)
    - [DeployInputField.EmailBox](#cloud.deploymentmanager.autogen.DeployInputField.EmailBox)
    - [DeployInputField.EmailBox.Validation](#cloud.deploymentmanager.autogen.DeployInputField.EmailBox.Validation)
    - [DeployInputField.GceZoneDropdown](#cloud.deploymentmanager.autogen.DeployInputField.GceZoneDropdown)
    - [DeployInputField.GroupedBooleanCheckbox](#cloud.deploymentmanager.autogen.DeployInputField.GroupedBooleanCheckbox)
    - [DeployInputField.GroupedBooleanCheckbox.DisplayGroup](#cloud.deploymentmanager.autogen.DeployInputField.GroupedBooleanCheckbox.DisplayGroup)
    - [DeployInputField.IntegerBox](#cloud.deploymentmanager.autogen.DeployInputField.IntegerBox)
    - [DeployInputField.IntegerBox.Validation](#cloud.deploymentmanager.autogen.DeployInputField.IntegerBox.Validation)
    - [DeployInputField.IntegerDropdown](#cloud.deploymentmanager.autogen.DeployInputField.IntegerDropdown)
    - [DeployInputField.IntegerDropdown.ValueLabelsEntry](#cloud.deploymentmanager.autogen.DeployInputField.IntegerDropdown.ValueLabelsEntry)
    - [DeployInputField.StringBox](#cloud.deploymentmanager.autogen.DeployInputField.StringBox)
    - [DeployInputField.StringBox.Validation](#cloud.deploymentmanager.autogen.DeployInputField.StringBox.Validation)
    - [DeployInputField.StringDropdown](#cloud.deploymentmanager.autogen.DeployInputField.StringDropdown)
    - [DeployInputField.StringDropdown.ValueLabelsEntry](#cloud.deploymentmanager.autogen.DeployInputField.StringDropdown.ValueLabelsEntry)
    - [DeployInputSection](#cloud.deploymentmanager.autogen.DeployInputSection)
    - [DeployInputSpec](#cloud.deploymentmanager.autogen.DeployInputSpec)
    - [DeploymentPackageAutogenSpec](#cloud.deploymentmanager.autogen.DeploymentPackageAutogenSpec)
    - [DiskSpec](#cloud.deploymentmanager.autogen.DiskSpec)
    - [DiskSpec.DeviceName](#cloud.deploymentmanager.autogen.DiskSpec.DeviceName)
    - [DiskSpec.DiskSize](#cloud.deploymentmanager.autogen.DiskSpec.DiskSize)
    - [DiskSpec.DiskType](#cloud.deploymentmanager.autogen.DiskSpec.DiskType)
    - [ExternalIpSpec](#cloud.deploymentmanager.autogen.ExternalIpSpec)
    - [FirewallRuleSpec](#cloud.deploymentmanager.autogen.FirewallRuleSpec)
    - [GceMetadataItem](#cloud.deploymentmanager.autogen.GceMetadataItem)
    - [GceMetadataItem.TierVmNames](#cloud.deploymentmanager.autogen.GceMetadataItem.TierVmNames)
    - [GceMetadataItem.TierVmNames.AllVmList](#cloud.deploymentmanager.autogen.GceMetadataItem.TierVmNames.AllVmList)
    - [GceStartupScriptSpec](#cloud.deploymentmanager.autogen.GceStartupScriptSpec)
    - [GcpAuthScopeSpec](#cloud.deploymentmanager.autogen.GcpAuthScopeSpec)
    - [ImageSpec](#cloud.deploymentmanager.autogen.ImageSpec)
    - [InstanceUrlSpec](#cloud.deploymentmanager.autogen.InstanceUrlSpec)
    - [Int32List](#cloud.deploymentmanager.autogen.Int32List)
    - [Int32Range](#cloud.deploymentmanager.autogen.Int32Range)
    - [IpForwardingSpec](#cloud.deploymentmanager.autogen.IpForwardingSpec)
    - [LocalSsdSpec](#cloud.deploymentmanager.autogen.LocalSsdSpec)
    - [MachineTypeSpec](#cloud.deploymentmanager.autogen.MachineTypeSpec)
    - [MachineTypeSpec.MachineType](#cloud.deploymentmanager.autogen.MachineTypeSpec.MachineType)
    - [MachineTypeSpec.MachineTypeConstraint](#cloud.deploymentmanager.autogen.MachineTypeSpec.MachineTypeConstraint)
    - [MultiVmDeploymentPackageSpec](#cloud.deploymentmanager.autogen.MultiVmDeploymentPackageSpec)
    - [NetworkInterfacesSpec](#cloud.deploymentmanager.autogen.NetworkInterfacesSpec)
    - [OptionalInt32](#cloud.deploymentmanager.autogen.OptionalInt32)
    - [OptionalString](#cloud.deploymentmanager.autogen.OptionalString)
    - [PasswordSpec](#cloud.deploymentmanager.autogen.PasswordSpec)
    - [PostDeployInfo](#cloud.deploymentmanager.autogen.PostDeployInfo)
    - [PostDeployInfo.ActionItem](#cloud.deploymentmanager.autogen.PostDeployInfo.ActionItem)
    - [PostDeployInfo.ConnectToInstanceSpec](#cloud.deploymentmanager.autogen.PostDeployInfo.ConnectToInstanceSpec)
    - [PostDeployInfo.InfoRow](#cloud.deploymentmanager.autogen.PostDeployInfo.InfoRow)
    - [SingleVmDeploymentPackageSpec](#cloud.deploymentmanager.autogen.SingleVmDeploymentPackageSpec)
    - [StackdriverSpec](#cloud.deploymentmanager.autogen.StackdriverSpec)
    - [StackdriverSpec.Logging](#cloud.deploymentmanager.autogen.StackdriverSpec.Logging)
    - [StackdriverSpec.Monitoring](#cloud.deploymentmanager.autogen.StackdriverSpec.Monitoring)
    - [TierVmInstance](#cloud.deploymentmanager.autogen.TierVmInstance)
    - [VmTierSpec](#cloud.deploymentmanager.autogen.VmTierSpec)
    - [VmTierSpec.TierInstanceCount](#cloud.deploymentmanager.autogen.VmTierSpec.TierInstanceCount)
    - [ZoneSpec](#cloud.deploymentmanager.autogen.ZoneSpec)
  
    - [ApplicationStatusSpec.StatusType](#cloud.deploymentmanager.autogen.ApplicationStatusSpec.StatusType)
    - [DeployInputSection.Placement](#cloud.deploymentmanager.autogen.DeployInputSection.Placement)
    - [ExternalIpSpec.Type](#cloud.deploymentmanager.autogen.ExternalIpSpec.Type)
    - [FirewallRuleSpec.Protocol](#cloud.deploymentmanager.autogen.FirewallRuleSpec.Protocol)
    - [FirewallRuleSpec.TrafficSource](#cloud.deploymentmanager.autogen.FirewallRuleSpec.TrafficSource)
    - [GcpAuthScopeSpec.Scope](#cloud.deploymentmanager.autogen.GcpAuthScopeSpec.Scope)
    - [InstanceUrlSpec.Scheme](#cloud.deploymentmanager.autogen.InstanceUrlSpec.Scheme)
  
- [Scalar Value Types](#scalar-value-types)



<a name="DeploymentPackageAutogenSpec"></a>
<p align="right"><a href="#top">Top</a></p>

## DeploymentPackageAutogenSpec



<a name="cloud.deploymentmanager.autogen.AcceleratorSpec"></a>

### AcceleratorSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| types | [string](#string) | repeated | One or more accelerator types. If the list contains exactly one item, the type will not be selectable by the user. Currently avaiable types: &#34;nvidia-tesla-k80&#34;, &#34;nvidia-tesla-p100&#34;, &#34;nvidia-tesla-v100&#34;. |
| default_type | [string](#string) |  | Default accelerator type, which should be one of those listed above |
| default_count | [int32](#int32) |  | Default number of accelerators. Currently, only values of 0, 1, 2, 4, and 8 are supported. |
| min_count | [int32](#int32) |  | Minimum number of accelerators (inclusive) that may be selected. Currently, only values of 0, 1, 2, 4, and 8 are supported. |
| max_count | [int32](#int32) |  | Maximum number of accelerators (inclusive) that may be selected. Currently, only values of 0, 1, 2, 4, and 8 are supported. |






<a name="cloud.deploymentmanager.autogen.ApplicationStatusSpec"></a>

### ApplicationStatusSpec
Specifies how to monitor application installation status in order to
detect when the application is ready or has failed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ApplicationStatusSpec.StatusType](#cloud.deploymentmanager.autogen.ApplicationStatusSpec.StatusType) |  | Defines how to monitor application installation status. Required. |
| waiter | [ApplicationStatusSpec.WaiterSpec](#cloud.deploymentmanager.autogen.ApplicationStatusSpec.WaiterSpec) |  | Required if type is WAITER. |






<a name="cloud.deploymentmanager.autogen.ApplicationStatusSpec.WaiterSpec"></a>

### ApplicationStatusSpec.WaiterSpec
Specifies how the waiter is setup.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| waiter_timeout_secs | [int32](#int32) |  | Timeout when the waiter fails itself in absence of status signals. Required. |
| script | [ApplicationStatusSpec.WaiterSpec.ScriptSpec](#cloud.deploymentmanager.autogen.ApplicationStatusSpec.WaiterSpec.ScriptSpec) |  | Optional integration with the VM to signal the waiter via startup-script. If the script spec is not present, the application is expected to directly signal the waiter. |






<a name="cloud.deploymentmanager.autogen.ApplicationStatusSpec.WaiterSpec.ScriptSpec"></a>

### ApplicationStatusSpec.WaiterSpec.ScriptSpec
Specifies the integration with the VM to signal the waiter via
startup-script.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| check_timeout_secs | [int32](#int32) |  | If not set, the waiter timeout will be used. |
| check_script_content | [string](#string) |  | Optional bash script to check the status. This should return 0 if the application is ready, 1 if the application is not yet ready but the check should be retried, and greater than 1 if the check has failed permanently. If the script is not present, the waiter is signaled as soon as the VM finishes booting. |
| disable_startup_script_url | [bool](#bool) |  | If true, the generated template will include an empty &#34;startup-script-url&#34; VM metadata value. This effectively disables project-wide startup_script_url settings which took precedence over instance-level startup_script settings in older versions of the Google instance init logic. This option is not necessary for images that use Google base packages newer than June, 2016. See b/31729022 for more context. TODO(volkman): Remove this once all images (Brocade, etc.) have moved to the new base package version. |






<a name="cloud.deploymentmanager.autogen.BooleanExpression"></a>

### BooleanExpression
Allows to build an expression which value should be evaluated to boolean.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| has_external_ip | [BooleanExpression.ExternalIpAvailability](#cloud.deploymentmanager.autogen.BooleanExpression.ExternalIpAvailability) |  |  |
| boolean_deploy_input_field | [BooleanExpression.BooleanDeployInputField](#cloud.deploymentmanager.autogen.BooleanExpression.BooleanDeployInputField) |  |  |






<a name="cloud.deploymentmanager.autogen.BooleanExpression.BooleanDeployInputField"></a>

### BooleanExpression.BooleanDeployInputField
Uses the value of a deploy input field of type boolean.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the deploy input field. Required. |
| negated | [bool](#bool) |  | If true, negate the value of the input field. |






<a name="cloud.deploymentmanager.autogen.BooleanExpression.ExternalIpAvailability"></a>

### BooleanExpression.ExternalIpAvailability
Allows to specify a condition based on external ip configuration for
a single instance (for single vm) or all instances in a tier (multi vm).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| negated | [bool](#bool) |  | Specifies if expression is based on external IP being available or not. |
| tier | [string](#string) |  | Multi-vm&#39;s tier name. It is required for multi vm spec. |






<a name="cloud.deploymentmanager.autogen.DeployInputField"></a>

### DeployInputField



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the field. Letters, numbers, dashes, and underscores only. Required. |
| title | [string](#string) |  | Title of the field. Required. |
| description | [string](#string) |  | Optional description appearing below the title. |
| tooltip | [string](#string) |  | Optional tooltip. |
| level | [int32](#int32) |  | Indicates the &#34;advanced&#34; level of the input property. Level 0 (default) will always be shown. Level 1 corresponds to one expansion (user clicks &#34;show advanced options&#34; or &#34;more options&#34;). Higher levels correspond to further expansions, or they may be collapsed to level 1 by the UI implementation. Optional. |
| required | [bool](#bool) |  | If required, the user must input a valid value. |
| boolean_checkbox | [DeployInputField.BooleanCheckbox](#cloud.deploymentmanager.autogen.DeployInputField.BooleanCheckbox) |  |  |
| grouped_boolean_checkbox | [DeployInputField.GroupedBooleanCheckbox](#cloud.deploymentmanager.autogen.DeployInputField.GroupedBooleanCheckbox) |  |  |
| integer_box | [DeployInputField.IntegerBox](#cloud.deploymentmanager.autogen.DeployInputField.IntegerBox) |  |  |
| integer_dropdown | [DeployInputField.IntegerDropdown](#cloud.deploymentmanager.autogen.DeployInputField.IntegerDropdown) |  |  |
| string_box | [DeployInputField.StringBox](#cloud.deploymentmanager.autogen.DeployInputField.StringBox) |  |  |
| string_dropdown | [DeployInputField.StringDropdown](#cloud.deploymentmanager.autogen.DeployInputField.StringDropdown) |  |  |
| zone_dropdown | [DeployInputField.GceZoneDropdown](#cloud.deploymentmanager.autogen.DeployInputField.GceZoneDropdown) |  |  |
| email_box | [DeployInputField.EmailBox](#cloud.deploymentmanager.autogen.DeployInputField.EmailBox) |  |  |






<a name="cloud.deploymentmanager.autogen.DeployInputField.BooleanCheckbox"></a>

### DeployInputField.BooleanCheckbox
A checkbox for a boolean value.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [bool](#bool) |  |  |






<a name="cloud.deploymentmanager.autogen.DeployInputField.EmailBox"></a>

### DeployInputField.EmailBox
A specialized textbox for email addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [string](#string) |  |  |
| validation | [DeployInputField.EmailBox.Validation](#cloud.deploymentmanager.autogen.DeployInputField.EmailBox.Validation) |  |  |
| placeholder | [string](#string) |  | A placeholder to hint the user what to enter here. If not specified, user@example.com is used. |
| test_default_value | [string](#string) |  | This attribute is used as field&#39;s value in automated tests. Defaults to default-user@example.com if this field is required unless the default_value or this field is explicitly present. If present, it overrides the default_value. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.EmailBox.Validation"></a>

### DeployInputField.EmailBox.Validation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| description | [string](#string) |  | Description shown to the user if the input value fails validation. If not specified, a default message is used. |
| regex | [string](#string) |  | Optional pattern. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.GceZoneDropdown"></a>

### DeployInputField.GceZoneDropdown
A dropdown with GCE zones as values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [OptionalString](#cloud.deploymentmanager.autogen.OptionalString) |  |  |






<a name="cloud.deploymentmanager.autogen.DeployInputField.GroupedBooleanCheckbox"></a>

### DeployInputField.GroupedBooleanCheckbox
A checkbox displayed next to other checkboxes under a common group title.
The group is intended for display purposes only. This is not a radio group;
the checkboxes are still independently selectable by the user.

The first checkbox in the group should define the display_group.
The immediately following GroupedBooleanCheckboxes without a display_group
are part of such group. In other words, the group ends when either
a different field type or one with a display_group is encountered.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [bool](#bool) |  |  |
| display_group | [DeployInputField.GroupedBooleanCheckbox.DisplayGroup](#cloud.deploymentmanager.autogen.DeployInputField.GroupedBooleanCheckbox.DisplayGroup) |  |  |






<a name="cloud.deploymentmanager.autogen.DeployInputField.GroupedBooleanCheckbox.DisplayGroup"></a>

### DeployInputField.GroupedBooleanCheckbox.DisplayGroup



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the group. Required. Convention is to use an UPPERCASE_UNDERSCORE name. |
| title | [string](#string) |  | Title of the group. Required. |
| description | [string](#string) |  | Optional description appearing below the title. |
| tooltip | [string](#string) |  | Optional tooltip. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.IntegerBox"></a>

### DeployInputField.IntegerBox
A textbox for entering an integer, with optional validation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [OptionalInt32](#cloud.deploymentmanager.autogen.OptionalInt32) |  |  |
| validation | [DeployInputField.IntegerBox.Validation](#cloud.deploymentmanager.autogen.DeployInputField.IntegerBox.Validation) |  |  |
| placeholder | [string](#string) |  |  |
| test_default_value | [OptionalInt32](#cloud.deploymentmanager.autogen.OptionalInt32) |  | This attribute is used as field&#39;s value in automated tests. If present, it overrides the default_value. For required fields without default_value, it is required to set test_default_value. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.IntegerBox.Validation"></a>

### DeployInputField.IntegerBox.Validation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| description | [string](#string) |  | Description shown to the user if the input value fails validation. |
| min | [OptionalInt32](#cloud.deploymentmanager.autogen.OptionalInt32) |  | Optional inclusive minimum value. |
| max | [OptionalInt32](#cloud.deploymentmanager.autogen.OptionalInt32) |  | Optional inclusive maximum value. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.IntegerDropdown"></a>

### DeployInputField.IntegerDropdown
A dropdown with integer values, with optional labels.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [int32](#int32) | repeated |  |
| default_value_index | [OptionalInt32](#cloud.deploymentmanager.autogen.OptionalInt32) |  |  |
| value_labels | [DeployInputField.IntegerDropdown.ValueLabelsEntry](#cloud.deploymentmanager.autogen.DeployInputField.IntegerDropdown.ValueLabelsEntry) | repeated | Optional labels for values (not indices). If a value does not have a corresponding label, its numeric string is used. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.IntegerDropdown.ValueLabelsEntry"></a>

### DeployInputField.IntegerDropdown.ValueLabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [int32](#int32) |  |  |
| value | [string](#string) |  |  |






<a name="cloud.deploymentmanager.autogen.DeployInputField.StringBox"></a>

### DeployInputField.StringBox
A textbox for entering a string, with optional validation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [string](#string) |  |  |
| validation | [DeployInputField.StringBox.Validation](#cloud.deploymentmanager.autogen.DeployInputField.StringBox.Validation) |  |  |
| placeholder | [string](#string) |  |  |
| test_default_value | [string](#string) |  | This attribute is used as field&#39;s value in automated tests. If present, it overrides the default_value. For required fields without default_value, it is required to set test_default_value. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.StringBox.Validation"></a>

### DeployInputField.StringBox.Validation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| description | [string](#string) |  | Description shown to the user if the input value fails validation. |
| regex | [string](#string) |  | Optional pattern. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.StringDropdown"></a>

### DeployInputField.StringDropdown
A dropdown with string values, with optional labels.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [string](#string) | repeated |  |
| default_value_index | [OptionalInt32](#cloud.deploymentmanager.autogen.OptionalInt32) |  |  |
| value_labels | [DeployInputField.StringDropdown.ValueLabelsEntry](#cloud.deploymentmanager.autogen.DeployInputField.StringDropdown.ValueLabelsEntry) | repeated | Optional labels for values. If a value does not have a corresponding label, its raw value will be displayed. |






<a name="cloud.deploymentmanager.autogen.DeployInputField.StringDropdown.ValueLabelsEntry"></a>

### DeployInputField.StringDropdown.ValueLabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="cloud.deploymentmanager.autogen.DeployInputSection"></a>

### DeployInputSection



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| placement | [DeployInputSection.Placement](#cloud.deploymentmanager.autogen.DeployInputSection.Placement) |  | Required. |
| name | [string](#string) |  | Section name, required if this section is a custom one. Must be unique among all sections. Convention is to use an UPPERCASE_UNDERSCORE name. |
| tier | [string](#string) |  | For Placement.TIER, this specifies the required tier name. |
| title | [string](#string) |  | Section title, required if this section is a custom one. |
| description | [string](#string) |  | Optional description appearing below the title. |
| tooltip | [string](#string) |  | Optional tooltip. |
| fields | [DeployInputField](#cloud.deploymentmanager.autogen.DeployInputField) | repeated |  |






<a name="cloud.deploymentmanager.autogen.DeployInputSpec"></a>

### DeployInputSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sections | [DeployInputSection](#cloud.deploymentmanager.autogen.DeployInputSection) | repeated | One or more sections containing input fields. |






<a name="cloud.deploymentmanager.autogen.DeploymentPackageAutogenSpec"></a>

### DeploymentPackageAutogenSpec
Top level spec.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| single_vm | [SingleVmDeploymentPackageSpec](#cloud.deploymentmanager.autogen.SingleVmDeploymentPackageSpec) |  | For solution deploying a single VM. |
| multi_vm | [MultiVmDeploymentPackageSpec](#cloud.deploymentmanager.autogen.MultiVmDeploymentPackageSpec) |  | For solution deploying multiple tiers of VMs. |
| version | [string](#string) |  | Human readable version of the deployment package. |






<a name="cloud.deploymentmanager.autogen.DiskSpec"></a>

### DiskSpec
Specifies a persistent disk.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disk_size | [DiskSpec.DiskSize](#cloud.deploymentmanager.autogen.DiskSpec.DiskSize) |  |  |
| disk_type | [DiskSpec.DiskType](#cloud.deploymentmanager.autogen.DiskSpec.DiskType) |  |  |
| display_label | [string](#string) |  | A short descriptive label for this disk. Optional for boot disk; default is &#39;Boot disk&#39;. Optional if this is the only one additional disk; default is &#39;Data disk&#39;. Required otherwise. |
| device_name_suffix | [DiskSpec.DeviceName](#cloud.deploymentmanager.autogen.DiskSpec.DeviceName) |  | Specifies the device name suffix. Ignored for boot disk. Optional. The final device name will be a concatenation of an instance name with the specified device name. |
| prevent_auto_deletion | [bool](#bool) |  | Whether to leave the disk when the instance is deleted. Ignored for boot disk. |






<a name="cloud.deploymentmanager.autogen.DiskSpec.DeviceName"></a>

### DiskSpec.DeviceName



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Device name specified as a constant string. Optional. By default takes the value of disk&#39;s name. |
| name_from_deploy_input_field | [string](#string) |  | Specifies a deploy input field name from which the device name should be read. |






<a name="cloud.deploymentmanager.autogen.DiskSpec.DiskSize"></a>

### DiskSpec.DiskSize



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_size_gb | [int32](#int32) |  | The default disk size in GB. Required. |
| min_size_gb | [int32](#int32) |  | Specifies the min disk size allowed in GB. |
| max_size_gb | [int32](#int32) |  | Specifies the max disk size allowed in GB. Not supported yet. |
| not_configurable | [bool](#bool) |  | Whether to restrict the user from changing away from the default settings. Not supported yet (always false). |






<a name="cloud.deploymentmanager.autogen.DiskSpec.DiskType"></a>

### DiskSpec.DiskType



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_type | [string](#string) |  | The default disk type. Required. See http://cloud.google.com/compute/docs/reference/latest/diskTypes. |
| not_configurable | [bool](#bool) |  | Whether to restrict the user from changing away from the default settings. Not supported yet (always false). |






<a name="cloud.deploymentmanager.autogen.ExternalIpSpec"></a>

### ExternalIpSpec
Defines how a VM is exposed on the Internet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_type | [ExternalIpSpec.Type](#cloud.deploymentmanager.autogen.ExternalIpSpec.Type) |  | Required. |
| not_configurable | [bool](#bool) |  | Whether to restrict the user from changing away from the default settings. |






<a name="cloud.deploymentmanager.autogen.FirewallRuleSpec"></a>

### FirewallRuleSpec
Specifies a firewall rule.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| protocol | [FirewallRuleSpec.Protocol](#cloud.deploymentmanager.autogen.FirewallRuleSpec.Protocol) |  | The IP Protocol that the firewall rule allows. Required. |
| port | [string](#string) |  | The target ports on the VM, which could be a single port number like &#34;80&#34; or a port range like &#34;32768-40000&#34;. |
| default_off | [bool](#bool) |  | Specifies that by default it should be off. Applicable to TrafficSource.PUBLIC only. |
| not_configurable | [bool](#bool) |  | Whether to restrict the user from changing away from the default settings. Not supported yet (always false). Applicable to TrafficSource.PUBLIC only. |
| allowed_source | [FirewallRuleSpec.TrafficSource](#cloud.deploymentmanager.autogen.FirewallRuleSpec.TrafficSource) |  | Source of traffic, wrapping range/tags with friendly enum. Defaults to TrafficSource.PUBLIC. |






<a name="cloud.deploymentmanager.autogen.GceMetadataItem"></a>

### GceMetadataItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | Metadata item key. Required. |
| value | [string](#string) |  | Static metadata item value. |
| tier_vm_names | [GceMetadataItem.TierVmNames](#cloud.deploymentmanager.autogen.GceMetadataItem.TierVmNames) |  | The value is the names of one or more VMs in a tier. |
| value_from_deploy_input_field | [string](#string) |  | Value referenced from deploy input field. Should specify existing input field&#39;s name. |






<a name="cloud.deploymentmanager.autogen.GceMetadataItem.TierVmNames"></a>

### GceMetadataItem.TierVmNames



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tier | [string](#string) |  | The tier name. |
| vm_index | [int32](#int32) |  | 0-based index of a VM. A negative index can be used, with -1 referring the last, -2 the second last, etc. |
| all_vms | [GceMetadataItem.TierVmNames.AllVmList](#cloud.deploymentmanager.autogen.GceMetadataItem.TierVmNames.AllVmList) |  | All VM names as a string list. |






<a name="cloud.deploymentmanager.autogen.GceMetadataItem.TierVmNames.AllVmList"></a>

### GceMetadataItem.TierVmNames.AllVmList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delimiter | [string](#string) |  | Delimiter for the VM names in the list, for example a comma. Required. |






<a name="cloud.deploymentmanager.autogen.GceStartupScriptSpec"></a>

### GceStartupScriptSpec
Specifies the initial startup script for a VM instance.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bash_script_content | [string](#string) |  | Specifies a complete startup script. If waiter specifies its check script, those two will be combined with a software_status_script.py resource. Required. |






<a name="cloud.deploymentmanager.autogen.GcpAuthScopeSpec"></a>

### GcpAuthScopeSpec
Specifies access to an GCP API on the VM. This effectively
configures the corresponding scope under the VM&#39;s service account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| scope | [GcpAuthScopeSpec.Scope](#cloud.deploymentmanager.autogen.GcpAuthScopeSpec.Scope) |  | API scope. Required. |
| default_off | [bool](#bool) |  | Specifies that by default it should be off. |
| not_configurable | [bool](#bool) |  | Whether to restrict the user from changing away from the default settings. Not supported yet (always false). |






<a name="cloud.deploymentmanager.autogen.ImageSpec"></a>

### ImageSpec
Specifies a disk image resource.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [string](#string) |  | The GCP project containing the image. Required. |
| name | [string](#string) |  | The name of the image. Required. |
| label | [string](#string) |  | A descriptive label for this image, useful in a list of images for the user to select from. |






<a name="cloud.deploymentmanager.autogen.InstanceUrlSpec"></a>

### InstanceUrlSpec
Specifies a URL used for accessing the application on the VM.
The domain is implied as the VM instance address.
Currently the machine IP is used, but that might change to
another endpoint in the future.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| scheme | [InstanceUrlSpec.Scheme](#cloud.deploymentmanager.autogen.InstanceUrlSpec.Scheme) |  | The URL scheme. |
| port | [int32](#int32) |  | The URL port. |
| path | [string](#string) |  | The URL path, without the leading forward slash. |
| query | [string](#string) |  | The URL query, without the leading question mark. |
| fragment | [string](#string) |  | The URL fragment, without the leading hash sign. |
| tier_vm | [TierVmInstance](#cloud.deploymentmanager.autogen.TierVmInstance) |  | Specifies a VM from tier whose address would be used. Required in a multi-VM configuration. |






<a name="cloud.deploymentmanager.autogen.Int32List"></a>

### Int32List



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [int32](#int32) | repeated |  |






<a name="cloud.deploymentmanager.autogen.Int32Range"></a>

### Int32Range



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start_value | [int32](#int32) |  | The inclusive starting value. Required. |
| end_value | [int32](#int32) |  | The inclusive ending value. Required. |






<a name="cloud.deploymentmanager.autogen.IpForwardingSpec"></a>

### IpForwardingSpec
Specifies if the VM can route IP packets.
See http://cloud.google.com/compute/docs/networking#eventualconsistency.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_off | [bool](#bool) |  | Specifies that by default it should be off. |
| not_configurable | [bool](#bool) |  | Whether to restrict the user from changing away from the default settings. Not supported yet (always false). |






<a name="cloud.deploymentmanager.autogen.LocalSsdSpec"></a>

### LocalSsdSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| count | [int32](#int32) |  | Specifies the number of local SSDs to be attached to a vm instance. |
| count_from_deploy_input_field | [string](#string) |  | Specifies the number of local SSDs by referencing a value from a deploy input field. |






<a name="cloud.deploymentmanager.autogen.MachineTypeSpec"></a>

### MachineTypeSpec
Specifies the default machine type, and any size constraints which
restrict what the user can select.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_machine_type | [MachineTypeSpec.MachineType](#cloud.deploymentmanager.autogen.MachineTypeSpec.MachineType) |  | Specifies the machine type that should be selecteed by default. Required. |
| minimum | [MachineTypeSpec.MachineTypeConstraint](#cloud.deploymentmanager.autogen.MachineTypeSpec.MachineTypeConstraint) |  | Specifies the minimum requirement for a user-selected machine type. |
| maximum | [MachineTypeSpec.MachineTypeConstraint](#cloud.deploymentmanager.autogen.MachineTypeSpec.MachineTypeConstraint) |  | Specifies the minimum requirement for a user-selected machine type. Not supported yet (no max). |
| not_configurable | [bool](#bool) |  | Whether to restrict the user from changing away from the default settings. Not supported yet (always false). |






<a name="cloud.deploymentmanager.autogen.MachineTypeSpec.MachineType"></a>

### MachineTypeSpec.MachineType
Specifies a machine type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| gce_machine_type | [string](#string) |  | A predefined or custom machine type. Required. See http://cloud.google.com/compute/docs/machine-types. |






<a name="cloud.deploymentmanager.autogen.MachineTypeSpec.MachineTypeConstraint"></a>

### MachineTypeSpec.MachineTypeConstraint
Specifies an upper- or lower-bound constraint.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cpu | [int32](#int32) |  |  |
| ram_gb | [float](#float) |  |  |






<a name="cloud.deploymentmanager.autogen.MultiVmDeploymentPackageSpec"></a>

### MultiVmDeploymentPackageSpec
Specifies a solution that deploys Multiple VMs.
Next ID: 9


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tiers | [VmTierSpec](#cloud.deploymentmanager.autogen.VmTierSpec) | repeated | One or more tiers. |
| site_url | [InstanceUrlSpec](#cloud.deploymentmanager.autogen.InstanceUrlSpec) |  | Declares a URL to access the deployed application. |
| admin_url | [InstanceUrlSpec](#cloud.deploymentmanager.autogen.InstanceUrlSpec) |  | Declares a URL to administer the deployed application. |
| passwords | [PasswordSpec](#cloud.deploymentmanager.autogen.PasswordSpec) | repeated | Defines how to generate passwords at deployment time. |
| post_deploy | [PostDeployInfo](#cloud.deploymentmanager.autogen.PostDeployInfo) |  | Customizes post-deploy information displayed to the user. This helps get the user started with using the deployed solution. |
| deploy_input | [DeployInputSpec](#cloud.deploymentmanager.autogen.DeployInputSpec) |  | Customizes additional inputs configured by user prior to deployment. Currently, the configured values can be passed through to the VM via metadata items. |
| zone | [ZoneSpec](#cloud.deploymentmanager.autogen.ZoneSpec) |  | Customizes the zone selector. |
| stackdriver | [StackdriverSpec](#cloud.deploymentmanager.autogen.StackdriverSpec) |  | Integration with Stackdriver. |






<a name="cloud.deploymentmanager.autogen.NetworkInterfacesSpec"></a>

### NetworkInterfacesSpec
Network interfaces configuration for this solution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| min_count | [int32](#int32) |  | Minimum number of Network Interfaces (defaults to 1). |
| max_count | [int32](#int32) |  | Maximum number of Network Interfaces (can&#39;t exceed 8 and if not specified, will take the value of min_count). |
| external_ip | [ExternalIpSpec](#cloud.deploymentmanager.autogen.ExternalIpSpec) |  |  |
| labels | [string](#string) | repeated | Label that will be in front of each Network Interface (according to the index in this list). If the list is greater than min_count, the last label will be used to name all networks added beyond min_count. |






<a name="cloud.deploymentmanager.autogen.OptionalInt32"></a>

### OptionalInt32



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [int32](#int32) |  |  |






<a name="cloud.deploymentmanager.autogen.OptionalString"></a>

### OptionalString



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  |  |






<a name="cloud.deploymentmanager.autogen.PasswordSpec"></a>

### PasswordSpec
Specifies a generated password and username combination.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metadata_key | [string](#string) |  | Specifies the name of the metadata entry, whose value contains the generated password, accessible to the VM. Must be unique per password spec in the same package spec. Required. |
| length | [int32](#int32) |  | The length of the generated password. Required. |
| allow_special_chars | [bool](#bool) |  | Whether special characters should be included in the generated password. |
| username | [string](#string) |  | Specifies a static username accompanying this password. |
| username_from_deploy_input_field | [string](#string) |  | Specifies that the username should come from a deploy input field whose name is specified here. |
| display_label | [string](#string) |  | A label describing the purpose of this username/password. Required, unless this is the only password, in which case the label defaults to &#34;Admin&#34;. |
| generate_if | [BooleanExpression](#cloud.deploymentmanager.autogen.BooleanExpression) |  | Specifies a condition to decide if password should be generated or not. Optional. If it is not specified, the password is generated. |






<a name="cloud.deploymentmanager.autogen.PostDeployInfo"></a>

### PostDeployInfo
Customizes post-deploy information displayed to the user.
This information helps get the user started with using the deployed solution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action_items | [PostDeployInfo.ActionItem](#cloud.deploymentmanager.autogen.PostDeployInfo.ActionItem) | repeated |  |
| info_rows | [PostDeployInfo.InfoRow](#cloud.deploymentmanager.autogen.PostDeployInfo.InfoRow) | repeated |  |
| connect_button_label | [string](#string) |  | Optional label to use for the button that connects to the VM. Deprecated in favor of connect_button.display_label. |
| connect_button | [PostDeployInfo.ConnectToInstanceSpec](#cloud.deploymentmanager.autogen.PostDeployInfo.ConnectToInstanceSpec) |  |  |






<a name="cloud.deploymentmanager.autogen.PostDeployInfo.ActionItem"></a>

### PostDeployInfo.ActionItem
Specifies an action item for the user to take.
Text content fields can either contain non-localized en-US text or a
reference (e.g. @ACTION_1_HEADING) into a localized text file. The latter
is not yet implemented, so only en-US text for now.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| heading | [string](#string) |  | Summary heading for the item. UTF-8 text. No markup. At most 64 characters. Required. |
| description | [string](#string) |  | Longer description of the item. UTF-8 text. HTML &lt;code&gt;&amp;lt;a href&amp;gt;&lt;/code&gt; tags only. At most 512 characters. Optional. At least one of description or snippet is required. |
| snippet | [string](#string) |  | Fixed-width formatted code snippet. Accepts string expressions. UTF-8 text. No markup. At most 512 characters. Optional. At least one of description or snippet is required. |
| show_if | [BooleanExpression](#cloud.deploymentmanager.autogen.BooleanExpression) |  | Specify the condition to display this action item. Optional. |






<a name="cloud.deploymentmanager.autogen.PostDeployInfo.ConnectToInstanceSpec"></a>

### PostDeployInfo.ConnectToInstanceSpec
Specifies a connect button configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tier_vm | [TierVmInstance](#cloud.deploymentmanager.autogen.TierVmInstance) |  | Specifies a VM from tier whose address would be used. Required in a multi-vm configuration. Mustn&#39;t be specified for a single-vm. |
| display_label | [string](#string) |  | Optional label to use for the button that connects to the VM. |






<a name="cloud.deploymentmanager.autogen.PostDeployInfo.InfoRow"></a>

### PostDeployInfo.InfoRow
Specifies a row in the application info table.
Text content fields can either contain non-localized en-US text or a
reference (e.g. @ROW_1_HEADING) into a localized text file. The latter
is not yet implemented, so only en-US text for now.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| label | [string](#string) |  | Row label. Required. Accepts string expressions. UTF-8 text. No markup. At most 64 characters. |
| value | [string](#string) |  | Row value. Accepts string expressions. UTF-8 text. HTML &lt;code&gt;&amp;lt;a href&amp;gt;&lt;/code&gt; tags only. At most 128 characters. |
| value_from_deploy_input_field | [string](#string) |  | Row value referenced from deploy input field. Should specify existing input field&#39;s name. |
| show_if | [BooleanExpression](#cloud.deploymentmanager.autogen.BooleanExpression) |  | Specify the condition to display this row. Optional. |






<a name="cloud.deploymentmanager.autogen.SingleVmDeploymentPackageSpec"></a>

### SingleVmDeploymentPackageSpec
Specifies a solution that deploys a single VM.
Next ID: 22


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| images | [ImageSpec](#cloud.deploymentmanager.autogen.ImageSpec) | repeated | Defines the disk images. If there are more than one, the user can select which image to deploy with. The 1st image is the default. Required. |
| machine_type | [MachineTypeSpec](#cloud.deploymentmanager.autogen.MachineTypeSpec) |  | Specifies the default machine type, and any size constraints which restrict what the user can select. Will use defaults if not specified. |
| boot_disk | [DiskSpec](#cloud.deploymentmanager.autogen.DiskSpec) |  | Defines boot disk. Will use defaults if not specified. |
| local_ssds | [LocalSsdSpec](#cloud.deploymentmanager.autogen.LocalSsdSpec) |  | Specifies additionally added local SSD disks (with a default naming convention). |
| additional_disks | [DiskSpec](#cloud.deploymentmanager.autogen.DiskSpec) | repeated | Defines additional persistent disks. Optional. |
| ip_forwarding | [IpForwardingSpec](#cloud.deploymentmanager.autogen.IpForwardingSpec) |  | If not specified, IP forwarding is forced off and not user-configurable. |
| network_interfaces | [NetworkInterfacesSpec](#cloud.deploymentmanager.autogen.NetworkInterfacesSpec) |  | Spec to define Multi/Single NIC(s) usage for this solution. |
| firewall_rules | [FirewallRuleSpec](#cloud.deploymentmanager.autogen.FirewallRuleSpec) | repeated | Specifies the default firewall rules to access the deployed application. They could be off by default, but should still be listed so that the user can get instructions on how to enable them post-deployment. |
| site_url | [InstanceUrlSpec](#cloud.deploymentmanager.autogen.InstanceUrlSpec) |  | Declares a URL to access the deployed application. |
| admin_url | [InstanceUrlSpec](#cloud.deploymentmanager.autogen.InstanceUrlSpec) |  | Declares a URL to administer the deployed application. |
| passwords | [PasswordSpec](#cloud.deploymentmanager.autogen.PasswordSpec) | repeated | Defines how to generate passwords at deployment time. |
| gcp_auth_scopes | [GcpAuthScopeSpec](#cloud.deploymentmanager.autogen.GcpAuthScopeSpec) | repeated | Declares what GCP APIs should be available to the VM. |
| gce_startup_script | [GceStartupScriptSpec](#cloud.deploymentmanager.autogen.GceStartupScriptSpec) |  | Specifies a startup script for a VM instance. |
| application_status | [ApplicationStatusSpec](#cloud.deploymentmanager.autogen.ApplicationStatusSpec) |  | Defines how to determine the application installation status in post-deployment. This tells when the application is ready for consumption. |
| external_ip | [ExternalIpSpec](#cloud.deploymentmanager.autogen.ExternalIpSpec) |  | Defines how the VM is accessible from the Internet. Will use defaults if not specified. DEPRECATED! Use NetworkInterfacesSpec instead. |
| post_deploy | [PostDeployInfo](#cloud.deploymentmanager.autogen.PostDeployInfo) |  | Customizes post-deploy information displayed to the user. This helps get the user started with using the deployed solution. |
| gce_metadata_items | [GceMetadataItem](#cloud.deploymentmanager.autogen.GceMetadataItem) | repeated | Customizes metadata items on a GCE VM instance. |
| accelerators | [AcceleratorSpec](#cloud.deploymentmanager.autogen.AcceleratorSpec) | repeated | Attach accelerator hardware (GPU) to the VM. Currently at most one accelerator spec is supported. |
| deploy_input | [DeployInputSpec](#cloud.deploymentmanager.autogen.DeployInputSpec) |  | Customizes additional inputs configured by user prior to deployment. Currently, the configured values can be passed through to the VM via metadata items. |
| zone | [ZoneSpec](#cloud.deploymentmanager.autogen.ZoneSpec) |  | Customizes the zone selector. |
| stackdriver | [StackdriverSpec](#cloud.deploymentmanager.autogen.StackdriverSpec) |  | Integration with Stackdriver. |






<a name="cloud.deploymentmanager.autogen.StackdriverSpec"></a>

### StackdriverSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| logging | [StackdriverSpec.Logging](#cloud.deploymentmanager.autogen.StackdriverSpec.Logging) |  | Shows a checkbox that enable Stackdriver Logging. |
| monitoring | [StackdriverSpec.Monitoring](#cloud.deploymentmanager.autogen.StackdriverSpec.Monitoring) |  | Shows a checkbox that enable Stackdriver Monitoring. |






<a name="cloud.deploymentmanager.autogen.StackdriverSpec.Logging"></a>

### StackdriverSpec.Logging



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_on | [bool](#bool) |  | Specifies that by default it should be on. |






<a name="cloud.deploymentmanager.autogen.StackdriverSpec.Monitoring"></a>

### StackdriverSpec.Monitoring



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_on | [bool](#bool) |  | Specifies that by default it should be on. |






<a name="cloud.deploymentmanager.autogen.TierVmInstance"></a>

### TierVmInstance
Identifies a specific VM in a tier.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tier | [string](#string) |  | Name of the tier. |
| index | [int32](#int32) |  | 0-based index of the VM in the tier. A negative index can be used, with -1 referring the last, -2 the second last, etc. |






<a name="cloud.deploymentmanager.autogen.VmTierSpec"></a>

### VmTierSpec
A tier consists of one or more VMs of the same type. Each VM is
uniquely identified by its index.
Next ID: 18


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Unique name for this tier. Only lowercases. Required. |
| title | [string](#string) |  | Display title for this tier. Required. |
| instance_count | [VmTierSpec.TierInstanceCount](#cloud.deploymentmanager.autogen.VmTierSpec.TierInstanceCount) |  | Defines the number of VM instances in this tier. |
| images | [ImageSpec](#cloud.deploymentmanager.autogen.ImageSpec) | repeated | Defines the disk images. If there are more than one, the user can select which image to deploy with. The 1st image is the default. Required. |
| machine_type | [MachineTypeSpec](#cloud.deploymentmanager.autogen.MachineTypeSpec) |  | Specifies the default machine type, and any size constraints which restrict what the user can select. Will use defaults if not specified. |
| boot_disk | [DiskSpec](#cloud.deploymentmanager.autogen.DiskSpec) |  | Defines boot disk. Will use defaults if not specified. |
| additional_disks | [DiskSpec](#cloud.deploymentmanager.autogen.DiskSpec) | repeated | Defines additional persistent disks to attach to each VM. Optional |
| local_ssds | [LocalSsdSpec](#cloud.deploymentmanager.autogen.LocalSsdSpec) |  | Specifies additionally added local SSD disks (with a default naming convention). |
| ip_forwarding | [IpForwardingSpec](#cloud.deploymentmanager.autogen.IpForwardingSpec) |  | If not specified, IP forwarding is forced off and not user-configurable. |
| network_interfaces | [NetworkInterfacesSpec](#cloud.deploymentmanager.autogen.NetworkInterfacesSpec) |  | Spec to define Multi/Single NIC(s) usage for this solution. |
| gcp_auth_scopes | [GcpAuthScopeSpec](#cloud.deploymentmanager.autogen.GcpAuthScopeSpec) | repeated | Declares what GCP APIs should be available to the VM. |
| gce_startup_script | [GceStartupScriptSpec](#cloud.deploymentmanager.autogen.GceStartupScriptSpec) |  | Specifies a startup script for each VM instance in this tier. |
| application_status | [ApplicationStatusSpec](#cloud.deploymentmanager.autogen.ApplicationStatusSpec) |  | Defines how to determine that VMs in this tier are ready to serve. The entire deployment is ready to serve when all tiers are. |
| external_ip | [ExternalIpSpec](#cloud.deploymentmanager.autogen.ExternalIpSpec) |  | Defines how the VMs are accessible from the Internet. Will use defaults if not specified. DEPRECATED! Use NetworkInterfacesSpec instead. |
| gce_metadata_items | [GceMetadataItem](#cloud.deploymentmanager.autogen.GceMetadataItem) | repeated | Customizes metadata items on each GCE VM instance. |
| accelerators | [AcceleratorSpec](#cloud.deploymentmanager.autogen.AcceleratorSpec) | repeated | Attach accelerator hardware (GPU) to the VM. Currently at most accelerator is supported. |
| firewall_rules | [FirewallRuleSpec](#cloud.deploymentmanager.autogen.FirewallRuleSpec) | repeated | Specifies the default firewall rules to access the VMs in this tier. They could be off by default, but should still be listed so that the user can get instructions on how to enable them post-deployment. |






<a name="cloud.deploymentmanager.autogen.VmTierSpec.TierInstanceCount"></a>

### VmTierSpec.TierInstanceCount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [int32](#int32) |  | The default number of instances. Must satisfy the constraint. |
| range | [Int32Range](#cloud.deploymentmanager.autogen.Int32Range) |  | Specifies a range of contiguous values. |
| list | [Int32List](#cloud.deploymentmanager.autogen.Int32List) |  | Explicitly lists out the supported values. |
| tooltip | [string](#string) |  | Optional. Specify the tooltip text. If not specified, it will get a default value. |
| description | [string](#string) |  | Optional. Field&#39;s description. |






<a name="cloud.deploymentmanager.autogen.ZoneSpec"></a>

### ZoneSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_zone | [string](#string) |  | Sets the default zone. |
| whitelisted_zones | [string](#string) | repeated | Lists the zones that are allowed to be used in this DM package. If list is empty, all zones are allowed. |
| whitelisted_regions | [string](#string) | repeated | Lists the regions that are allowed to be used in this DM package. Only the zones that belong to the specified regions will be allowed to be used. |





 


<a name="cloud.deploymentmanager.autogen.ApplicationStatusSpec.StatusType"></a>

### ApplicationStatusSpec.StatusType
Defines how to monitor application installation status.

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 |  |
| LEGACY_DETECTOR | 1 | Deprecated. |
| WAITER | 2 | Uses runtime config waiter to block the deployment until the application finishes installing. Accompanied by WaiterSpec. |



<a name="cloud.deploymentmanager.autogen.DeployInputSection.Placement"></a>

### DeployInputSection.Placement
Defines where this section should be placed.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PLACEMENT_UNSPECIFIED | 0 |  |
| MAIN | 1 | The predefined untitled section that appears at the top. Image, zone, machine type fields are in this section. Custom fields will appear below the predefined ones. Only one such section can be defined. |
| CUSTOM_TOP | 2 | A custom section sitting right below the main section. Sections of this type appear in the order they are defined. |
| CUSTOM_BOTTOM | 3 | A custom section sitting after all other sections. Sections of this type appear in the order they are defined. |
| TIER | 4 | In a multi-VM configuration, each VM tier has its own section. This placement enables adding input fields into a tier section. |



<a name="cloud.deploymentmanager.autogen.ExternalIpSpec.Type"></a>

### ExternalIpSpec.Type
How the VM is exposed on the Internet.

| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNSPECIFIED | 0 |  |
| NONE | 1 | The VM is not accessible from the Internet. |
| EPHEMERAL | 2 | The VM is accessible from the Internet with an ephemeral IP. |



<a name="cloud.deploymentmanager.autogen.FirewallRuleSpec.Protocol"></a>

### FirewallRuleSpec.Protocol
The IP Protocol that the firewall rule allows.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PROTOCOL_UNSPECIFIED | 0 |  |
| TCP | 1 |  |
| UDP | 2 |  |
| ICMP | 3 |  |



<a name="cloud.deploymentmanager.autogen.FirewallRuleSpec.TrafficSource"></a>

### FirewallRuleSpec.TrafficSource
Description of network source.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SOURCE_UNSPECIFIED | 0 | Defaults to Public to keep backward compatibility. |
| PUBLIC | 1 | Applies to traffic incoming from the internet, with default source range of 0.0.0.0/0 or configurable by end user. This rule can be enabled or disabled by end user. |
| TIER | 2 | Applies to traffic between instances within a tier. |
| DEPLOYMENT | 3 | Applies to traffic between all instances in the deployment. |



<a name="cloud.deploymentmanager.autogen.GcpAuthScopeSpec.Scope"></a>

### GcpAuthScopeSpec.Scope
API scope.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SCOPE_UNSPECIFIED | 0 |  |
| CLOUD_PLATFORM_READONLY | 1 | https://www.googleapis.com/auth/cloud-platform.read-only |
| CLOUD_PLATFORM | 2 | https://www.googleapis.com/auth/cloud-platform |
| COMPUTE_READONLY | 3 | https://www.googleapis.com/auth/compute.readonly |
| COMPUTE | 4 | https://www.googleapis.com/auth/compute |
| SOURCE_READ_WRITE | 5 | https://www.googleapis.com/auth/source.read_write |
| PROJECTHOSTING | 6 | https://www.googleapis.com/auth/projecthosting |



<a name="cloud.deploymentmanager.autogen.InstanceUrlSpec.Scheme"></a>

### InstanceUrlSpec.Scheme
The URL scheme. Required.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SCHEME_UNSPECIFIED | 0 |  |
| HTTP | 1 |  |
| HTTPS | 2 |  |


 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

