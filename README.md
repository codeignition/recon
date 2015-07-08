# recon
A tool to detect attributes of a linux machine à la ohai of chef.

It is at a very early stage and not yet feature complete. So, obviously the APIs are going to change every day. So don't use it for anything important.

### Disclaimer

I can only test this on an Ubuntu machine. So, I implicitly added linux only build tags for packages. This may not work on your machine. Be warned!

### Installation

If you have Go installed and workspace setup,

~~~sh
go get github.com/hariharan-uno/recon
~~~

Test whether its installed by running 

~~~sh
recon -addr=":8080"
~~~

and then open http://localhost:8080. If addr flag isn't specified, it serves on :3030.

### License

BSD 3-clause "New" or "Revised" license
