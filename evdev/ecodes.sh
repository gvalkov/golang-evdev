#!/bin/bash

# Generate a go source file that exports macros from
# /usr/include/linux/input.h


header=${1:-/usr/include/linux/input.h}
tmpl=${2:-./ecodes.go.tmpl}
[[ ! -e $header ]] && echo "no such file: $header" && exit 1

function codes () {
    awk '
    /#define +(KEY|ABS|REL|SW|MSC|LED|BTN|REP|SND|ID|EV|BUS|SYN)_/ {
        print $2, "=", $3;
    }' ${header}
}

python -c "
import os, sys

prefixes = 'KEY ABS REL SW MSC LED BTN REP SND ID EV BUS SYN'.split()
codes = '''$(codes)'''.split('\n')

ctx = {
   'uname'   : '$(uname -s -r -v -m)',
   'codes'   : '\n'.join('    %s' % i for i in codes),
   'codemap' : '\n'.join('    \"%s\" : %s,' % (i,i) for i in (j.split()[0] for j in codes)),
}

tmpl = open('${tmpl}').read()
tmpl = tmpl % ctx

sys.stdout.write(tmpl)
" | tail -n +3
