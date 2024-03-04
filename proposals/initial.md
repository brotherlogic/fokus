<!-----



Conversion time: 0.331 seconds.


Using this Markdown file:

1. Paste this output into your source file.
2. See the notes and action items below regarding this conversion run.
3. Check the rendered output (headings, lists, code blocks, tables) for proper
   formatting and use a linkchecker before you publish this page.

Conversion notes:

* Docs to Markdown version 1.0β35
* Sun Mar 03 2024 16:00:38 GMT-0800 (PST)
* Source doc: Fokus
----->



# Fokus

<p style="text-align: right">
Brotherlogic</p>


<p style="text-align: right">
2024-02-27</p>


<p style="text-align: right">
Draft</p>



### Abstract

Rewrite the local focus into something more Kubernetes friendly


### Rationale

Fokus is a method for determining which github issue I should focus on at any given time - using whatever signals it has available to it (i.e. what time it is, where I am, what the state of the house is etc.).

Effectively this works through pluggable modules which (a) declare a weight indicating their preference to be active, and also a selection module which when given a set of github issues makes a selection, returning both the issue and some metadata about the module that selected it. Oldest issue resolve ties in the case of equal weighting


### Auth

This will be publicly reachable, so basic token based authentication, token is placed in kubernetes secrets and validated within the context of the caller. Alerts setup for failures


### Modules



1. Not getting to it Module

    If any issue is over 2 weeks old, flag it here. Weight is 0.95

2. Hacking Module

    Between the hours of 1600 and 1800, returns oldest issue in a subset of components (e.g. gramophile, githubridge etc.). Weight is 0.9

3. Cleaning Module

    Active if the cleaner is switched on. Returns an issue in the recordcleaner component which has “needs cleaning in it”. If the top issue is labelled “in-progress”, returns nothing. Weight is 1.0.



### Build Process



1. Create the repo
2. Add this proposal
3. Create skeleton of server, with build settings
4. Configures server in flux
5. Define proto for fokus
6. Build out base API with no module selection
7. Write CLI - replaces focus_cli get
8. Build basic grafana dashboard
9. Alerting module in place to support auth failures
10. Not getting to it module is in place and running
11. Cleaning module is in place and running
12. Hacking module is in place and running