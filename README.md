### GoLabs ###


# Preface

As i have spent much time on this project, ive come to find that
as a codebase grows it does infact take more and more increasing
effort to maintain updates and overall inprovement. Also as
this is my first mediumish scale project, at least the first thats
gotten off the ground, im still learning and will, inevitably, make
bad design/business decisions. all that i ask is that you just be
patient with me.


# The Premise

Up until now, theres not really been any good command line solutions 
for messing around with the Go language. Let's face it, other than
just flat out writing an entire module/project, or at the very least
a little script which when youre just testing out a few functions 
or a couple new design ideas, theres a good bit of boilerplate involved.
Golabs sets out to accomplish a few things promised by other Golang
REPL/shells, but this time, actually delivering. Dont get me wrong the
others "worked", but the second you press a "HOME" key or "DEL" key
youre in for a rude awakening. but i dont blame them for opting out
of diving all the way into raw terminal mode, and fully implementing
it all yourself. after all, this did take a good litle while. but alas
i have no life, so low level terminal implementation it is lol.

# Prerequisites

as of now the only requirements that i know of for sure are as follows:

1. A Unix style OS. (Im using android 11, arm with termux environment)

2. Golang programming language installed on your system.

3. GoImports installed, and the binary put into your $PATH or your
$GOPATH so long as its also in your $PATH.

# Installation

git clone the repo. move into the downloaded directory then run
'go build .', after that there should be a binary called labs. run
it with ./labs. or optionally just do 'go run main.go' as the project
is far from conplete. this section will be updated when the project gets
closer to being somewhere i could call packagable.

TODO:
	implement a toggled concurrently running debug mode 
	instead of having to press ctrl-D everytime you need stats

	generally add more features/commands

	more robust help command
	
	fix structs/functions

	redesign evaluation command.

