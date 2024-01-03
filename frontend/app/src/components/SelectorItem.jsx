import React from "react"

export default function Item({itemName, isSelected, onClick}){
    return (
        <div className={`scale-item ${isSelected && 'selected'}`} 
             onClick={() => onClick(itemName)}
        >
            {itemName}
        </div>
    )
}