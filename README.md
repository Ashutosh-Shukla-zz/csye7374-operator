# Simple s3 operator
A Kubernetes operator to create s3 folders.

This custom operator will create a IAM user with the specified username. This username will have access to the folder created in the S3 bucket. The mandatory parameters for the custom operator are:

1)S3 bucket name

2)Operator name and namespace where the operator pod will run

3)Secret which contains IAM user name who will access to S3 folder

To create instances of this custom operator:

1)Create instance of the custom operator. Sample CR
```
apiVersion: csye7374/v1alpha1
kind: folder
metadata:
  name: ashutosh
  namespace: some_namespace
spec:
  username: ashutosh
  userSecret:
    name: ashutosh-secret

```
