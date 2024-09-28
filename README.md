# SyncKor

A sync server for the e-reader app [KOReader](https://koreader.rocks/).  
Made to be a drop-in replacement for the official KOReader sync server.

## Demo


## Quick Start
* The project has an official Docker image published on github container registry.  
* The repository is located [HERE](https://github.com/atran25/synckor/pkgs/container/synckor)  
* The environment variables that are needed can be found below in the [Environment Variables](#environment-variables) section.
* For setting up the s3 bucket or s3 compatible storage, please reference the detailed guide in the [litestream documentation](https://litestream.io/guides/)
## Getting Started

### Dependencies

* [Go](https://go.dev/)
* [litestream](https://litestream.io/)
* [SQLite](https://www.sqlite.org/)
* [MinIO](https://github.com/minio/minio) or any other s3 compatible storage

### How to Setup the Application (Local)
* Local setup doesn't include litestream
1. Setup a local MinIO object storage.
```bash
docker run -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"
```
2. Follow the guide on [litestream](https://litestream.io/getting-started/#setting-up-minio) to create the bucket in the MinIO object storage.
3. Run the make command to start the application.
```bash
make run
```
### How to Setup the Application (Docker)
1. Setup a local MinIO object storage.
```bash
docker run -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"
```
2. Follow the guide on [litestream](https://litestream.io/getting-started/#setting-up-minio) to create the bucket in the MinIO object storage.
3. Use the make command to build the docker image.
```bash
make docker-build
```
4. Edit the docker-compose.yml file to make sure the environment variables are correct.

5. Use the make command to run the docker image.
```bash
make docker-run
```

### Environment Variables
| Key                          	| Description                                  	| Example                                        	| Note                                                                                                                	|   	|
|------------------------------	|----------------------------------------------	|------------------------------------------------	|---------------------------------------------------------------------------------------------------------------------	|---	|
| PORT                         	| The port that the server will listen on      	| 8050                                           	|                                                                                                                     	|   	|
| REGISTRATION_ENABLED         	| Whether or not registration is enabled       	| true                                           	|                                                                                                                     	|   	|
| LITESTREAM_ACCESS_KEY_ID     	| The access key id for the object storage     	| minioadmin                                     	| The key will be different depending on which object storage you use, check the litestream guides to see the changes 	|   	|
| LITESTREAM_SECRET_ACCESS_KEY 	| The secret access key for the object storage 	| minioadmin                                     	| The key will be different depending on which object storage you use, check the litestream guides to see the changes 	|   	|
| REPLICA_URL                  	| The url for the object storage               	| s3://synckor-bkt-test.localhost:9000/db.sqlite 	|                                                                                                                     	|   	|

## License

This project is licensed under the [MIT License](LICENSE) - see the LICENSE file for details