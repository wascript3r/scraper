import React, { useState, useEffect } from 'react';
import axios from 'axios';
// import ScraperLogo from '../Components/Logo.js';
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

    // our brand name
    // <ScraperLogo />
  
    return (
        <div className='items-container sidebar-nav container'>
            {/* <span class="navbar-brand mb-0 h1"></span> */}
            <nav className="navbar navbar-expand-lg navbar-dark bg-secondary">
                <button type="button" className="navbar-toggler" data-toggle="collapse" data-target="#navbarSupportedContent">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div className="collapse navbar-collapse" id="navbarSupportedContent">
                    <ul className="navbar-nav mr-auto">
                        {ebayItems.map(item => (
                            <li key={item.name} onClick={() => setID(item.id)} className="nav-item btn btn-outline-light button">
                                <a className="nav-link text-white">{item.name}</a>
                                {/* <hr /> */}
                            </li>
                        ))}
                    </ul>
                </div>
            </nav>
        </div>
    )
}