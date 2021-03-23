#!/usr/bin/env bash
_platformcc_completions()
{
    cmd=$(printf " %s" "${COMP_WORDS[@]:1}")
    COMPREPLY+=($(pcc _ac "$cmd"))
}
complete -F _platformcc_completions pcc
