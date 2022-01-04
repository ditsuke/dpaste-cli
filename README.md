# dpaste-cli

**dpaste-cli** is a command line interface client for [dpaste.com](https://dpaste.com).

I made this client as an exercise in learning Go and developing my first "useful" command line app.

## Why this client?
It would be remiss of me to ignore the existence of other CLI clients for dpaste, so why this client? 
- [x] Single binary with no dependencies, no interpreter required. 
- [x] Support a UNIX workflow, pipes by default (see [dpaster][dpaster])
- [ ] Support persistent configuration.

## Installation
@todo

## Usage
- Pipes: `cat <your_file> | dpaste-cli`
- Pass a file arg: `dpaste-cli create --file <your_file>`

## Configuration
@todo

## Contributing
You will need a local or remote Go development environment. While I use the GoLand IDE by JetBrains, VSCode also does 
a good job and Vim has an excellent [plugin][vim-go] for Go.

For remote development, my strong recommendation is [Gitpod][gitpod]. [Replit][replit] is also an excellent choice.

This repository is open to PRs :)

[dpaster]: https://github.com/xvm32/dpaster
[vim-go]: https://github.com/fatih/vim-go
[gitpod]: https://www.gitpod.io/
[replit]: https://replit.com/
