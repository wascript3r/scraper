import React from 'react';
import ItemInformation from '../Components/ItemInformation.js';
import '../Visualization/ItemInformationContainer.css';

export default function ItemInforamtionContainer({givenID}) {
    return (
        <div className='itemInformation-container'>
            <ItemInformation id={givenID}/>
            
        </div>
    )
}