# Capture terminal screen

Because ANSI escape sequences can move the cursor, there are many ways to draw a terminal screen. Complex screens are often drawn different parts at a time, or overwrite old values.

This program captures the output and prints it top to bottom.

(Still in early development; the code is neither complete, battle tested, nor pretty.)

## Development

I experimented with a strict type of TCR (Test && Commit || Reset) in Vim:

```
vim -c 'source vimrc' -O main.go main_test.go
```

In this setup, TCR is called on every save, and on every "Leave Insert mode".

My experience of this has been mixed. I think the Vim script and setup works very well, automatically formatting, giving test feedback and committing or resetting the code appropriately. It's nice that the Go compiler is so fast.

On the other hand, the Go compiler is also quite strict, giving some challenges with "x declared but not used". I had to add [extract variable](https://github.com/fvictorio/vim-extract-variable) functionality at a minimum.

A general TCR challenge is thst "test first" is difficult. I usually implemented a change and then the test, often using Undo to get the test right. Changing the behaviour of already implemented and test-covered code is really not easy (unless the tests are deleted). It's like every function is an API you cannot break, it's challenging, but interesting.

To summarize, I like the general TCR setup, but I think InsertLeave activation might be too masochistic.
