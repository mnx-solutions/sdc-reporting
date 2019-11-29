# sdc-reporting
Reporting and billing for sdc

# Design Concept
https://github.com/joyent/sdc-hagfish-watcher runs on each of the CN's, and outputs usage reports to /var/log/usage every hour.  `sdc-reporting` will read these files, summarize and send them to the reporting database to be used for billing and reporting.

# TODO
- Automate the reading, summary, and output of the usage files
- Write a microservice to serve the data for `sdc-portal` to show billing details in the portal.
- Integrate into the Chargebee invoicing process

# hagfish watcher 
Details on hagfish can be founder here: https://eng.joyent.com/usage/.

# update-billing
Currently reads billing information from chargebee and inserts into the billing table.  Pricing is based on monthly prices.

# hagfish-reader
The hagfish-reader reads the usage logs in /var/log/usage/ and insert parts of the records into a database to be used with hagfish-reporter

Currently we import owner_uuid, billing_id, vm_uuid.  From this data we read a billing table to calculate per minute pricing

# hagfish-reporter
hagfish-reporter is used to calculate costs, and network bandwidth usage.
