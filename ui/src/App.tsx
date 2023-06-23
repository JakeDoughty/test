import { Route, Routes } from 'react-router-dom';
import './App.css';
import NavBar from './components/NavBar';
import ApplicationsPage from './components/ApplicationsPage';
import EventsPage from './components/EventsPage';
import { Application, createBackendConnection } from './components/backend';
import { useState } from 'react';

function App() {
  const backend = createBackendConnection('http://localhost:3000');

  const [applications, setApplications] = useState<Application[]>([]);

  return (
    <div className="container">
      <NavBar />
      <Routes>
        <Route
          path="/"
          element={<EventsPage appId={} backend={backend} />}
        ></Route>
        <Route path="/applications" element={<ApplicationsPage />}></Route>
      </Routes>
    </div>
  );
}

export default App;
