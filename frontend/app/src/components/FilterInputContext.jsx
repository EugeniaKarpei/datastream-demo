import React from 'react'

const FilterInputContext = React.createContext({
  inputValue: '',
  updateInputValue: () => {} // default value
})

export default FilterInputContext