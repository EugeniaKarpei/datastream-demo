import React, { useContext } from "react"
import { AppContext } from "../App"
import FilterInputContext from "./FilterInputContext"

export default function FilterInput(){
    const { openDropdown } = useContext(AppContext)
    const { inputValue, updateInputValue } = useContext(FilterInputContext)

    return (
        <input className="filter-input" 
               name="filter"
               value={inputValue}
               placeholder="filter by" 
               onFocus={() => openDropdown(true)}
               onChange={e => updateInputValue(e.target.value)}
        />
    )
}