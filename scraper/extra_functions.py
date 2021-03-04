import json
import main

def json_register_listing(writer_object, id, search_query, title, currency, condition, seller_id, photo_links, country, region, shipping_to):
    register_listing = {
        "id": str(id),
        "searchQueryID": str(search_query),
        "title": str(title),
        "currency": str(currency),
        "condition": str(condition),
        "sellerID": str(seller_id),
        "photos": photo_links,
        "country": country,
        "region": region,
        "shipping": shipping_to        
    }

    register_listing_json = json.dumps(register_listing)
    writer_object.write(register_listing_json)
    writer_object.write("\n")

def json_add_sold_history(writer_object, item_id, history_users, history_other, price, quantity):
    # print(f"json users: {history_users}")
    sold_history = {}
    sold_history_array = []
    sec_count = 0
    for count, user in enumerate(history_users):
        sold_history["user"] = str(user)
        sold_history["price"] = history_other[sec_count]
        sold_history["quantity"] = history_other[sec_count + 1]
        sold_history["date"] = history_other[sec_count + 2]
        sold_history_array.append(sold_history)
        sec_count += 2
    add_listing_history = {
        "id": str(item_id),
        "soldHistory": sold_history_array,
        "currentPrice": price,
        "remainingQuantity": quantity
    }
    add_listing_history_json = json.dumps(add_listing_history)
    writer_object.write(add_listing_history_json)
    writer_object.write("\n")

def empty_file(output_name):
    with open(output_name, 'w') as writer:
        writer.write("")

if __name__ == "__main__":
    main.main()