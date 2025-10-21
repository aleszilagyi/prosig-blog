# prosig-blog

## Dependencies

- Docker [getting started](https://www.docker.com/get-started/)

## How to run it

- Some common commands were added to `Makefile`, check it for some common commands to support running and testing this project

### Format code:

- To keep the code in a standardized format, please use:

```shell
make format
```

### Regen mocks:

- This should download the official uber `mockgen` and run mocks regen

```shell
make mocks
```

### Run tests:

- The recommended command to run tests to guarantee that they are not interdependent and race conditions are checked:

```shell
go test ./... -race -shuffle=on
```

- Alternatively you can run:

```shell
make tests
```

### Local build:

- This should create a `_build` dir, build and output the binary in `_build`

```shell
make local-build
```

### Run Application in Docker Container:

- To start containers:

```shell
make run-docker-local
```

- To shut down containers:

```shell
make down-docker-local
```

### Common localhost test requests:

- To create a new post:

```shell
make request-post-post
```

- To get all posts with comment count:

```shell
make request-get-posts
```

- To get the post with id 1 (if exists):

```shell
make request-get-post-1
```

- To fail when trying to get the post with invalid id:

```shell
make request-get-post-fail
```

- To add comment to post with id 1 (if exists):

```shell
make request-post-comment-post-1
```

- To fail adding comment with an empty content:

```shell
make request-post-comment-post-validation-fail
```

- To fail adding comment with an invalid post id:

```shell
make request-post-comment-post-postid-fail
```

- To fail when creating a new post due to empty title:

```shell
make request-post-post-fail
```
