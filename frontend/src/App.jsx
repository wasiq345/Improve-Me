import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import React from 'react'
import Dashboard from './dashboard'
import Login from './login'
import Register from './register'
import Profile from './profile'
import ProtectedRoute from './protected_route'
import CreateNote from './create_note'
import GetNotes from './get_notes'
import ReadNote from './read_note'

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path='/' element={<Dashboard />} />
        <Route path='/LoginUser' element={<Login />} />
        <Route path='/RegisterUser' element={<Register />} />
        <Route
          path='/Profile/:username'
          element={
            <ProtectedRoute>
              <Profile />
            </ProtectedRoute>
          } />
        <Route
          path='/Profile/:username/Createnote'
          element={
            <ProtectedRoute>
              <CreateNote />
            </ProtectedRoute>
          } />
        <Route
          path='/Profile/:username/GetNotes'
          element={
            <ProtectedRoute>
              <GetNotes />
            </ProtectedRoute>
          } />
        <Route
          path='/Profile/:username/ReadNote/:noteId'
          element={
            <ProtectedRoute>
              <ReadNote />
            </ProtectedRoute>
          } />
      </Routes>
    </Router>
  )
}

