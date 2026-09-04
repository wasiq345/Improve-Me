import { BrowserRouter as Router, Routes, Route} from 'react-router-dom'
import React from 'react'
import Dashboard from './dashboard'
import Login from './login'
import Register from './register'
import Profile from './profile'
import ProtectedRoute from './protected_route'

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path='/' element={<Dashboard />} />
        <Route path='/LoginUser' element={<Login />} />
        <Route path='/RegisterUser' element={<Register />} />
        <Route 
          path='/Profile' 
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            } />
      </Routes>
    </Router>
  )
}

