import json
import requests
import logging
from logging.handlers import RotatingFileHandler
from datetime import datetime

# variables
header_authorization = {'Authorization': 'Bearer oneltsecret'}
info_file = "log/json_info.log"
MAX_BYTES = 50000
BACKUP_COUNT = 5

request_get = "http://91.225.104.238:3000/api/queries/get"
request_exists = "http://91.225.104.238:3000/api/listing/exists"
request_register_main = "http://91.225.104.238:3000/api/listing/register"
requests_update_listing = "http://91.225.104.238:3000/api/listing/history/add"
requests_update_history = "http://91.225.104.238:3000/api/listing/soldHistory/add"
# logging
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
main_logger = logging.getLogger(__name__)    # sets main logger
main_logger.setLevel(logging.INFO)    # sets main logger
info_logging = RotatingFileHandler(info_file, maxBytes=MAX_BYTES, backupCount=BACKUP_COUNT)
info_logging.setLevel(logging.INFO)
info_logging.setFormatter(formatter)
main_logger.addHandler(info_logging)

# Current time
# cur = "21-04-22 10:06:34"
send_ids = []

def get_request():
    '''
    Get request, which items need to be scraped
    '''
    # main_logger.info(f"Sending request to {request_get}")
    response = requests.get(request_get, headers=header_authorization)
    data = response.json()
    return data

def check_if_exists(item):
    '''
    Checks if item already exists in database
    '''
    # main_logger.info(f"Checking if item exists id: {item}")
    query = {"id": item}
    response = requests.post(request_exists, json=query, headers=header_authorization)
    data = response.json()
    return data

def main_info(item_id, search_query, title, currency, condition, seller_id, photo_links, country, region, shipping_to):
    # main_logger.info(f"Registering {item_id} main info")
    with open("test.txt", "a") as w:
        w.write(f"{item_id:13} - {title:150} - {currency:10} - {condition:10} - {seller_id:} - {photo_links} - {country} - {region} - {shipping_to} - {search_query }\n")
    if_exists = check_if_exists(item_id)
    if if_exists["error"] == None:
        data = if_exists["data"]
        if data["exists"] == False:
            shipping = {}
            shipping_array = []
            locations = {}
            locations_array = []
            print(f"{currency} - {condition}")
            for ship in shipping_to:
                shipping["country"] = ship
                shipping["region"] = None
                shipping_array.append(shipping)
            for reg in region:
                locations["country"] = country
                locations["region"] = reg
                locations_array.append(locations)
            query = {
                "id": item_id,
                "searchQueryID": int(search_query),
                "title": title,
                "currency": currency,
                "condition": condition,
                "sellerID": seller_id,
                "photos": photo_links,
                "location": locations_array,
                "shipping": shipping_array,
            }
            response = requests.post(request_register_main, json=query, headers=header_authorization)
            data = response.json()
            print(f"Registering response: {data}")
        elif data["exists"] == True:
            main_logger.info(f"Preke jau egzistuoja duonbazeje - {item_id}, atnaujinama info")



def listing_history(item_id, price, remaining):
    if item_id not in send_ids:
        send_ids.append(item_id)
        price = price.replace("$", "")
        now = datetime.now()
        cur = now.strftime("%Y-%m-%d %H:%M:%S")
        if remaining is None:
            remaining = 1
        else:
            remaining = remaining.replace(" available", "")
            remaining = remaining.replace("More than ", "")
            remaining = remaining.replace("Last one", "1")
            remaining = remaining.replace("3 lots (3 items per lot)", "3")
        # print(f"listing_history | {item_id} - {price} - {remaining} - {cur}")
        # main_logger.info(f"Atnaujinama {item_id} informacija")
        query = {
            "listingID": item_id,
            "price": float(price),
            "remainingQuantity": int(remaining),
            "parsedDate": cur
        }
        # query = json.dumps(sample)
        print(type(query))
        response = requests.post(requests_update_listing, json=query, headers=header_authorization)
        data = response.json()
        print(f"Listing history response: {data}")
    else:
        print(f"Trying to send clone of {item_id}")

def sold_history(item_id, users, other_info):
    print(f"sold_history | {item_id} - {users} - {other_info}")

if __name__ == "__main__":
    import main
    main.main()