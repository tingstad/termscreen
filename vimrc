augroup TestAutoGroup
    autocmd TextChanged,InsertLeave *.go call TestCommitResetGo()
augroup END
function! TestCommitResetGo()
    write
    if system('git ls-files --modified --deleted --error-unmatch >/dev/null 2>&1 && printf changed') == 'changed'
        let l:buffer = bufnr('%')
        if exists('g:mybuf') && bufwinnr(g:mybuf) > -1
            execute bufwinnr(g:mybuf) . 'wincmd w'
        else
            botright 9 new
            let g:mybuf = bufnr('%')
        endif
        execute 'read! { date && go fmt && go build && go test ;} && git commit --all --message=TCR || git reset --hard HEAD'
        execute bufwinnr(l:buffer) . 'wincmd w'
        unlet l:buffer
        edit
        redraw!
    endif
endfunction
