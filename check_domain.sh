#!/usr/bin/env bash

domain=$1
whois $domain | grep -m 1 -Eo '(Expir[a-z]* Date: [0-9-]*|No match)' | sed 's/.*: //'
