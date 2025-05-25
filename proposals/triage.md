# Triage

The process of triaging bugs getting issues ready for processing

## Basic settings

The basic triage process is to set parameters around each issue:

1. Type:
   1. Process
   1. Research
   1. Code

1. Priorities
   1. P0 - Work on straight away
   1. P1 - Urgent, but not P0
   1. P2 - Not Urgent

### Process

Something to do outside of the system, these are things that we do
outside of the system. Broken down further into

1. 5 minutes [Small]
1. 30 minutes [Medium]
1. 2 hours [Large]
1. > 2 hours [XL]

These are then bin packed into spare time

### Research

These are issues that lead to other issues - stuff where we need
to understand something or look into something. This is at least
30 minutes, most likely longer

1. 2 hours [Large]
1. > 2 hours [XL]

### Code

We only code in transit. These are first in, first out type affairs
and code should be pre-broken down into workable tasks. Code has
to be code, not things that require internet. Effectively we should be
able to define a unit test at the outset.

Code issues are blessed when ready to run. If not blessed, then we
break it down further.

## Post Triage

Triage is interactive and we do it once a day, at the start of the day
and each bug needs to be worked through

## Process

1. Fokus Triage pulls up each untriaged issue in turn.
1. Fokus Triage supports type labels
1. Fokus Triage supports priority labels
1. Fokus Triage supports size labels
1. For code issues, we validate that the code issue is suitable and is blessed.

## Tasks

1. Fokus Triage prints an untriaged issue
1. Fokus sets type label
1. Fokus sets sublabels
1. Fokus loops
