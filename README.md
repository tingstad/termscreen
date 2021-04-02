# Capture terminal screen

Because ANSI escape sequences can move the cursor, there are many ways to draw a terminal screen. Complex screens are often drawn different parts at a time, or overwriting old values.

This program captures the output and prints it top to bottom.

## Development

I experimented with a strict type of TCR (Test && Commit || Reset) in Vim:

```
vim -c 'source vimrc' -O main.go main_test.go
```

In this setup, TCR is called on every save, and on every "Leave Insert mode".

My experience of this has been mixed. I think the Vim script and setup works very well, automatically formatting, giving test feedback and committing or resetting the code appropriately. It's nice that the Go compiler is so fast.

On the other hand, the Go compiler is also quite strict, giving some challenges with "x declared but not used". I had to add [extract variable](https://github.com/fvictorio/vim-extract-variable) functionality as a minimum.

Another challenge I experienced was that 
