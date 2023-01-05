# Read OpenAPI and generate CFN Schema

## to generate cfn template based on open api(swagger.json) for the resources  mentioned in mapping.json

Add an entry into mapping.json and run the below command

    make schema

## to find the differences with the latest spec(swagger.json) for the resources mentioned in mapping.json

We can update swagger.json manually and run the below command to find out the differences. You should see the difference in schema generated in json file(diff.json)

Use command

    make compare
