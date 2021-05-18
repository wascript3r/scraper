import json
import requests
import logging
from logging.handlers import RotatingFileHandler
from datetime import datetime
import re

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

# Current time
# cur = "21-04-22 10:06:34"
send_ids = []
history_ids = []

def get_request():
    '''
    Get request, which items need to be scraped
    '''
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
    if_exists = check_if_exists(item_id)
    if if_exists["error"] == None:
        data = if_exists["data"]
        if data["exists"] == False:
            shipping = {}
            shipping_array = []
            locations = {}
            locations_array = []
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
            print(f"main_info response: {response.json()}")


def listing_history(item_id, price, remaining):
    if item_id not in send_ids:
        send_ids.append(item_id)
        price = price.replace("$", "")
        now = datetime.now()
        cur = now.strftime("%Y-%m-%d %H:%M:%S")
        query = {
            "listingID": item_id,
            "price": float(price),
            "remainingQuantity": int(remaining),
            "parsedDate": cur
        }
        # query = json.dumps(sample)
        response = requests.post(requests_update_listing, json=query, headers=header_authorization)
        print(f"listing_history response: {response.json()}")

def sold_history(item_id, users, other_info, product_price, quantity):
    if item_id not in history_ids:
        sold_history_array = []
        count = 0
        for counts, user in enumerate(users):
            sold_history = {}
            # price parsing
            price = other_info[count]
            if price.isdigit():
                break
            
            price = price.replace(",", "")
            trim = re.compile(r'[^\d.,]+')
            price = trim.sub("", price)
            # price = price
            if price.lower().islower() or "--" in price or not price:
                count += 3
                continue

            price = float(price)
            # date parsing
            date_count = other_info[count + 2]
            date = date_count.split()[0]
            date = datetime.strptime(date, "%b-%d-%y")
            date = date.strftime("%Y-%m-%d")
            date = f"{date} {date_count.split()[1]}"
            sold_history["user"] = str(user)
            sold_history["price"] = price
            sold_history["quantity"] = int(other_info[count + 1])
            sold_history["date"] = date
            sold_history_array.append(sold_history)
            del sold_history
            count += 3
        query = {
            "listingID": str(item_id),
            "soldHistory": sold_history_array
        }
        with open("history.json", "a") as w:
            w.write(json.dumps(query, indent=4))
            w.write("\n")
        history_ids.append(item_id)
        response = requests.post(requests_update_history, json=query, headers=header_authorization)
        print(f"sold_history response: {response.json()}")

if __name__ == "__main__":
    import main
    main.main()
    