import React, { useState, useEffect } from 'react';
import axios from 'axios';
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
        <div className='left-item-information'>
            <p>Item Name: {ebayItemName.name}</p>
            <hr />
            <p>Item Currency: {ebayItemName.currency}</p>
            <hr />
            <p>Item URL link:  
                <a href={ebayItemName.url} rel='noreferrer' target='_blank'>
                    <span>  Ebay.com/{ebayItemName.name}</span>
                </a>
            </p>
            <hr />
            <p>Tracking Item From: {historyItems[0]}</p>
            <hr />
        </div>
        <div className='right-item-information'>
            <p>Current Average Price: {ebayItemName.currentAvgPrice}</p>
            <hr />
            <p>Current Average Sold Price: {ebayItemName.currentAvgSoldPrice}</p>
            <hr />
            <p>Current Remaining Quantity: {ebayItemName.currentRemainingQuantity}</p>
            <hr />
            <p>Current Sold Quantity: {ebayItemName.currentSoldQuantity}</p>
            <hr />
        </div>
        </>
    )
}
