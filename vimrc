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
