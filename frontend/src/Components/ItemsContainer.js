import React, { useState, useEffect } from 'react';
import axios from 'axios';
import ScraperLogo from '../Components/Logo.js';
import '../Visualization/ItemsContainer.css';

export default function ItemsContainer({ID, setID}) {
    const [ebayItems, setEbayItems] = useState([]);
    useEffect(() => {
     
         /* This function gets the data from the API URL */
        const fetchData = async () =>{
             const result =  await axios('http://91.225.104.238:3000/api/queries/get');
             setEbayItems(result.data.data.queries);    
          };
          fetchData();
       }, [])

  
    return (
        <div className='items-container'>
            <ScraperLogo />
            <div className='objects'>
                {ebayItems.map(item => (
                    <li key={item.name} onClick={() => setID(item.id)} >
                         <div>{item.name}</div>
                         <hr />
                    </li>
        ))}
            </div>
        </div>
    )
}