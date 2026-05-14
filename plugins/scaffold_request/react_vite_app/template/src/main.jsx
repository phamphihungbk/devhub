import React from 'react'
import { createRoot } from 'react-dom/client'
import './style.css'

function App() {
  return (
    <main>
      <h1>[[SERVICE_NAME]]</h1>
      <p>React frontend running in [[ENVIRONMENT]].</p>
    </main>
  )
}

createRoot(document.getElementById('root')).render(<App />)
