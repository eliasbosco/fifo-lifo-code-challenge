# unicorn

* it takes some time until a unicorn is produced, the request is blocked on requesters site and he needs to wait 

* to improve the situation adjust the code, so that the requester is receiving a request-id, with this request-id he can poll and validate if unicorns are produced

* if the unicorn is produced it should be returned though using fifo principle

* adjust the code, so that every x seconds a new unicorn is produced at put to a store, which can be used to fulfill the requests (LIFO Store)

* make sure, duplicate capabilities are not added to the unicorn

* improve the overall code

* if any requirements are not clear, compile meaningful assumptions

# code challenge outcome
* How-to *Install* and *Run* via Docker
```
docker build -t unicorn .
docker run -d -p 8888 -n unicorn unicorn
```

* Endpoints available:
```
# Request production of the amount unicorn informed in the body request
POST /api/set-unicorns-production

# Retrieve detail from the requested unicorn production
GET /api/get-request-detail/:request_id

# List all the unicorns production in queue
GET /api/get-all-request

# Delivery unicorn package fulfilling FIFO principle
PUT /api/delivery-package/:request_id

# Clean up queue
DELETE /api/clean-queue
```

> Remarks: Postman collection file ***Code Challenges.postman_collection.json*** follows in the project