# Watchtower

**NOTE**: Alpha version, release pipelines are WIP, if you're keen on running it locally you'll need to
install [wails](https://wails.io/).

Focused product view of your GitHub organisation. View repositories, pull requests, security vulnerabilities grouped by
product.

## Motivation

Working in organisations looking after multiple teams, it was hard to view the health of the product with multiple
microservices.
Watch tower provides a lite weight approach to grouping information by product with additional filters.

## Features

| Feature                         | Description                                                                                                                                            |
|---------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| Group products by organisations | Add multiple organisations and group products within those organisations                                                                               |
| Group repositories by Product   | User github topics to group repos by products                                                                                                          |
| Dashboard view                  | Quickly view the overall health of a group of products by viewing open PRs, security vulnerabilities and repositories                                  |
| Insights                        | Get insights at a org level on the for pull requests and security vulnerabilities                                                                      |
| Product view                    | Focused view for open PRs, security vulnerabilities and repositories                                                                                   |
| Notification on new PR / Issue  | Get notifications when new PRs or Issues are created in the watched repositories accross all orgs                                                      |                                                     |
| Local only                      | Data will never leave your device, settings page includes a kill switch to remove all data from the device. GITHUB PAT token only required read access |

## Roadmap

Issues, suggestions can be raised in the [issues tab](https://github.com/code-gorilla-au/watchtower/issues).

## Getting started

## List tasks

List all available tasks

```bash
task
```

### Install

Install all dependencies and tools

```bash
task go-install
task frontend-install
```

### Run

- Generates sqlc code
- Runs the generate command

```bash
task dev
```

### Test

#### Go

```bash
# Watch mode
task go-watch

# View coverage
task go-cover-html
```

### Frontend

```bash
# With coverage
task frontend-test

# Watch mode
task frontend-watch
```

### Lint

```bash
task go-lint
task frontend-lint
```

### CI tasks

```bash
task go-ci
task frontend-ci
```

### Reset state

- Removes local sql file
- Run dev task

```bash
task reset
```

