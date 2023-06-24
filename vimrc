" I experimented with a strict type of TCR (Test && Commit || Reset) in Vim
"
" vim -c 'source vimrc' -O termscreen.go termscreen_test.go
"
" In this setup, TCR is called on every save, and on every "Leave Insert mode".
"
" My experience of this has been mixed. I think the Vim script and setup works very well,
" automatically formatting, giving test feedback and committing or resetting the code appropriately.
" It's nice that the Go compiler is so fast.
"
" On the other hand, the Go compiler is also quite strict,
" giving some challenges with "x declared but not used".
" I had to add [extract variable](https://github.com/fvictorio/vim-extract-variable) functionality at a minimum.
"
" A general TCR challenge is that "test first" is difficult.
" I usually implemented a feature and then the test, often using Undo to get the test right.
" Changing the behaviour of already implemented and test-covered code is really not easy (unless the tests are deleted).
" It's like every function is an API you cannot break, it's challenging, but interesting.
"
" To summarize, I like the general TCR setup, but I think InsertLeave activation might be too masochistic.

autocmd BufWritePost,InsertLeave *.go call TestCommitResetGo()

function! TestCommitResetGo()
    write
    let l:buffer = bufnr('%')
    if exists('g:mybuf') && bufwinnr(g:mybuf) > -1
        execute bufwinnr(g:mybuf) . 'wincmd w'
    else
        botright 9 new
        let g:mybuf = bufnr('%')
    endif
    noautocmd execute 'read! { date && go fmt && go build && go test ;} && git commit --all --message=TCR || git reset --hard HEAD'
    execute bufwinnr(l:buffer) . 'wincmd w'
    unlet l:buffer
    edit
    redraw!
endfunction
