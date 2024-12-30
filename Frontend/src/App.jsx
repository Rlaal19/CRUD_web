import { useState } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Dashboard from './assets/dashboaed';
import Edit from './assets/edit/edit';


function App() {

  return (
    <div>
      <BrowserRouter>
          <div className="navbar flex justify-between bg-blue-700   text-white py-1 text-3xl ">Employee</div>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="edit/:id" element={<Edit />} />
          </Routes>
        </BrowserRouter>
    </div>
  )
}

export default App
