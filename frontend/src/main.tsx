import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Toaster } from 'sonner';
import { SWRConfig } from 'swr';
import App from './App.tsx';
import './index.css';

function fetcher(url: string) {
  return fetch(`${import.meta.env.VITE_API_URL}${url}`).then(res => res.json());
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <div className="bg-neutral-950 text-white">
      <SWRConfig value={{ fetcher }}>
        <App />
        <Toaster />
      </SWRConfig>
    </div>
  </StrictMode>,
);
