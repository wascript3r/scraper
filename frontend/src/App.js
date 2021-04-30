import React, { useState } from 'react';
import ItemsContainer from './Components/ItemsContainer.js';
import ItemInforamtionContainer from './Components/ItemInformationContainer.js';

import './Visualization/App.css';

const App =() => {

  const [ID, setID] = useState(1);

  return (
    <div className="body-container">
      <ItemsContainer ID={ID} setID={value => setID(value)}/>
      <ItemInforamtionContainer givenID={ID}/>
    </div>
  );
}

export default App;