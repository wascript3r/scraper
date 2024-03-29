# python 3

import logging
from logging import FileHandler
import os
import re
# from logging.handlers import file

class LogFilter(object):
    def __init__(self, level):
        self.__level = level
    def filter(self, logRecord):
        return logRecord.levelno == self.__level

# variables
exception_file = "./log/exception.log"
info_file = "./log/info.log"
# remove old files
os.remove(exception_file)
os.remove(info_file)
MAX_BYTES = 50000
BACKUP_COUNT = 5
# formatters
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')

main_logger = logging.getLogger(__name__)    # sets main logger
main_logger.setLevel(logging.INFO)    # sets main logger
# info logger
# info_logging = RotatingFileHandler(info_file, maxBytes=MAX_BYTES, backupCount=BACKUP_COUNT)
info_logging = FileHandler(info_file)
info_logging.setLevel(logging.INFO)
info_logging.setFormatter(formatter)
info_logging.addFilter(LogFilter(logging.INFO))
main_logger.addHandler(info_logging)

# warning logger
# exception_logging = RotatingFileHandler(exception_file, maxBytes=MAX_BYTES, backupCount=BACKUP_COUNT)
exception_logging = FileHandler(exception_file)
exception_logging.setLevel(logging.ERROR)
exception_logging.setFormatter(formatter)
exception_logging.addFilter(LogFilter(logging.ERROR))
main_logger.addHandler(exception_logging)

try:
    import json_functions as jf     # import samples queries
    import aiohttp                  # import requests library - used in requesting webpage content
    from aiohttp import ClientSession
    import json                     # to parse and create json requests
    from bs4 import BeautifulSoup   # import beautifulsoup - to get info from page
    import gc                       # garbage collector
    import asyncio
    import requests
except Exception as er:
    main_logger.exception(f"Cant import module, exiting, error: {er}")
    import sys
    sys.exit()

class ShopItem():
    def __init__(self, item):
        self.search_query_id = item.split("|")[2]
        self.url = item.split("|")[0]  # done
        self.text = item.split("|")[1]    # done
        self.item_id = None # done
        self.price = ""   # done
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
    async def get_main_soup(self, url):
        try:
            async with ClientSession() as sesion:
                async with sesion.get(self.url) as response:
                    text = await response.read()
                return BeautifulSoup(text.decode('utf-8'), "lxml")
        except Exception as er:
            main_logger.exception(f"Can't download page, skipping this url: {url}")
            return "get_main_soup_download"
    async def collect_info(self):
        listing_page_soup = await self.get_main_soup(self.url)
        if listing_page_soup == "get_main_soup_download":
            main_logger.info("Can't download main listing page info")
            return
        # print(f"---getting info for {self.url[:40]}")
        # price and currency
        try:
            # item_id
            item_id = self.url.split("/")[-1]
            self.item_id = item_id[:12]
            scraped_price = listing_page_soup.select_one("span[itemprop=price]")
            # price and currency
            if scraped_price != None:
                scraped_price = scraped_price.text.split(" ", 1)
                currency = scraped_price[0]
                if currency == "US":
                    self.curreny = "USD"
                else:
                    self.curreny = scraped_price[0]
                price = scraped_price[1]
            price = price.replace(",", "")
            trim = re.compile(r'[^\d.,]+')
            self.price = trim.sub("", price)

            # remaining
            remaining_scraped = listing_page_soup.select_one("span#qtySubTxt span")
            if remaining_scraped != None:
                self.remaining = remaining_scraped.text.strip()

            # if keywords with abstract remaining - skip those f
            if "lots" in self.remaining:
                # print("lots in remaining, skipping")
                return
            if self.remaining == "More than 10 available":
                # print(f"Skipping: {self.item_id} - More than 10 available")
                return

            if self.remaining is None:
                self.remaining = 1
            else:
                self.remaining = self.remaining.replace("Last one", "1")
                self.remaining = self.remaining.replace(" available", "")
            # condition
            condition = listing_page_soup.select_one("div[itemprop=itemCondition]")
            if condition != None:
                if condition.text == "--not specified":
                    self.condition = "not specified"
                elif "New" in condition.text or "new" in condition.text:
                    self.condition = "new"
                elif "Used" in condition.text or "used" in condition.text:
                    self.condition = "used"
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
                    if shipping == "remove_unecessary_remove":
                        main_logger.info(f"collect_info: Couldn't remove unecessary info, url: {self.url}")
                    shipping = shipping.strip()
                    self.shipping_to.append(shipping)
            # photo_links
            photo_links_scraped = listing_page_soup.select("td.tdThumb div img")
            for photo in photo_links_scraped:
                photo = photo["src"].replace("64", "500")
                self.photo_links.append(photo)
            if not self.photo_links:
                photo = listing_page_soup.select_one("img[itemprop=image]")
                photo = photo["src"].replace("96", "500")
                self.photo_links.append(photo)
        except Exception as er:
            main_logger.exception(f"Couldn't get main info, skipping this url: {self.url}")
            return "collect_info_main"

        # checking if has sold history
        try:
            solds_link = listing_page_soup.select_one("a.vi-txt-underline")
            if solds_link != None:
                self.is_secondary_info = True
                self.sold_history_url = solds_link["href"]
                await self.get_secondary_info()
            await self.export_json()
        except Exception as er:
            main_logger.exception(f"Couldn't find if it has secondary info, skipping this url: {self.url}")
            return "collect_info_secondary"
        
    async def get_secondary_info(self):
        '''
        Gets information from item's sold history, such as user's id, price, quantity and date when this item was sold to that specific user
        '''
        try:
            sold_history = requests.get(self.sold_history_url)
            sold_history_soup = BeautifulSoup(sold_history.content, "html.parser")
        except Exception as er:
            main_logger.exception(f"Couldn't download sold_history page, skipping this secondary url: {self.sold_history_url}")
            
        try:
            if sold_history_soup.select_one("#splashRCP_ct") != None:
                self.secondary_info_headers.append("CAPTCHA")
                self.secondary_info_user.append("CAPTCHA")
                self.secondary_info_other.append("CAPTCHA")
                return True
            # headers
            correct_table = sold_history_soup.select_one(".BHbidSecBorderGrey")
            if correct_table != None:
                headers = correct_table.select(".tabHeadDesignFont")
                for header in headers:
                    header_text = header.text
                    header_text = self.remove_unecessary(header_text)
                    if header_text == "remove_unecessary_remove":
                        main_logger.info(f"get_secondary_info: Couldn't remove unecessary info, url: {self.url}")
                    self.secondary_info_headers.append(header_text)
            else:
                headers = sold_history_soup.select(".tabHeadDesignFont")
                for header in headers:
                    header_text = header.text
                    header_text = self.remove_unecessary(header_text)
                    if header_text == "remove_unecessary_remove":
                        main_logger.info(f"get_secondary_info: Couldn't remove unecessary info, url: {self.url}")
                    self.secondary_info_headers.append(header_text)
            
            # user id
            user_id = sold_history_soup.select("td.onheadNav")
            if user_id != None:
                for user in user_id:
                    user = user.text.strip()[:5]
                    self.secondary_info_user.append(user)

            # other info
            other_info = sold_history_soup.select("td.contentValueFont")
            counter = 0
            for count, info in enumerate(other_info):
                if info.text == "Accepted" or info.text == "Declined" or info.text == "Expired":
                    counter += 1
                else:
                    self.secondary_info_other.append(info.text)
            self.secondary_info_user = self.secondary_info_user[:-counter or None]
        except Exception as er:
            main_logger.exception(f"Couldn't get sold history, skipping for this url: {self.sold_history_url}")
            return "get_secondary_info_sold"

    def remove_unecessary(self, from_which):
        '''
        Removes all unecessary symbols from string.
        '''
        try:
            from_which = from_which.replace("\r", "")
            from_which = from_which.replace("\n", "")
            from_which = from_which.replace("\t", "")
            from_which = from_which.replace("\xa0", "")
            from_which = from_which.replace("|", "")
            from_which = from_which.replace("See exclusions", "")
            from_which = from_which.replace("See details", "")
            return from_which
        except Exception as er:
            main_logger.exception("Couldn't remove unecessasy info from string, skipping")

    async def export_json(self):
        try:
            registered = jf.main_info(self.item_id, self.search_query_id, self.text, self.curreny, self.condition, self.seller_id, self.photo_links, self.country, self.region, self.shipping_to)
            jf.listing_history(self.item_id, self.price, self.remaining)
            if not self.secondary_info_user or self.secondary_info_user[0] == "CAPTCHA":
                pass
                # print(f"Hehe, cant upload {self.item_id}")
            else:
                jf.sold_history(self.item_id, self.secondary_info_user, self.secondary_info_other, self.price, self.remaining)
        except Exception as er:
            import traceback
            traceback.print_exc()
            # print(f"Couldn't send json to server, item_id: {self.item_id} er: {er}")
            main_logger.exception(f"Couldn't send json to server, er: {er}")
            return "export_json_json"

def main_tasker(url, search_id):
    links = []  # where to store item listing urls
    try:
        search_page = requests.get(url)
    except Exception as er:
        main_logger.exception("Couldn't get search page content, skipping this url")
        return "main_tasker_search_content"
    if search_page.status_code == 200:
        try:
            search_page_parse = BeautifulSoup(search_page.content, "html.parser")
            listing_links = search_page_parse.select("a.s-item__link")  # get all item listing in search page
        except Exception as er:
            main_logger.exception("Couldn't convert page to soup object")
            return "main_tasker_soup"
        try:
            for count, link in enumerate(listing_links):
                url = link["href"]

                text = link.text.replace("New Listing", "")
                text = text.strip()
                info = f"{url}|{text}|{search_id}"
                links.append(info)
        except Exception as er:
            main_logger.exception("Couldn't enumerate through links, skipping this url")
            return "main_tasker_links"
        tasks = []
        loop = asyncio.get_event_loop()
        try:
            for link in range(len(links)):
                current_item = ShopItem(links[link])
                if link == len(links) - 1:
                    loop.run_until_complete(current_item.collect_info())
                loop.create_task(current_item.collect_info())
        except Exception as er:
            main_logger.exception(f"Couldn't create jobs and execute them, skipping this: {links[link]}")
    return 1

def main():
    main_logger.info("Starting program")
    print("starting")
    request_to_parse = jf.get_request()     # get JSON query and convert to python data types
    data_to_parse = request_to_parse["data"]
    from timeit import default_timer as timer
    start = timer()
    for count, request in enumerate(data_to_parse["queries"]):
        # print(f"{request['id']} -> {request['url']}")
        main_logger.info(f"Scraping search query - id: {request['id']}, url: {request['url']}")
        main_tasker(request["url"], request['id'])
    end = timer()
    main_logger.info(f"Ending program, elapsed time: {end-start}")

if __name__ == "__main__":
    main()