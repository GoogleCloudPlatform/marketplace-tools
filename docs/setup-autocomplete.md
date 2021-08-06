## Overview

Editing `mpdev` yaml files can be error-prone. IDEs such as VSCode and Intellij
offer YAML auto-complete using JSON Schema files.

### VSCode

1. Install the YAML extension for VSCode.
   https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml
1. Follow the instructions to associate YAML files with a JSON schema
   [here](https://github.com/redhat-developer/vscode-yaml#associating-a-schema-to-a-glob-pattern-via-yamlschemas).
   Your `settings.json` file should be similar to the following:
   ```json
    yaml.schemas: {
        "https://raw.githubusercontent.com/GoogleCloudPlatform/marketplace-tools/master/jsonschema/mpdev.jsonschema": [
            "*/*.yaml"
          ]
    }
    ```

### Intellij

Follow the instructions
[here](https://www.jetbrains.com/help/idea/json.html#ws_json_schema_add_custom)
to add a custom JSON schema.

The URL of the `mpdev` JSON schema is
https://raw.githubusercontent.com/GoogleCloudPlatform/marketplace-tools/master/jsonschema/mpdev.jsonschema
