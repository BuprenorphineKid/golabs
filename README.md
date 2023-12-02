### GoLabs ###

#THIS IS A WORK IN PROGRESS#

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

1. A Unix style OS. (While building this program, I have switched through
    using android 11 arm64, android 7 arm7ve, android 7 arm6,
    and many other android variants with termux environment.
    Also have tested and run on ubuntu, but havent gotten the chance
    to test on any other distros..)
   

3. Golang programming language installed on your system.

4. GoImports installed, and the binary put into your $PATH or your
$GOPATH so long as its also in your $PATH.
	(optional but auto-importing will not work)

5.tbox installed.. in $PATH. tbox is a custom program i wrote that I
literally just use to draw boxes in the terminal lol.

I've included a script entitled "install.sh" that should take care of
this part for you. Of course it is to be run with bash or your shell
of choice.


# Installation

git clone this repo. with 'git clone https://github.com/BuprenorphineKid/GoLabs'
move into the downloaded directory then run './install.sh',
after that there should be a binary called labs. run it with ./labs. 

TODO:

	more debugging options

	add menus and other options
	
	more robust help command
	
	|implement autotab support

	flesh out the output window

 	andf many many many more.


  I hope you enjoy the program and that it helps you all in some way, i welcome
  any and all feedback/suggestions/requests/issues as this app is far from being complete <3
  i love you. :)
