# github-activity-cli

A simple command-line tool written in Go that shows the recent public activity of any GitHub user.

It uses only the Go standard library, so there are no external dependencies.

## Features

- Fetches recent public events from the GitHub Events API
- Accepts the username as a command-line argument or asks for it interactively
- Prints a readable one-line summary for each event
- Handles common errors (user not found, API errors, network failures)

## Requirements

- [Go](https://go.dev/dl/) 1.18 or newer

## Installation

```bash
git clone https://github.com/<your-username>/github-activity.git
cd github-activity
go build -o github-activity
```

## Usage

Pass the username as an argument:

```bash
./github-activity-cli lorenzoalessi
```

Or run it without arguments and type the username when prompted:

```bash
./github-activity-cli
Enter the GitHub username: lorenzoalessi
```

You can also run it directly without building:

```bash
go run . lorenzoalessi
```

### Example output

```
- Pushed to master in lorenzoalessi/linux
- Starred some-user/some-repo
- Opened a pull request in some-org/some-repo
- Commented on an issue in some-org/some-repo
```

## Supported events

| Event type          | Output                                    |
|---------------------|-------------------------------------------|
| `PushEvent`         | Pushed to `<branch>` in `<repo>`          |
| `IssuesEvent`       | `<Action>` an issue in `<repo>`           |
| `IssueCommentEvent` | Commented on an issue in `<repo>`         |
| `PullRequestEvent`  | `<Action>` a pull request in `<repo>`     |
| `WatchEvent`        | Starred `<repo>`                          |
| `ForkEvent`         | Forked `<repo>`                           |
| `CreateEvent`       | Created a `<branch/tag/repo>` in `<repo>` |

Any other event type is printed as `<EventType> in <repo>`.

## How it works

1. Reads the username from the command-line arguments or from standard input.
2. Sends a `GET` request to `https://api.github.com/users/<username>/events`.
3. Decodes the JSON response into a slice of Go structs.
4. Converts each event into a human-readable sentence and prints it.

## Notes and limitations

- The tool uses the **unauthenticated** GitHub API, which is limited to 60 requests per hour per IP address.
- The Events API only returns recent public activity (up to 30 events per page, and only events from the last 90 days).
- GitHub's `PushEvent` payload no longer includes the number of commits, so the tool shows the branch that was pushed to instead.

## Project structure

```
.
├── go.mod
├── main.go      # CLI entry point, API call and event formatting
└── README.md
```

## Ideas for future improvements

- Add an HTTP client timeout
- Add a `--limit` flag to cap the number of events shown
- Group consecutive pushes to the same repository
- Show the commit count using the compare API
- Support authentication with a personal access token
- Split the code into packages and add unit tests

## License

This project is released under the [MIT License](LICENSE).

## Acknowledgments

Thanks to [Claude](https://claude.ai) for this README