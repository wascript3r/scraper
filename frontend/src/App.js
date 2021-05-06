import React, { useState } from 'react';
import ItemsContainer from './Components/ItemsContainer.js';
import ScraperLogo from './Components/Logo.js';
import ItemInforamtionContainer from './Components/ItemInformationContainer.js';

// import './Visualization/App.css';

const App =() => {

  const [ID, setID] = useState(1);

  return (
    <div className="body-container container-fluid">
        <div className="row">
            <div className="col-sm-2 bg-secondary">
                <ScraperLogo />
                <ItemsContainer ID={ID} setID={value => setID(value)}/>
            </div>
            <div className="col-sm-10 bg-light">
                <ItemInforamtionContainer givenID={ID}/>
            </div>
        </div>
    </div>
  );
}

export default App;