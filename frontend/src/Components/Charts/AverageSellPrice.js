import React, { useState, useEffect } from 'react';
import { Line } from 'react-chartjs-2';

import '../../Visualization/AverageSellPrice.css';

export default function AverageSellPrice(props){
    const [historyData, setData] = useState([])

    useEffect(() => {
     
       const takeData = () =>{
            setData(props.data.soldHistory);    
         };
         takeData();
      }, [props.data.soldHistory])

const data = {
  labels: historyData?.slice(0).reverse().map(item => (
      item.date
  )),
  datasets: [
    {
      label: 'Vidutinė pardavimo kaina',
      data:  historyData?.slice(0).reverse().map(item => (
        item.avgPrice
    )),
      fill: false,
      backgroundColor: 'rgb(255, 99, 132)',
      borderColor: 'rgba(255, 99, 132, 0.2)',
    },
  ],
};
const options = {
    maintainAspectRatio: false, 
    scales: {
      y:{
          ticks: {
            color: "black",
            font:{
                weight: "bold"
            }
          }
        },
      x:{
          ticks:{
              color: "black",
              font:{
                  weight: "bold",
              },
          }
      },
    },
  };

  return(
      <div className='averageSellPrice-chart'>
          <Line data={data} options={options} />
      </div>
  )
}