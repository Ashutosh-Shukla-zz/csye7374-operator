# csye7374-operator
A Kubernetes operator to create s3 folders.

This custom operator will create a IAM user with the specified username. This username will have access to the folder created in the S3 bucket. The mandatory parameters for the custom operator are:

1)S3 bucket name

2)Operator name and namespace where the operator pod will run

3)Secret which contains IAM user name who will access to S3 folder

## Team Information

| Name | NEU ID | Email Address |
| --- | --- | --- |
| Shubhankar Dandekar| 001439467| dandekar.s@husky.neu.edu |
| Mitali Salvi|001630137  | salvi.mi@husky.neu.edu|
| Ashutosh Shukla|001449285 | shukla.as@husky.neu.edu|
| Lakshit Talreja|001475200 |talreja.l@husky.neu.edu |

To create instances of this custom operator:

1)Install this operator using helm chart from the following link:
```
https://github.com/mitali-salvi/csye7374-operator-helm-chart 
```
2)Create instance of the custom operator. Sample CR
```
apiVersion: csye7374/v1alpha1
kind: folder
metadata:
  name: mitalisalvi
  namespace: some_namespace
spec:
  username: mitalisalvi
  userSecret:
    name: mitali-secret

```
