# todo-api-go
Yet Another ToDo App - I'm using this one to learn how to build web services in Go (and to play with Kubernetes on my Pi Cluster).

The ToDo API application exposes a RESTful API, with basic create, retrieve, update and delete functionality. Currently a ToDo resource contains a description and a "completed" boolean flag. There is no authentication. The application currently supports SQLite, MySQL and Mongo databases.

## Build

```shell script
$ go build -o todo-api
```

## Test

```shell script
$ go test
```

## Run

### Environment Variables

You will need to export the following environment variables:

* `CONNECTION_STRING` : 
  * for `sqlite3` this will be the path to the database file. It will be created if it does not exist.
  * for `MySQL` this will be a connection string with the form: `<USERNAME>:<PASSWORD>@(<HOST>)/todo?charset=utf8&parseTime=True&loc=Local`
* `HOST_ADDRESS` : This should the IP address and port on in the form of `<HOST_IP>:<PORT>`

Example:

```shell script
export CONNECTION_STRING=test.db
export HOST_ADDRESS=127.0.0.1:8000
```

### Usage

```bash
$ ./todo-api --help                                                                                                                                                                                                                                               *[helm] 
Usage of ./todo-api:
  -db string
        Database to use. Options are: "sqlite3", "mysql" and "mongo"
```

#### Example

Assuming you exported the environment variables as in the example above, you can run the application with a SQLite database using the following command:

```shell script

$ ./todo-api --db sqlite3
2020/05/09 10:15:31 Starting web server on 127.0.0.1:8000

```

## CRUD Examples

In the examples below I use [curl](https://curl.haxx.se/) to issue the HTTP requests and [jq](https://stedolan.github.io/jq/) to to present a nicely formatted response.

### Create

```bash
curl -s -H "Content-Type: application\json" \
--request POST \
--data '{"description": "Buy milk", "completed": false}' \
http://127.0.0.1:8000/todo | jq
```

Which produces the output:

```json
{
  "Id": "1",
  "Description": "Buy milk",
  "Completed": false
}

```

### Retrieve

```bash
curl -s http://127.0.0.1:8000/todo/1 | jq
```

Output:

```json
{
  "Id": "1",
  "Description": "Buy milk",
  "Completed": false
}
```

### Update

```bash
curl -s -H "Content-Type: application\json" \
--request PUT \
--data '{"description": "Buy milk", "completed": true}' \
http://127.0.0.1:8000/todo/1 | jq
```

Output:

```json
{
  "Id": "1",
  "Description": "Buy milk",
  "Completed": true
}
```

### Delete

```bash
curl -s --request DELETE http://127.0.0.1:8000/todo/1 | jq
```

Output:

```json
{
  "result": "success"
}
```

## Multi-platform Docker images

If you wish to run the ToDo API in Docker or as part of a Kubernetes deployment you will need to build a Docker image for the application. This is fairly simple for x86-based nodes, you would just build the image using the `Dockerfile` with an appropriate tag and push it to your repository. To use the newly built image you should change the repository details in the Helm chart `todo/values.yaml` file.

If you wish to run the application on a Raspberry Pi (as I am interested in doing), you will need to build an ARM compatible image. I'm doing this with the aid of the Docker [buildx](https://github.com/docker/buildx) plugin, a multi-stage Dockerfile (`Dockerfile.multi_arch`) and a helper script (`gobuild_multi_arch.sh`).

Instructions to get set up for multi-platform Docker builds can be found here: https://www.docker.com/blog/getting-started-with-docker-for-arm-on-linux/

The command you can then use to build my multi-platform image would look like:

```bash
docker buildx build --platform linux/amd64,linux/arm/v7 -t <USERNAME>/<REPOSITORY>:<TAG> -f Dockerfile.multi-arch --push .
```

## Kubernetes deployments with Helm charts

Assuming you have [Helm](https://helm.sh/docs/intro/install/) installed, the ToDo API application can be installed as a Kubernetes application with either a SQLite or MySQL database.

The default `values.yml` expects the service to exposed via a LoadBalancer, so you may need to change that to suit your specific needs.

### ToDo API Helm install with SQLite

```bash
helm install todo --generate-name
NAME: todo-1589016378
LAST DEPLOYED: Sat May  9 10:26:18 2020
NAMESPACE: default 
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace todo svc -w todo-1589016378'
  export SERVICE_IP=$(kubectl get svc --namespace default todo-1589016378 --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
  echo http://$SERVICE_IP:80

```

### ToDo API Helm install with MySQL

First create the MySQL root password as a secret:

```bash
kubectl create secret generic mysql-password-secret --from-literal=password='YOUR_PASSWORD_HERE'
```

Then install the MySQL chart:

```bash
$ helm install mysql --generate-name
NAME: mysql-1589016929
LAST DEPLOYED: Sat May  9 10:35:29 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=mysql,app.kubernetes.io/instance=mysql-1589016929" -o jsonpath="{.items[0].metadata.name}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace todo port-forward $POD_NAME 8080:80

```

And finally the ToDo API chart:

```bash
$ helm install todo --generate-name -f todo/values-mysql.yaml
NAME: todo-1589016959
LAST DEPLOYED: Sat May  9 10:36:00 2020
NAMESPACE: default 
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace default svc -w todo-1589016959'
  export SERVICE_IP=$(kubectl get svc --namespace default todo-1589016959 --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
  echo http://$SERVICE_IP:80

```