import json
import requests

header_authorization = {'Authorization': 'Bearer oneltsecret'}

def get_request():
    '''
    Get request, which items need to be scraped
    '''
    response = requests.get("http://91.225.104.238:3000/api/queries/get", headers=header_authorization)
    data = response.json()
    return data

def check_if_exists(item):
    '''
    Checks if item already exists in database
    '''
    print(item)
    query = {"id": item}
    response = requests.post("http://91.225.104.238:3000/api/listing/exists", json=query, headers=header_authorization)
    data = response.json()
    return data

def main_info(id, search_query, title, currency, condition, seller_id, photo_links, country, region, shipping_to):
    if_exists = check_if_exists(id)
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
                "id": id,
                "searchQueryID": int(search_query),
                "title": title,
                "currency": currency,
                "condition": condition,
                "sellerID": seller_id,
                "photos": photo_links,
                "location": locations_array,
                "shipping": shipping_array,
            }
            print(query)
            response = requests.post("http://91.225.104.238:3000/api/listing/register", json=query, headers=header_authorization)
            data = response.json()
            error = data["error"]
            print(f"Response: {error}")

def sold_history():
    pass

if __name__ == "__main__":
    import main
    main.main()