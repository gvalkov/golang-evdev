#!/usr/bin/env python
# -*- coding: utf-8; -*-

'''
Generate a go source file that exports macros from evdev header files.
'''

from __future__ import print_function

import os, sys, re
import string


# The default header file locations to try.
SOURCES = [
    '/usr/include/linux/input.h',
    '/usr/include/linux/input-event-codes.h',
]

if sys.argv[1:]:
    SOURCES = sys.argv[1:]

# The golang template file in which we inject the macros.
template_path = os.path.join(os.path.dirname(__file__), '../ecodes.go.template')
TEMPLATE = string.Template(open(template_path).read())


#-----------------------------------------------------------------------------
MACRO_REGEX = r'#define +((?:KEY|ABS|REL|SW|MSC|LED|BTN|REP|SND|ID|EV|BUS|SYN|FF)_\w+)\s+(\w+)'
MACRO_REGEX = re.compile(MACRO_REGEX)

def get_uname():
    uname = list(os.uname())
    del uname[1]
    return ' '.join(uname)

def parse_header(fh):
    for line in fh:
        macro = MACRO_REGEX.search(line)
        if macro:
            yield macro.groups()

def get_ecodes(headers):
    all_macros = []
    for header in headers:
        try:
            fh = open(header)
        except (IOError, OSError):
            continue
        all_macros += parse_header(fh)

    if not all_macros:
        sys.exit('no input macros found in: %s' % ' '.join(SOURCES))
    return all_macros

ECODES = get_ecodes(SOURCES)

context = {
    'UNAME': get_uname(),
    'CODES':   os.linesep.join('\t%s = %s' % i for i in ECODES),
    'CODEMAP': os.linesep.join('\t"%s": %s,' % (i, i) for i, _ in ECODES),
}

print(TEMPLATE.safe_substitute(context))
