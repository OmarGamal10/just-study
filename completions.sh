#!/bin/bash

_focus_completions()
{
  local cur opts

  cur="${COMP_WORDS[COMP_CWORD]}"
  opts="on off status"
  
  COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
}

complete -F _focus_completions study