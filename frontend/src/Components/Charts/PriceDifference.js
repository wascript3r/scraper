import React, { useState, useEffect } from 'react';
import { Line } from 'react-chartjs-2';

import '../../Visualization/PriceDifference.css';

export default function PriceDifference(props){

    const [historyData, setData] = useState([])

    useEffect(() => {
     
       const takeData = () =>{
            setData(props.data.history);    
         };
         takeData();
      }, [props.data.history])


const data = {
  labels: historyData?.map(item => (
      item.date
  )),
  datasets: [
    {
      label: 'Kainų pokytis',
      data:  historyData?.map(item => (
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
        <div className='price-chart'>
            <Line data={data} options={options} />
        </div>
    )
}