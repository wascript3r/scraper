import json

sample_request_query = {
    "error": None,
    "data": {
        "queries": [
            {
                "id": 1,
                "url": "https://www.ebay.com/sch/i.html?_from=R40&_trksid=p992.m570.l1313&_nkw=bitcoin&_sacat=0"
            },
            {
                "id": 2,
                "url": "https://www.ebay.com/sch/i.html?_from=R40&_trksid=p992.m570.l1313&_nkw=litecoin&_sacat=0"
            },
            {
                "id": 3,
                "url": "https://www.ebay.com/sch/i.html?_from=R40&_trksid=p992.m570.l1313&_nkw=dogecoin&_sacat=0"
            }
        ]
    }
}

sample_request_query_json = json.dumps(sample_request_query)