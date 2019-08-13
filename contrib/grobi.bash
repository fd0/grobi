# Bash completion script for grobi

_grobi_complete_rules()
{
    if [ "${#COMP_WORDS[@]}" -eq 3 ]; then
        COMPREPLY=($(compgen -W "$(grobi rules)" -- "${COMP_WORDS[2]}"))
    fi
}

_grobi_completions()
{
    if [ "${#COMP_WORDS[@]}" -eq 2 ]; then
        COMPREPLY=($(compgen -W "apply rules show update version watch" -- "${COMP_WORDS[1]}"))
    else
        command=${COMP_WORDS[1]}

        case $command in
            apply) _grobi_complete_rules ;;
            *) ;;
        esac
    fi
}

complete -F _grobi_completions grobi
