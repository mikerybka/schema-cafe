import React from 'react';
import ReactDOM from 'react-dom/client';

const data = JSON.parse(document.getElementById('data')!.textContent!);

function App() {
    return <div>{JSON.stringify(data)}</div>
}

const root = ReactDOM.createRoot(document.getElementById('root')!);
root.render(<App />);
