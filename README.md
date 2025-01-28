# prismahelper

A program designed for Prism Launcher/MultiMC that helps to save important files. Initially
[a cmd script](https://gist.github.com/lottuce-yami/fb9c81112b2e1d72edd9436628349de3) that moved screenshots to a common
directory - now rewritten in Go to offer portability and advanced features.

## Installation

1. Download a release from the Releases tab.
2. Put the downloaded executable in the launcher root. To find it, open your launcher and click on Folders -> Launcher
Root.
3. Open your launcher's Settings -> Custom Commands and put the downloaded executable filename in the Post-exit command
field.

Now every time you exit your instance, prismahelper will move screenshots from the instance to a common directory in the
launcher root.

## Building

Make sure you have Go 1.23.4 installed and run `go build`.
