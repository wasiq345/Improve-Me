import { BrowserRouter as Router, Routes, Route} from 'react-router-dom'
import React from 'react'
import Dashboard from './dashboard'
import Login from './login'
import Register from './register'

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path='/' element={<Dashboard />} />
        <Route path='/LoginUser' element={<Login />} />
        <Route path='/RegisterUser' element={<Register />} />
      </Routes>
    </Router>
  )
}

