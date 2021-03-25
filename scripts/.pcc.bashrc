#!/usr/bin/env bash
_platformcc_completions()
{
    cmd=$(printf " %s" "${COMP_WORDS[@]:1}")
    COMPREPLY+=($(pcc _ac "$cmd"))
}
complete -o default -F _platformcc_completions pcc
export PATH=~/.pcc:$PATH