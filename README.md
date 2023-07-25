# article-processor
article-processor is a simple service that retrieves articles from external endpoints and stores them in the database.
The service has a scheduled interval on which it sends out get requests to the external endpoints to retrieve the article list.
The article list comes in XML format then it is mapped to the service local structures and stored in MongoDB if the article is not stored there yet.

The service provides two endpoints to get the List of the articles stored in the database and to get single articles by their ID
Chi router was chosen as a lightweight solution with easy to use features to handle the HTTP services.

## Requirements
* GO 1.20+ (might build on lower versions as well)
* MongoDB

## Configuration
The service configurations are stored in the conf.yaml file. There is also the conf_test.yaml file for the test environment configurations.

Main configurations:
* Service port
* Core MongoDB configurations
* LogPath to store logs
* External endpoints for articles
* Interval for how often the service should check for new articles

## Running the service
1. Clone the repository
2. Navigate to the cloned directory
3. Make sure you have MongoDB running locally on port `:27017`
4. Simply run the service with `go run main.go` while in the main project directory. You can also build the service.

## Endpoints

### GET HEALTH
* Simple GET request to check if the service is running
  `http://localhost:3000/api/health`

### GET ARTICLE LIST
* GET request that retrieves all articles from the database.
  `http://localhost:3000/api/article/list`

### GET ARTICLE BY ID
* GET request that retrieves a specific article by the ID.
  `http://localhost:3000/api/article/{id}`

## Testing
Testing is only done in the articles directory, inside the `article_test.go` file. Currently there is only one test, but it replicates the core logic of this service and covers a few test cases.
More test cases with different outcomes should be created additionally the scheduler should be tested.
To test it you can simply run `go test -v ./...` from the root of directory of the project.

##TODO:
* Containerize service
* Add id logic for different clubs
* Increase test coverage
* Add authentication
* Potentially add checks to see if the articles have been updated in the external source