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

# Installation

## hagfish-reader

hagfish-reader needs to be installed on all compute nodes (CNs). This is accomplished by a cron entry:

```
# crontab -l
10 * * * * /opt/local/bin/run_hagfish-reader.sh 2>&1 >> /var/log/hagfish-reader.log
```

Ensure `/opt/local/bin/run_hagfish-reader.sh`, `/opt/local/bin/hagfish-reader`, and `/opt/local/etc/.env` are installed and configure.

Review the headnode for configuration details on `.env`:

Though it should include:

```
export DB_NAME="reporting"
export DB_PASS=""
export DB_USER="reporting"
export DB_HOST=""
export DB_PORT="3306"
```

Optionally, on a node that also updates billing packages you will need entries for chargebee
```
export CHARGEBEE_KEY=""
export CHARGEBEE_SITE=""
```

