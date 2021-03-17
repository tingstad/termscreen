augroup TestAutoGroup
    autocmd TextChanged,InsertLeave *.go call TestCommitResetGo()
augroup END
function! TestCommitResetGo()
    write
    if system('git ls-files --modified --deleted --error-unmatch >/dev/null 2>&1 && printf changed') == 'changed'
        if exists('g:mybuf') && bufwinnr(g:mybuf) > -1
            let l:buffer = bufnr('%')
            silent! execute bufwinnr(g:mybuf) . 'windo quit!'
            execute bufwinnr(l:buffer) . 'wincmd w'
            unlet l:buffer
        endif
        botright 9 new
        let g:mybuf = bufnr('%')
        execute 'read! { date && go fmt && go build && go test ;} && git commit --all --message=TCR || git reset --hard HEAD'
        " goto previous window:
        wincmd p
        edit
        redraw!
    endif
endfunction
