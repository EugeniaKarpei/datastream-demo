import React from "react"
import { ResponsiveContainer, LineChart, Line, CartesianGrid, XAxis, YAxis, Tooltip, Legend } from 'recharts'


export default function Chart({data}){

    return <div className="chart-container">
    <ResponsiveContainer width="100%" height={350}>
      <LineChart data={data} margin={{ top: 30, right: 35, bottom: 10, left: 0 }}>
        <Line type="monotone" dataKey="value" stroke="#8884d8" />
        <CartesianGrid stroke="#ccc" strokeDasharray="3 3"/>
        <XAxis dataKey={"scale"} />
        <YAxis  tickCount={6} />
        <Tooltip />
        <Legend />
      </LineChart>
    </ResponsiveContainer>
  </div>
}