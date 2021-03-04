# python 3

import sample_queries as sq # import samples queries
import requests # import requests library - used in requesting webpage content
import json # to parse and create json requests
from bs4 import BeautifulSoup   # import beautifulsoup - to get info from page
import extra_functions as ef    # extra functions, used only in this sprint
import gc   # garbage collector
import logging  # do I need to tell anything more on this one??

# variables
logging.basicConfig(level=logging.INFO)
main_logger = logging.getLogger("main")

class ShopItem():
    def __init__(self, url, text, id):
        self.search_query_id = id
        self.url = url  # done
        self.text = text    # done
        self.item_id = None # done
        self.price = None   # done
        self.curreny = None #done
        self.remaining = None   # done
        self.condition = None   # done
        self.seller_id = None   # done
        self.country = None     # done
        self.region = list()    # done
        self.shipping_to = list()   # done
        self.photo_links = list()    # done
        self.sold_history_url = None           # done
        self.is_secondary_info = False
        self.secondary_info_headers = list()    # first item is always user id, then it depends
        self.secondary_info_user = list()
        self.secondary_info_other = list()
    def collect_info(self):
        listing_page = requests.get(self.url)
        listing_page_soup = BeautifulSoup(listing_page.content, "html.parser")
        # price and currency
        scraped_price = listing_page_soup.select_one("span[itemprop=price]")
        if scraped_price != None:
            scraped_price = scraped_price.text.split(" ", 1)
            self.curreny = scraped_price[0]
            self.price = scraped_price[1]

        # remaining
        remaining_scraped = listing_page_soup.select_one("span#qtySubTxt span")
        if remaining_scraped != None:
            self.remaining = remaining_scraped.text.strip()

        # condition
        condition = listing_page_soup.select_one("div[itemprop=itemCondition]")
        if condition != None:
            self.condition = condition.text

        # seller_id
        seller_id = listing_page_soup.select_one("span.mbg-nw")
        if seller_id != None:
            self.seller_id = seller_id.text

        # location
        locations = listing_page_soup.select_one("span[itemprop=availableAtOrFrom]")
        if locations != None:
            locations = locations.text.split(", ")
            locations_length = len(locations)
            self.country = locations[-1]
            for count, loc in enumerate(locations):
                if count == locations_length - 1:
                    break
                self.region.append(loc)

        # shipping to
        shipping_to_scraped = listing_page_soup.select_one("span[itemprop=areaServed]")
        if shipping_to_scraped != None:
            shipping_to_scraped = shipping_to_scraped.text.split(", ")
            # remove what's not needed (white spaced, new lines etc)
            for shipping in shipping_to_scraped:
                shipping = self.remove_unecessary(shipping)
                shipping = shipping.strip()
                self.shipping_to.append(shipping)

        # item_id
        item_id = self.url.split("/")[-1]
        self.item_id = item_id[:12]

        # photo_links
        photo_links_scraped = listing_page_soup.select("td.tdThumb div img")
        for photo in photo_links_scraped:
            photo = photo["src"].replace("64", "500")
            self.photo_links.append(photo)

        # checking if has sold history
        solds_link = listing_page_soup.select_one("a.vi-txt-underline")
        if solds_link != None:
            self.is_secondary_info = True
            self.sold_history_url = solds_link["href"]
            self.get_secondary_info()
        
    def get_secondary_info(self):
        '''
        Gets information from item's sold history, such as user's id, price, quantity and date when this item was sold to that specific user
        '''
        print(f"{self.text} {self.url}")
        sold_history = requests.get(self.sold_history_url)
        sold_history_soup = BeautifulSoup(sold_history.content, "html.parser")

        # headers
        correct_table = sold_history_soup.select_one(".BHbidSecBorderGrey")
        if correct_table != None:
            headers = correct_table.select(".tabHeadDesignFont")
            for header in headers:
                header_text = header.text
                header_text = self.remove_unecessary(header_text)
                self.secondary_info_headers.append(header_text)
        else:
            headers = sold_history_soup.select(".tabHeadDesignFont")
            for header in headers:
                header_text = header.text
                header_text = self.remove_unecessary(header_text)
                self.secondary_info_headers.append(header_text)
        
        # user id
        user_id = sold_history_soup.select("td.onheadNav")
        if user_id != None:
            for user in user_id:
                user = user.text.strip()[:5]
                self.secondary_info_user.append(user)

        # print(f"users before end deletion: {self.secondary_info_user}")
        # other info
        other_info = sold_history_soup.select("td.contentValueFont")
        for count, info in enumerate(other_info):
            if info.text == "Accepted" or info.text == "Declined" or info.text == "Expired":
                how_many_to_leave = int(count / (len(self.secondary_info_headers) - 1))
                self.secondary_info_user = self.secondary_info_user[:how_many_to_leave]
                # print(f"to_delete - {how_many_to_leave}, count - {count}, len(headers): {len(self.secondary_info_headers)}")
                break
            self.secondary_info_other.append(info.text)

        # print(f"users after end deletion: {self.secondary_info_user}")

    def remove_unecessary(self, from_which):
        '''
        Removes all unecessary symbols from string.
        '''
        from_which = from_which.replace("\r", "")
        from_which = from_which.replace("\n", "")
        from_which = from_which.replace("\t", "")
        from_which = from_which.replace("\xa0", "")
        from_which = from_which.replace("|", "")
        from_which = from_which.replace("See exclusions", "")
        from_which = from_which.replace("See details", "")
        return from_which

    def export_json(self):
        with open("register_listing.json", "a", encoding="utf8") as writer:
            ef.json_register_listing(writer, self.item_id, self.search_query_id, self.text, self.curreny, self.condition, self.seller_id, self.photo_links, self.country, self.region, self.shipping_to)
        if self.is_secondary_info == True:
            with open("add_listing_history.json", "a", encoding="utf8") as writer:
                ef.json_add_sold_history(writer, self.item_id, self.secondary_info_user, self.secondary_info_other, self.price, self.remaining)

def get_search_links(url, search_id):
    search_page = requests.get(url)
    if search_page.status_code == 200:
        search_page_parse = BeautifulSoup(search_page.content, "html.parser")
        listing_links = search_page_parse.select("a.s-item__link")  # get all item listing in search page
        for count, link in enumerate(listing_links):
            url = link["href"]
            text = link.text
            # main_logger.info(f"text: {text}, url: {url[:30]}")
            print(f"id: {str(search_id)}")
            current_item = ShopItem(url, text, search_id)
            current_item.collect_info()
            current_item.export_json()
            print(f"count: {count}, gc: {gc.get_count()}")

def main():
    print(gc.get_count())
    ef.empty_file("register_listing.json")     # empty file in which to save JSON
    ef.empty_file("add_listing_history.json")   # same as above
    request_to_parse = json.loads(sq.sample_request_query_json)     # get JSON query and convert to python data types
    data_to_parse = request_to_parse["data"]
    print(gc.get_count())
    # print(data_to_parse)
    for count, request in enumerate(data_to_parse["queries"]):
        main_logger.info(f"id: {request['id']}, url: {request['url']}")
        get_search_links(request["url"], request['id'])
    print("")
    print(gc.get_count())

if __name__ == "__main__":
    main()