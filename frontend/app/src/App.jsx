import React, { useState, useEffect } from "react"
import Header from "./components/Header.jsx"
import FilterInput from "./components/FilterInput.jsx"
import FilterDropdown from "./components/FilterDropdown.jsx"
import Filter from "./components/Filter.jsx" 
import AggregateBySelector from "./components/AggregateBySelector.jsx"
import SelectorItem from "./components/SelectorItem.jsx"
import useWebSocket, { ReadyState } from 'react-use-websocket'
import Chart from "./components/Chart.jsx"

const AppContext = React.createContext()
export { AppContext }

function App() {
  const DEFAULT_SCALE = "Monthly"
  const DEFAULT_AGGREGATE_BY = "Sum"
  
  const [data, setData] = useState([])
  const [open, setOpen] = useState(false)
  const [location, setLocation] = useState("")

  const scaleItems = ["Daily", "Weekly", "Monthly"]
  const [currentScale, setCurrentScale] = useState(DEFAULT_SCALE)

  const aggregateByItems = ["Sum", "Avg", "Count"]
  const [currentAggregateByItem, setCurrentAggregateByItem] = useState(DEFAULT_AGGREGATE_BY)

  const [filters, setFilters] = useState([])

  const dataSocketUrl = 'ws://localhost:8080/getData'
  
  
  const { sendJsonMessage, readyState, lastJsonMessage } = useWebSocket(dataSocketUrl, {
    onOpen: () => console.log('opened'),
    //Will attempt to reconnect on all close events, such as server shutting down
    shouldReconnect: (closeEvent) => true,
    onClose: () => console.log('closed')
  })

  useEffect(() => {
    if (readyState === ReadyState.OPEN){
      sendJsonMessage({
        filters: [...filters],
        scale: currentScale,
        aggregator: currentAggregateByItem
      });
    }

  }, [readyState, currentScale, currentAggregateByItem, filters])

  useEffect(() => {
    if (lastJsonMessage){
      const newData = lastJsonMessage.map(item => {
        return {
          scale: getScale(item.timestamp),
          value: item.value
        }
      })
      
      setData((prev) => prev = newData)
    }

  }, [lastJsonMessage])

  function getScale(timestamp){
    const date = new Date(timestamp)
    if (currentScale === "Monthly"){
      return date.getMonth() + 1
    }else {
      return `${date.getDate()}/${date.getMonth() + 1}`
    }
  }

  function renderScaleByItems(){
    return scaleItems.map((item) => {
      return <SelectorItem key={item}
                           itemName={item} 
                           isSelected={item === currentScale}
                           onClick={(value) => setCurrentScale(value)}
              />
    })
  }

  function renderAggregateByItems(){
    return aggregateByItems.map((item) => {
      return <SelectorItem key={item}
                           itemName={item}
                           isSelected={item === currentAggregateByItem}
                           onClick={(value) => setCurrentAggregateByItem(value)}
              />
    })
  }

  function openDropdown(isOpen){
        setOpen(isOpen)
  }

  function updateFilters(value){
    if (!filters.includes(value)){
      setFilters(prev => [...prev, value])
      if (value.includes('location')){
        setLocation(value.split(':')[1])
      }
    }
  }

  function removeFilter(value){
    const newFilters = filters.filter(filter => filter !== value)
    setFilters(newFilters)
    if (value.includes(location)){
      setLocation("")
    }
  }

  function handleMainAreaClick(e){
    const targetTagName = e.target.tagName 
    if(open && targetTagName !== 'INPUT'){
      openDropdown(false)
    }
  }

  return (
    <>
      <Header />
      <main onClick={e => handleMainAreaClick(e)}>
        <div className="chart-description">
          <p>Sales Data {location !== "" && `for ${location}`}</p>
        </div>
        <Chart data={data} />
        <AggregateBySelector>
          {renderAggregateByItems()}
        </AggregateBySelector>
        <AggregateBySelector>
          {renderScaleByItems()}
        </AggregateBySelector>
        <AppContext.Provider value={{open, openDropdown, updateFilters, removeFilter}}>
          <Filter filters={filters}>
            <FilterInput />
            <FilterDropdown />
          </Filter>
        </AppContext.Provider>
      </main>
    </>
  )
}

export default App;