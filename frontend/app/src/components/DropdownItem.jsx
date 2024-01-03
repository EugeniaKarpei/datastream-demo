import React, { useContext } from "react"
import { FilterContext } from "./Filter"
import { AppContext } from "../App"
import FilterInputContext from "./FilterInputContext"

export default function DropdownItem({value}){
    // const { openDropdown } = useContext(FilterContext)
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