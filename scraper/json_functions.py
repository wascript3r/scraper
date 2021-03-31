import json
import requests
import logging
from logging.handlers import RotatingFileHandler

# variables
header_authorization = {'Authorization': 'Bearer oneltsecret'}
info_file = "log/info.log"
MAX_BYTES = 50000
BACKUP_COUNT = 5

request_get = "http://91.225.104.238:3000/api/queries/get"
request_exists = "http://91.225.104.238:3000/api/listing/exists"
request_register_main = "http://91.225.104.238:3000/api/listing/register"
# logging
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
main_logger = logging.getLogger(__name__)    # sets main logger
main_logger.setLevel(logging.INFO)    # sets main logger
info_logging = RotatingFileHandler(info_file, maxBytes=MAX_BYTES, backupCount=BACKUP_COUNT)
info_logging.setLevel(logging.INFO)
info_logging.setFormatter(formatter)
main_logger.addHandler(info_logging)

def get_request():
    '''
    Get request, which items need to be scraped
    '''
    main_logger.info(f"Sending request to {request_get}")
    response = requests.get(request_get, headers=header_authorization)
    data = response.json()
    return data

def check_if_exists(item):
    '''
    Checks if item already exists in database
    '''
    main_logger.info(f"Checking if item exists id: {item}")
    query = {"id": item}
    response = requests.post(request_exists, json=query, headers=header_authorization)
    data = response.json()
    return data

def main_info(item_id, search_query, title, currency, condition, seller_id, photo_links, country, region, shipping_to):
    main_logger.info(f"Registering {item_id} main info")
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
            error = data["error"]
            print(f"Response: {error}")
        elif data["exists"] == True:
            print(f"Preke jau egzistuoja duonbazeje - {item_id}")

def sold_history():
    pass

if __name__ == "__main__":
    import main
    main.main()