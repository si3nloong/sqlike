## Contributing to sqlike

## Installation

1. Ensure you have minimum version 1.15 of [go](https://golang.org/dl/) installed.
2. Ensure you have minimum version 5.7 of [mysql](https://dev.mysql.com/downloads/installer/) installed.
3. After fork the repository, you may start your development.

## Pull requests

#### Your first pull request

If you have decided to contribute code back to upstream by opening a pull request. You've invested a good chunk of time, and we appreciate it. We will do our best to work with you and get the PR looked at.

Working on your first Pull Request? You can learn how from this free video series:

[How to Contribute to an Open Source Project on GitHub](https://egghead.io/courses/how-to-contribute-to-an-open-source-project-on-github)

#### Proposing a change

If you would like to request a new feature or enhancement but are not yet thinking about opening a pull request, you can also file an issue with feature template.

If you're only fixing a bug, it's fine to submit a pull request right away but we still recommend that you file an issue detailing what you're fixing. This is helpful in case we don't accept that specific fix but want to keep track of the issue.

#### Sending a pull request

Small pull requests are much easier to review and more likely to get merged. Make sure the PR does only one thing, otherwise please split it by module.

Please make sure the following is done when submitting a pull request:

1. Fork the [repository](https://github.com/si3nloong/sqlike/v2) and create your branch from master.
2. Describe your test plan in your pull request description. Make sure to test your changes.
3. Make sure your tests pass `go test ./...`.
4. All pull requests should be opened against the master branch.

## Testing

For testing, we are using package `github.com/stretchr/testify`, please ensure you use the same package for your unit testing as well.

## License

By contributing to sqlike, you agree that your contributions will be licensed under its [MIT license](https://github.com/si3nloong/sqlike/v2/blob/master/LICENSE).
