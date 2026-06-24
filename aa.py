import requests
import time
from pypinyin import lazy_pinyin  # Pinyin conversion library

# Tencent Finance qt.gtimg.cn code format:
# - A-shares: sh / sz + 6-digit security code
# - HK stocks: hk + 5-digit HKEX code (e.g. Hua Hong H-share hk01347); different from A-share codes
#   Same group example: HK "Hua Hong Semiconductor" 01347 -> hk01347; STAR "Hua Hong" 688347 -> sh688347
# API response code field: A-shares use sh/sz+digits; HK stocks use digits only
STOCK_NAME_MAPPING = {
    "sh601138": "Foxconn Industrial Internet",
    "sz300985": "Zhiyuan New Energy",
    "sz300059": "East Money Information",
    "sz300601": "Kangtai Biologics",
    "sz000158": "Changshan Beiming",
    "sz000066": "China Great Wall",
    "sz300592": "HuaKai Ebuy",
    "sh600199": "Jin Seed Wine",
    "sh600398": "Hailan Home",
    "sh601006": "Daqin Railway",
    "sz300527": "ST Yingji",
    "sz002424": "ST Bailing",
    "sz301058": "COFCO Technology and Engineering",
    "sh563230": "Fuguo Satellite ETF",
    "sh588780": "Guolian Sci-Tech Chip Design ETF",
    "sh688521": "VeriSilicon",
    "sh688110": "Dosilicon",
    "sh688047": "Loongson",
    "sz000858": "wuliangye",
    "sh601668": "zhong guo jianzhu",
    "sh601166": "xingye Bank",
    "sh601818": "Everbright Bank",
    "sz300474": "Jingjia Micro",
    "sh688256": "Cambricon",
    "sh688041": "Hygon Information",
    "sh688981": "SMIC",
    "sz300498": "Wens Foodstuff",
    "sz300097": "Zhiyun",
    "sh603029": "Swan",
    "sz000564": "Gongxiao Daji",
    "sz003042": "Sino-Agri United",
    "sh600693": "Dongbai Group",
    "sz002251": "Bubugao",
    "sz000759": "Zhongbai Group",
    "sh601933": "Yonghui Superstores",
    "sh601012": "LONGi Green Energy",
    "sh600340": "China Fortune Land",
    "sz300476": "Shenghong Technology",
    "sh000001": "Shanghai Composite Index",
    "sz399001": "Shenzhen Component Index",
    "sz399006": "ChiNext Index",
    "hk01347": "Hua Hong Semiconductor (HK)",
}

def fetch_stock_data(stock_codes):
    """
    Fetch real-time stock data using Tencent Finance API
    Format: http://qt.gtimg.cn/q=sh601138,sz300985
    """
    url = "http://qt.gtimg.cn/q=" + ",".join(stock_codes)
    try:
        response = requests.get(url, timeout=5)
        response.raise_for_status()
        return response.text
    except requests.RequestException as e:
        print(f"Request failed: {e}")
        return None

def parse_stock_data(raw_data):
    """
    Parse the data returned by Tencent API
    """
    if not raw_data:
        return []
    
    stock_list = []
    for line in raw_data.split(";"):
        if not line.strip():
            continue
        
        data_str = line.split("=")[-1].strip('"')
        fields = data_str.split("~")
        
        if len(fields) < 5:
            continue
        
        stock_code = fields[2]
        # HK stocks are often returned as 5-digit codes; align with hk01347 in the mapping table
        lookup_key = (
            f"hk{stock_code}" if (len(stock_code) == 5 and stock_code.isdigit()) else stock_code
        )
        english_name = STOCK_NAME_MAPPING.get(lookup_key, STOCK_NAME_MAPPING.get(stock_code, fields[1]))
        # print(fields)
        stock_info = {
            "name": english_name,
            "code": stock_code,
            "price": fields[3],
            "open": fields[4],
            "high": fields[5],
            "low": fields[6],
            "volume": fields[7],
            "amount": fields[8],
            "buy1": fields[9],
            "sell1": fields[10],
            "lowest": fields[34],
            "highest": fields[33],
        }
        stock_list.append(stock_info)
    
    return stock_list

def display_stock_data(stock_list):
    """Display stock data with English names converted to pinyin"""
    # print("\n{:<30} {:<10} {:<10} {:<10} {:<15}".format(
    #     "Stock (Pinyin)", "Code", "Price", "Change", "Volume (lots)"))
    # print(stock_list)
    print("\n{:<20} {:<8} {:<6} {:<3} {:<8} {:<8} {:<15}".format(
        "name (Pinyin)", "current", "Change", "Rotio", "lowest", "highest", "Volume (lots)"))
    print("-" * 75)
    
    for stock in stock_list:
        change = float(stock['price']) - float(stock['open'])
        rotio = (change / float(stock['open'])) * 100 if float(stock['open']) != 0 else 0
        pinyin_name = " ".join(lazy_pinyin(stock['name']))  # Convert to pinyin
        # print(stock)
        print("{:<20} {:<8} {:<+6.2f} {:<+3.2f}% {:<8} {:<8} {:<15,}".format(
            pinyin_name,
            stock['price'],
            change,
            rotio,
            stock['lowest'],
            stock['highest'],
            int(stock['volume']) if stock['volume'] else 0
        ))

def main():
    stock_codes = list(STOCK_NAME_MAPPING.keys())
    
    while True:
        raw_data = fetch_stock_data(stock_codes)
        if raw_data:
            stock_list = parse_stock_data(raw_data)
            display_stock_data(stock_list)
        else:
            print("Failed to fetch data")
        
        # print("\n" + "=" * 75)
        # print(f"Last update: {time.strftime('%Y-%m-%d %H:%M:%S')}")
        # print("=" * 75 + "\n")
        
        time.sleep(2)  # Refresh every 10 seconds

if __name__ == "__main__":
    main()
