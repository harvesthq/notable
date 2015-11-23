# Harvest Notable

A simple Slack interface to disseminate small but important bits of information to the team.

## Example usage

```
> /notable Chad found Cade's coat!
# Got it!

> /notable The soup, the soup. The soup! #big
# Got it!
```

## Development

You'll need to have a working Go environment (including `$GOPATH` and the proper directory structure).
It also uses Redis to persist data across processes, so you'll need that too.
To build the binaries, execute

```
go install ./...
```

in the project root and that should put `notable_server` in your `$GOPATH/bin` to run. Normal execution
goes something like:

```
env PORT=8080 SLACK_API_TOKEN=... $GOPATH/bin/notable_server
```

There's also the `send_digest` binary that's used by the Heroku scheduler to send out and clear any notes nightly.
