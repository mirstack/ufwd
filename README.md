# ufwd - UDP Forwarding Daemon

This is super simple (~ 100 loc) and reliable UDP forwarder/proxy written in Golang. 

## Installation

The package has no external dependencies and is easily go-installable:

    $ go install gitbub.com/mirstack/ufwd
    
You can also download it and build manually.

    $ git clone https://github.com/mirstack/ufwd.git
    $ cd ufwd
    $ go build . && go install .

## Usage

Here's an example of a forwarding server:

    $ ufwd -addr :514 -dest syslog.domain.com:514
    
In this example all the local UDP connections on port `514` will be passed through to
a syslog server running on port `514` of `syslog.domain.com` host.

Check `ufwd -help` for list of all available options.

## Hacking

Not much to say. If you wanna hack on `ufwd` just clone the repo and play with the
code. You can run the tests at any time with standard `go test` tool:

    $ go test .

## Contribute

1. Fork the project.
2. Write awesome code (in a feature branch).
3. Test, test, test...!
4. Commit, push and send Pull Request.
