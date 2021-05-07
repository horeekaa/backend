# Horeekaa Backend
==========================================

A [Go](https://golang.org/) backend using these services:
- [Graphql](https://graphql.org/) with the flavour of [GqlGen](gqlgen.com)
- [Google Firebase Authentication](https://firebase.google.com/products/auth)
- [Google Cloud Platform's App engine](https://cloud.google.com/appengine/)
- [Google Cloud Platform's Cloud Storage](https://cloud.google.com/storage/)
- [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)

## Language

This project uses Golang for better multi-threading backend performance.

## Project Structure

```

```

## Initialize App

1. The Project requires 2 config files needed in order to run. 
- Please setup the all the project variables as listed on `/commons/configs/strings.go`.
Put the variable values on a `.env` file inside the root folder so the code can consume it when it runs.
- Please put the `Firebase Service Account JSON file` within folder `/commons/assets`

DO NOT include `.env` and `Firebase Service Account JSON file` in the git repository under any circumstances.

2. Run `gqlgen generate` to get the all the auto-generated models from gqlgen schemas.