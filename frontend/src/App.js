import React from 'react';
import ItemsContainer from './Components/ItemsContainer.js';
import ItemInforamtionContainer from './Components/ItemInformationContainer.js';

import './Visualization/App.css';

const App =() => {

  const ItemsList = ItemsContainer();
  const ItemInformation = ItemInforamtionContainer();
  return (
    <div className="body-container">
      {ItemsList}
      {ItemInformation}
    </div>
  );
}

export default App;
