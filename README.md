# csye7374-operator
A Kubernetes operator to create s3 folders

## Team Information

| Name | NEU ID | Email Address |
| --- | --- | --- |
| Shubhankar Dandekar| 001439467| dandekar.s@husky.neu.edu |
| Mitali Salvi|001630137  | salvi.mi@husky.neu.edu|
| Ashutosh Shukla|001449285 | shukla.as@husky.neu.edu|
| Lakshit Talreja|001475200 |talreja.l@husky.neu.edu |

To create instances of this custom operator:
1)Install this operator using helm chart using the following command:
```sh
helm install . --namespace operator-demo \
--name csye7374-operator-helm-chart \
--set aws.accessKeyId=AWS_ACCESS_KEY \
--set aws.secretAccessKey=AWS_SECRET_KEY \
--set aws.bucketName=BUCKET_NAME_WHERE_FOLDER_IS_CREATED \
--set aws.region=AWS_REGION \
--set image.repository=DOCKER_IMAGE \
--set imageCredentials.username=DOCKER_USERNAME \
--set imageCredentials.password=DOCKER_PASSWORD 
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
