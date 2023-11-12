# RSS-Reader

A simple RSS Reader

## Build

```bash
go build ./...
```

## Usage

Running the application will start a server on port 3000 where it will listen for "GET" requests with a json payload of:

```json
{
"links":["rss feed link 1","rss feed link 2"]
}
```

Where "rss feed link" is a valid link to a RSS feed.
If there is at least one valid link to a feed, the server will parse it and return a json object of all of the objects that were downloaded from it which contains an aray of RSS Items:
```json
{
"items":[
{
"title":"title",
"source":"source_link",
"source_url":"source_url_link",
"link":"link",
"publish_date":"publish_date",
"description":"item_description"
},
{
"title":"title2",
"source":"source_link2",
"source_url":"source_url_link2",
"link":"link2",
"publish_date":"publish_date2",
"description":"item_description2"
}
]
}
```

## Running tests

To check whether the parser works as intended, run the following command in the root directory of the project:

```bash
go test ./...
```