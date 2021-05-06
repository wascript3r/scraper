import React, { useState, useEffect } from 'react';
import { Line } from 'react-chartjs-2';

import '../../Visualization/LeftQuantity.css'
export default function LeftQuantity(props){
    const [historyData, setData] = useState([])

    useEffect(() => {
     
       const takeData = () =>{
            setData(props.data.history);    
         };
         takeData();
      }, [props.data.history])


const data = {
  labels: historyData?.slice(0).reverse().map(item => (
      item.date
  )),
  datasets: [
    {
      label: 'Prekių likutis',
      data:  historyData?.slice(0).reverse().map(item => (
        item.remainingQuantity
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
        <div className='leftQuantity-chart'>
            <Line data={data} options={options} />
        </div>
    )
}