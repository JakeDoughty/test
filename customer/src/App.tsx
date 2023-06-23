import { useEffect, useState } from 'react';
import './App.css';
import { Link, useLocation } from 'react-router-dom';

function App() {
  const location = useLocation();
  const [count, setCount] = useState(0);
  const [savedLocation, setSavedLocation] = useState('');

  useEffect(() => {
    if (savedLocation != location.pathname) {
      console.log('location changed ', location);
      dispatchEvent(new Event('popstate'));
      setSavedLocation(location.pathname);
    }
  }, [location]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <>
      <div>
        <Link to="/">Home</Link> | <Link to="/private">Private</Link>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;
