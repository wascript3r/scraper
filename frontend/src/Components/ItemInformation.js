import React, { useState, useEffect } from 'react';
import axios from 'axios';
import PriceDifference from './Charts/PriceDifference';
import LeftQuantity from './Charts/LeftQuantity';
import AverageSellPrice from './Charts/AverageSellPrice';
import SoldAmount from './Charts/SoldAmount';

import '../Visualization/ItemInformation.css';

export default function ItemInformation(props){

    const [ebayItemName, setEbayItemName] = useState([]);
    const [historyItems, setHistoryItems] = useState([]);
    useEffect(() => {
        const fetchItemData = async () =>{
            const results =  await axios(
                {
                    method: 'post',
                    url: 'http://91.225.104.238:3000/api/query/stats',
                    data: {
                        id: props.id
                    }
                }
            )
            setEbayItemName(results.data.data); 
            setHistoryItems(results.data.data.history.map(({ date }) => date))  
            console.log(results.data.data)
         };
         fetchItemData();
    }, [props.id]);

    return(
        <>
        <div className='table-information table-responsive'>
            <table className="table table-striped table-hover table-bordered">
                <tbody>
                    <tr>
                        <th>Item Name:</th>
                        <td>
                            <div>{ebayItemName.name}</div>
                        </td>
                        <th>Current Average Price:</th>
                        <td>
                            <div>{ebayItemName.currentAvgPrice}</div>
                        </td>
                    </tr>
                    <tr>
                        <th>Item Currency:</th>
                        <td>
                            <div>{ebayItemName.currency}</div>
                        </td>
                        <th>Current Average Sold Price:</th>
                        <td>
                            <div>{ebayItemName.currentAvgSoldPrice}</div>
                        </td>
                    </tr>
                    <tr>
                        <th>Item URL link:</th>
                        <td>
                            <a href={ebayItemName.url} rel='noreferrer' target='_blank'>
                                <span>  Ebay.com/{ebayItemName.name}</span>
                            </a>
                        </td>
                        <th>Current Remaining Quantity: </th>
                        <td>
                            <div>{ebayItemName.currentRemainingQuantity}</div>
                        </td>
                    </tr>
                    <tr>
                        <th>Tracking Item From:</th>
                        <td>
                            <div>{historyItems[historyItems.length - 1]}</div>
                        </td>
                        <th>Current Sold Quantity:</th>
                        <td>
                            <div>{ebayItemName.currentSoldQuantity}</div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div className="row">
            <div className="col-md-12">
                <PriceDifference data={ebayItemName} />
            </div>
        </div>
        <div className="row">
            <div className="col-md-12">
            <LeftQuantity  data={ebayItemName} />
            </div>
        </div>
        <div className="row">
            <div className="col-md-12">
            <AverageSellPrice  data={ebayItemName} />
            </div>
        </div>
        <div className="row">
            <div className="col-md-12">
            <SoldAmount data={ebayItemName} />
            </div>
        </div>
        </>
    )
}
