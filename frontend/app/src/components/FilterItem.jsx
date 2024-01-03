import React from "react"
import { IoIosCloseCircleOutline } from "react-icons/io"
import { AppContext } from "../App"

export default function FilterItem({value}){
    const { removeFilter } = React.useContext(AppContext)
    
    return (
        <div className="filter-item">
            <p>{value}</p>
            <IoIosCloseCircleOutline 
               className="remove-filter-btn" 
               onClick={() => removeFilter(value)}
            />
        </div>
    )
}