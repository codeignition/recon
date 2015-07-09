# recon
A tool to detect attributes of a linux machine à la ohai of chef.

It is at a very early stage and not yet feature complete. So, obviously the APIs are going to change every day. So don't use it for anything important.

### Disclaimer

I can only test this on an Ubuntu machine. So, I implicitly added linux only build tags for packages. This may not work on your machine. Be warned!

### Installation

If you have Go installed and workspace setup,

```sh
go get github.com/hariharan-uno/recon
```

Test whether its installed by running

```
recon
```

and then open [http://localhost:3030?indent=1](http://localhost:3030?indent=1).

You can also specify an `addr` flag with the port to serve on. e.g. `-addr=":8080"`.
While debugging, the JSON output can be pretty dense and hard to read. So, the `?indent=1` query string
directs the server to output pretty JSON. If it is consumed by another application, you need not
append the `?indent=1` query string.

### License

BSD 3-clause "New" or "Revised" license
