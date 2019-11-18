#! python

import json
import sys
import pprint
import time

import chargebee

billing_plan = {}

# read package pricing from chargebee
#chargebee.configure("","")
#entries = chargebee.Plan.list({
#    "limit": 100,
#    "status[is]" : "active"
#    })
#for entry in entries:
#  plan = entry.plan
#  billing_plan[plan.id] = {"name": plan.name, "price": plan.price}
#print billing_plan

# for testing, set static package pricing
# pricing in cents, for the month based on 720 hours. 
# to get per minute pricing; per_minute_price = price / 720 / 60
billing_plan = {u'6c1329e0-e983-ea40-8d2d-98f9ee035d3c': {'price': 14000, 'name': u'c1.large'}, u'd43202b9-38dd-4484-d92f-8e9273b5fc4a': {'price': 3600, 'name': u's1.nano'}, u'58d53657-6690-e40c-bb20-d0d280817f44': {'price': 1064, 'name': u'g1.micro'}, u'99830774-59e7-ef98-8725-b454907dd458': {'price': 1800, 'name': u's1.pico'}, u'74e3c2d3-7371-cb7d-c97d-afd5b7754739': {'price': 22144, 'name': u'm1.large'}, u'0d843596-c0eb-4e69-e0a4-9b94417777a1': {'price': 100000, 'name': u'c1.3xlarge'}, u'2820b176-173d-617b-bb46-a8a8567c729c': {'price': 6736, 'name': u'm1.small-100'}, u'87b634f3-18c9-c989-f710-907dc06d39a4': {'price': 130592, 'name': u'g1.4xlarge'}, u'd156bfa4-072d-6ea8-de88-a526b270605b': {'price': 532, 'name': u'g1.nano'}, u'e16ef682-a166-ce8c-ee8e-bb0ddf02b033': {'price': 71296, 'name': u'g1.2xlarge'}, u'35f9dad9-77aa-e433-be13-83589627063d': {'price': 56000, 'name': u'c1.2xlarge'}, u'192f8b31-3fb0-e527-c665-d72d59276a69': {'price': 700, 'name': u's1.pico-ops'}, u'2b810a58-3d2f-6261-801f-8ff6065a5f6f': {'price': 8912, 'name': u'g1.medium'}, u'483b4cd6-90ca-61cf-ebd9-e1d5661c524a': {'price': 88576, 'name': u'm1.2xlarge'}, u'df26ba1d-1261-6fc1-b35c-f1b390bc06ff': {'price': 1750, 'name': u'c1.xsmall'}, u'c35ba654-678c-695d-a43f-d381260a7d70': {'price': 11072, 'name': u'm1.medium'}, u'dfce24f9-87cb-637a-a5b7-da5fbcd631bc': {'price': 250, 'name': u'g1.pico'}, u'cloud-hosting': {'price': 0, 'name': u'Cloud Hosting'}, u'f945cabf-2bbb-e19d-97ab-a9f8c1ee44be': {'price': 44288, 'name': u'm1.xlarge'}, u'3259329f-bcf0-6449-d64f-ecebb73030d5': {'price': 28000, 'name': u'c1.xlarge'}, u'eef3aee6-ca19-ea80-c905-821870ed0832': {'price': 17824, 'name': u'g1.large'}, u'893a7241-16f9-ef0a-9f8d-ee1a0a29328e': {'price': 144000, 'name': u'c1.4xlarge'}, u'75b134f9-ffe0-e463-ddc2-b1ffa0f35738': {'price': 3500, 'name': u'c1.small'}, u'6cfc63f1-f295-64a4-9973-f2d6747945c6': {'price': 2228, 'name': u'g1.xsmall'}, u'818cb60d-3af8-4918-9673-a013aafd692d': {'price': 35648, 'name': u'g1.xlarge'}, u'a731a718-5f8e-e1ea-a400-ab6bad4dcde7': {'price': 5536, 'name': u'm1.small'}, u'09664b0f-34d6-4eb3-cae2-cf275c3f3ba4': {'price': 165152, 'name': u'm1.4xlarge'}, u'9ffb1543-d7aa-6fe0-fa4b-b1e70beb3dca': {'price': 4456, 'name': u'g1.small'}, u'65689ae7-965f-6975-e38e-bc39b6f4b6d1': {'price': 7000, 'name': u'c1.medium'}}

billing_dict = {}

with open(sys.argv[1], 'r') as fp:
    line = fp.readline()
    cnt = 1
    while line:
        try:
            json_line = json.loads(line)
        except Exception, e:
            pass
        else:
            try:
                line_type = json_line['type']
            except:
                pprint.pprint(json_line)

            if line_type == 'summary':
               vm_count = json_line['vm_count'] 

            if line_type == 'usage':
                vm_uuid = json_line['uuid']
                owner_uuid = json_line['config']['attributes']['owner-uuid']
                billing_id = json_line['config']['attributes']['billing-id']
                timestamp = json_line['timestamp']

                # u'network_usage': {u'net0': {u'counter_start': u'2019-06-21T05:13:09.864Z',
                #              u'received_bytes': 12962028958,
                #              u'sent_bytes': 15514621839}},

                if not owner_uuid in billing_dict:
                    billing_dict[owner_uuid] = {}
               
                if not vm_uuid in billing_dict[owner_uuid]:
                    billing_dict[owner_uuid][vm_uuid] = []

                if billing_plan.get(billing_id, None):
                    price = float(billing_plan[billing_id]['price'])/720/60
                else:
                    price = 0
               
                billing_dict[owner_uuid][vm_uuid].append({'timestamp': timestamp, 'billing_id': billing_id, 'price': price})
            
        cnt += 1
        line = fp.readline()

for owner in billing_dict.keys():
    for vm in billing_dict[owner]:
        cost = sum([li['price'] for li in billing_dict[owner][vm]]) 
        print owner, vm, cost, billing_dict[owner][vm][-1]['timestamp']
