YAURLS (Yet Another URL Shortener)
====

yaurls is a pure Go implemtation of a URL Shortener service.

For some products on our store, [thingsforhackers](https://thingsforhackers.com), we wanted to have the ability to generate a short URL that would link back our main site for things like blog entries, product listings etc. There are many url Shortener services out there but most seem to not let you chose the short name. We already have a shortish domain name (njohn.uk) under our control, so all we needed was a simple web service that would re-direct the short urls to a full one.

In general we wanted to be able to support the following type of mapping:

`njohn.uk/go/<shortName> --> http://somelongurl.com/some/path.html`

E.g. njohn.uk/go/tl866 ---> https://thingsforhackers.com/tl866-universal-programmer.html



In its simplest form, this can be achieved by the following compoents:

1. A http server.

   This will be running on njohn.uk and will handle the required GET/PUT rest calls.

2. Persistent storage for the short name to long url mappings.

   Be able to remember that the shortname __tl866__ maps to the full url as shown abvove.

   A simple key value store will be used, with the __key__ being the short name, and the __value__ being the full URL.


The code in this module acheives exactly this.

## Getting Started

### Installing
To start using yaurls, assuming you've got a working Go installation, just run:
```sh
$ go get github.com/thingsforhackers/yaurls
```
This will get the program, along with its requirements, and install in your `$GOPATH\bin` directory.

### Running
The help output of yaurls is as follows:
```sh
njohn@tfh ~ $ yaurls --help
Usage of yaurls:
  -dbPath string
    	Path to Dbase file (default "/home/njohn/url.db")
  -portNum int
    	Port to listen on (default 8080)
  -updateToken string
    	Optional token required for DB modification operations
njohn@tfh ~ $
```

All of these are optional. The usage of __dbPath__ & __portNum__ options should be obvious enough, and we will cover __updateToken__ later on.

### Usage
*The following examples assume you are talking to a yaurls server running on the default port on your local machine.*

We can use curl to configure and test our yaurls server. Feel free to use you web browser to test the __GET__ operations.

First let's try a shortName lookup.

```sh
njohn@tfh ~ $ curl -X GET --include 127.0.0.1:8080/go/tech
HTTP/1.1 404 Not Found
Content-Type: text/html
Date: Fri, 20 May 2016 17:37:30 GMT
Content-Length: 32

<h1>Can not map tech to a URL</h1>
```

As expected this fails due to the dbase initially being empty.

We can use the following __PUT__ request to add a shortName to URL mapping.

```sh
njohn@tfh ~ $ curl -X PUT --include --header "X-Full-URL:http://www.bbc.co.uk/news/technology" 127.0.0.1:8080/go/tech
HTTP/1.1 201 Created
Content-Type: text/html
Date: Fri, 20 May 2016 17:51:31 GMT
Content-Length: 0
```
Here you can see that we've used the same url for the __shortname__ and have specified the full URL in the request header with the __X-Full-URL__ key.

Now let's issue the __GET__ request again.

```sh
njohn@tfh ~ $ curl -X GET --include 127.0.0.1:8080/go/tech
HTTP/1.1 302 Found
Location: http://www.bbc.co.uk/news/technology
Date: Fri, 20 May 2016 17:52:30 GMT
Content-Length: 59
Content-Type: text/html; charset=utf-8

<a href="http://www.bbc.co.uk/news/technology">Found</a>.
```

As a further test, point your web browser at http://127.0.0.1:8080/go/tech, you should get redirected to the BBC's technology homepage.

<TBC>
