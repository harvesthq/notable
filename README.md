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

## Configuration

Configuration is done via environment variables, either directly or through Heroku.

* `FROM_NAME` and `FROM_EMAIL` determine who the email is from
* `TO_NAME` and `TO_EMAIL` determine who the email is sent to
* `REDIS_URL` is normally provided by Heroku and gives connection details for the Redis instance
* `SLACK_CHANNEL` determines what room notes are broadcast to, defaults to "general"

## Contributors

* Danny Wen (danny@getharvest.com) &mdash; progenitor
* Jason Dew (jason@getharvest.com) &mdash; initial code development
* Katie Rose (katie@getharvest.com) &mdash; lots of documentation and testing
* Matthew Lettini (lettini@getharvest.com) &mdash; documentation page design
* Chris Moore (chrism@teamgaslight.com) &mdash; configurable broadcast channel
