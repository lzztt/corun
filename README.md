# corun
a parallel executor with a maximum number of concurrent tasks at any time

## A task
Imagine that we have a task, e.g. checking whether a domain name is available for register. We can run the following command
```bash
whois test.com
```
We may end up with an one-liner with filtered output:
```bash
$ whois test.com | grep -m 1 -Eo '(Expir[a-z]* Date: [0-9-]*|No match)' | sed 's/.*: //'
2019-06-17
$ whois test_non_exist.com | grep -m 1 -Eo '(Expir[a-z]* Date: [0-9-]*|No match)' | sed 's/.*: //'
No match
```
The above shows that test.com will be expired on 2019-06-17 and test_non_exist.com is available.

## What if we have a large number of similar tasks

We can put the one-liner into [a bash script](check_domain.sh) and loop through the input data.

### The execution time of one task
```bash
$ time bash check_domain.sh test.com
2019-06-17

real    0m6.601s
user    0m0.004s
sys     0m0.006s
$ time bash check_domain.sh test_non_exist.com
No match

real    0m3.275s
user    0m0.004s
sys     0m0.004s
```

If we have 1 million of these tasks, to loop through all tasks, we need 3.3 to 6.6 million seconds, which is 38 to 76 days.

## We need to parallel these tasks, but without DDOS ourselves.
```bash
$ ./corun --in=input.1000.txt --np=100 --out=output.1000.txt
```

The input file contains task identifiers and commands.
```bash
$ head input.1000.txt 
able_bay: bash check_domain.sh ablebay.com
about_bay: bash check_domain.sh aboutbay.com
above_bay: bash check_domain.sh abovebay.com
act_bay: bash check_domain.sh actbay.com
add_bay: bash check_domain.sh addbay.com
after_bay: bash check_domain.sh afterbay.com
again_bay: bash check_domain.sh againbay.com
age_bay: bash check_domain.sh agebay.com
ahead_bay: bash check_domain.sh aheadbay.com
air_bay: bash check_domain.sh airbay.com
```

The output file contains task identifiers and outputs. The order of outputs may be different from the order of input. Thus we need the task idenfier.
```bash
$ head output.1000.txt 
buy_bay: 2018-09-30
apply_bay: 2027-02-06
claim_bay: 2019-05-13
dark_bay: 2019-07-09
candy_bay: 2018-11-10
carry_bay: 2019-02-18
bus_bay: 2019-04-18
above_bay: No match
about_bay: No match
able_bay: 2019-03-08
```

## How much time it can save with --np=100?
When testing with 1000 example tasks and setting the maximum number of concurrent tasks equals to 100, the total execution time is 2.5 minutes. This includes the ramp-up and ramp-down time of the pipeline (the first 100 and last 100 tasks). It would take 1 to 2 hours if we run the tasks sequentially, based on the execution time of a single task (3.3 to 6.6 seconds).
```bash
ikki@dt:~/github/corun$ time ./corun --in=input.1000.txt --np=100 --out=/tmp/out1

real    2m36.791s
user    0m5.858s
sys     0m3.411s
```
