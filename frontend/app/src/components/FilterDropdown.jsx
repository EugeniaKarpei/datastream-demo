import React, { useState, useEffect, useContext } from "react"
// import { FilterContext } from "./Filter"
import { AppContext } from "../App"
import DropdownItem from "./DropdownItem"
import useWebSocket, { ReadyState } from 'react-use-websocket'
import FilterInputContext from "./FilterInputContext"

export default function FilterDropdown(){
    const { open } = useContext(AppContext)
    const [filterData, setfilterData] = useState()
    const filtersSocketUrl = 'ws://localhost:8080/getFilters'
    const { inputValue } = useContext(FilterInputContext)

    const { sendJsonMessage, readyState, lastJsonMessage } = useWebSocket(filtersSocketUrl, {
        onOpen: () => console.log('opened'),
        // Will attempt to reconnect on all close events, such as server shutting down
        shouldReconnect: (closeEvent) => true,
        onClose: () => console.log('closed')
    })

    useEffect(() => {
        if (readyState === ReadyState.OPEN){
          sendJsonMessage({
            query: inputValue
          })
        }
    
      }, [open, inputValue, sendJsonMessage, readyState])

      useEffect(() => {
        if (lastJsonMessage){
          setfilterData((prev) => prev = lastJsonMessage)
        }
    
      }, [lastJsonMessage])

    function renderDropdownItems(){
        const dropdownItems = filterData.map((data, i) => <DropdownItem key={i} value={data} />)
        return dropdownItems
    }

    return (
        <>
        {(open && filterData) && 
        <div className="filter-dropdown">
            {renderDropdownItems()}
        </div>
        }
        </>
        
    )
}