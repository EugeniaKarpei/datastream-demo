import React, { useContext } from "react"
import { AppContext } from "../App"
import FilterInputContext from "./FilterInputContext"

export default function DropdownItem({value}){
    const { updateFilters, openDropdown } = useContext(AppContext)
    const { updateInputValue } = useContext(FilterInputContext)

    function handleItemClick(){
        updateFilters(value)
        updateInputValue(value)
        openDropdown(false)
    }

    return (
        <div className="dropdown-item" onClick={() => handleItemClick()}>
            {value}
        </div>
    )
}