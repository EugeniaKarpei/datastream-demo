import React, { useState, createContext, useEffect } from "react"
import FilterItemsList from "./FilterItemsList"
import FilterItem from "./FilterItem"
import FilterInputContext from "./FilterInputContext"

const FilterContext = createContext()
export { FilterContext }

export default function Filter({children, filters}){
    // const [open, setOpen] = useState(false)
    const [inputValue, setInputValue] = useState("")

    useEffect(() => {
        if (filters.length === 0 ){
            setInputValue("")
        }
    }, [filters])

    function updateInputValue(newValue){
        setInputValue(newValue)
    }

    function renderFilters(){
        return filters.map((filter, i) => <FilterItem key={`filter-${i}`} value={filter} />)
    }

    return (
        <div className="filter-container">
            {filters.length > 0 && 
            <FilterItemsList>
                {renderFilters()}
            </FilterItemsList>}
            <FilterInputContext.Provider value={{inputValue, updateInputValue}}>
                {children}
            </FilterInputContext.Provider>
        </div>
    )
}