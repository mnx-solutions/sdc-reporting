# sdc-reporting
Reporting and billing for sdc

# Design Concept
https://github.com/joyent/sdc-hagfish-watcher runs on each of the CN's, and outputs usage reports to /var/log/usage every hour.  `sdc-reporting` will read these files, summarize and send them to the reporting database to be used for billing and reporting.

# TODO
- Automate the reading, summary, and output of the usage files
- Write a microservice to serve the data for `sdc-portal` to show billing details in the portal.
- Integrate into the Chargebee invoicing process

# hagfish watcher 
Details on hagfish can be founder here: https://eng.joyent.com/usage/
