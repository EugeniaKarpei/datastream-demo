# valery-datadog-datastream-demo

This project is a little demo of how I would approach handling metric charts with tags on them, doing filtering/partitioning/aggregations by tags and time periods and plotting charts similarly to how DataDog API web client works.

**Project structure**

|
|- cmd
|   |
|   |- server         // this is where server starter lives
|   |- testclient     // test client in Go, hitting locally started service API-s: /getData and /getFilters
|
|- data               // test dataset as a csv - some random online sales transactions for 2019. I like this
|                     // dataset becasuse it has trx dates and can be aggregated by time and few other fields
|                     // such as gender, location, product_catefory, coupon_status, coupon_code
|
|- frontend           // react project with all frontend code
|
|- internal           // backend, the most interesting part
|   |
|   |- api            // API endpoint handlers
|   |- config         // mostly some metadata related to csv dataset parsing
|   |- data           // data model (api, metrics, tags) + csv file reader
|   |- processor      // core of metric processing:
|                     // * _MetricProcessor_ with all internal datastructures supporting filtering by tags
|                     // * _partitioners_ to partition data into time-chunks and prepare for aggregation
|                     // * _aggregators_ to aggregate the data
|                     // * _triesearch_ - small Trie-based datastructure to search for tag names and values
|
|- Dockerfile         // Deployment file.

The backend is implemented as a web-socket service that potentially is able to handle real-time data streams. I used [gorilla/websocket](https://github.com/gorilla/websocket)https://github.com/gorilla/websocket as a server and [GIN](https://github.com/gin-gonic/gin)https://github.com/gin-gonic/gin as http router.

The frontend is a simple 1-page react app with chart component and few input fields.

**How metrics retrieval by tags works**

In order to provide fast metric lookup by tag:value pair we pre-compute tagged metrics as a nested map of

   map [tagName]   ->   map [tagValue]   ->   Metrics

where Metrics is essentially a collection of MetricRecords that is able to perform 2 things:
1. Quickly [O(1)] answer if certain MetricRecord is present in the collection (by metric record id) - this is important for filtering 
2. List all metric records inside the collection - used for metric data extraction for partitioning by time and aggregation

When _GetData_ request with filter comes to MetricProcessor we can face one of 3 possbile situations:
1. No filtering required (empty filters) - we use non-filtered metrics.
2. One filter is selected - we can get filtered data using nested map at O(1).
3. Two or more filters are selected - in this case we can still get data for each filter at O(1) but then we need to merge the results to find the intersection between all filters. This step potentially has linear complexity in this case and the alternative to it is to either precompute metrics with combined filters (similarly to how we did it for single filters but using compund keys (filter1:value1;filter2:value2;etc) or caching most popular combinations using same compound keys. I didn't implement any of these approaches here as it seems a little over-complicated for this demo.

After we gathered all metric points we do partitioning by time. The demo supports 3 scales of data aggregation granularity:
* Monthly
* Weekly
* Daily
Those are static scales, we could also use some dynamic partititoner here and calculate partition time period dynamically based on time interval that we are observing, but I decided to have time interval hardcoded (2019-01-01..2019-12-31) for this demo and therefore - used static time partitions.
Also since our data at this point (metric records) is likely sorted by time - time partitioning is the step that could potentially benefit from parallelisation.

Final step is aggregation of time-partitioned data. We support several aggregating functions:
* Sum
* Count
* Avg
They are pretty straightforward in this demo, but this is another step that can be parallelised.

**How tag names search works**

We are using a Trie data structure to store all tag name:value pairs in it to be able to quickly get list of tags and values by given name:value prefix.

**TODOs**
* P0: Finish frontend
* P1: Implement wildcards for tag search API
* P1: Write unit-tests
* P2: Add more metrics, add metric name to API.
* P2: Add multiple metrics on the same chart
* P2: Add partition_by tagName