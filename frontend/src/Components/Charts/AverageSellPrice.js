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
      backgroundColor: 'rgb(173, 83, 0)',
      borderColor: 'rgb(221, 137, 35)',
      pointBorderWidth: 5,
      pointBorderColor: 'rgb(173, 83, 0)',
      pointRadius: 3,
    },
  ],
};
const options = {
    maintainAspectRatio: false,
    plugins: {
      legend: {
        labels:{
          color: 'black',
          font:{
            size: 15,
            weight: "600",
          }
        }
      }
    },  
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