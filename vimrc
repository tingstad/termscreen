augroup TestAutoGroup
    autocmd TextChanged,InsertLeave *.go call TestCommitResetGo()
augroup END
function! TestCommitResetGo()
    if exists('g:mybuf') && bufwinnr(g:mybuf) > -1
        let l:buffer = bufnr('%')
        silent! execute bufwinnr(g:mybuf) . 'windo quit!'
        execute bufwinnr(l:buffer) . 'wincmd w'
        unlet l:buffer
    endif
    write
    botright 9 new
    let g:mybuf = bufnr('%')
    execute 'read! { date && go fmt && date && go test ;} && git commit --all --message=TCR || git reset --hard HEAD'
    " goto previous window:
    wincmd p
    edit
    redraw!
endfunction
