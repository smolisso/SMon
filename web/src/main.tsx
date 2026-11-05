import { StrictMode, useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';
import './style.css';

type Connection = 'loading' | 'ok' | 'unavailable';
type Database = 'loading' | 'ready' | 'not_configured' | 'unavailable_or_unmigrated' | 'unavailable';

function App() {
  const [api, setAPI] = useState<Connection>('loading');
  const [database, setDatabase] = useState<Database>('loading');
  const [updated, setUpdated] = useState<string>('');
  useEffect(() => {
    const controller = new AbortController();
    let busy = false;
    async function refresh() {
      if (busy) return;
      busy = true;
      const signal = AbortSignal.any([controller.signal, AbortSignal.timeout(5000)]);
      try {
        const [health, ready] = await Promise.allSettled([
          fetch('/api/v1/health', { signal }).then(async r => {
            const data = await r.json();
            if (!r.ok || data.status !== 'ok') throw new Error('health');
          }),
          fetch('/api/v1/ready', { signal }).then(async r => {
            const data = await r.json();
            const states: Database[] = ['ready', 'not_configured', 'unavailable_or_unmigrated'];
            if (![200, 503].includes(r.status) || !states.includes(data.database)) throw new Error('ready');
            return data.database as Database;
          }),
        ]);
        if (controller.signal.aborted) return;
        setAPI(health.status === 'fulfilled' ? 'ok' : 'unavailable');
        setDatabase(ready.status === 'fulfilled' ? ready.value : 'unavailable');
        setUpdated(new Date().toLocaleTimeString('it-IT'));
      } finally { busy = false; }
    }
    void refresh();
    const timer = window.setInterval(() => { void refresh(); }, 15000);
    return () => { controller.abort(); window.clearInterval(timer); };
  }, []);

  const databaseLabels: Record<Database, string> = {
    loading: 'Verifica in corso', ready: 'Pronto', not_configured: 'Da configurare',
    unavailable_or_unmigrated: 'Connessione o schema da verificare', unavailable: 'Non raggiungibile',
  };
  return <div className="layout">
    <aside>
      <a className="brand" href="/" aria-label="SMon, panoramica"><span className="brand-icon">S</span>SMon<span className="brand-dot">●</span></a>
      <p className="tagline">La tua esperienza. Il tuo monitoring.</p>
      <p className="nav-label">WORKSPACE</p>
      <nav aria-label="Navigazione principale"><a href="/" aria-current="page"><span aria-hidden="true">◉</span> Panoramica</a></nav>
      <div className="aside-footer"><span className="small-dot" /> Progetto open source<br /><span className="muted">di Stefano Molisso</span></div>
    </aside>
    <main>
      <header><div><p className="eyebrow">SMON / PANORAMICA</p><h1>Il punto di partenza.</h1><p className="intro">Un nuovo sistema di monitoraggio, costruito sulla tua esperienza.</p></div><span className="stage">Prima base</span></header>
      <section className="status-grid" aria-label="Stato del sistema" aria-live="polite">
        <article className="status-card"><p className="card-label">SERVER SMON</p><h2 className={api === 'ok' ? 'healthy' : ''}>{api === 'loading' ? 'Verifica in corso' : api === 'ok' ? 'Online' : 'Non raggiungibile'}</h2><p>API REST · verifica ogni 15 secondi</p></article>
        <article className="status-card"><p className="card-label">POSTGRESQL</p><h2 className={database === 'ready' ? 'healthy' : ''}>{databaseLabels[database]}</h2><p>Connessione e migrazione dello schema</p></article>
        <article className="status-card"><p className="card-label">MONITORAGGIO</p><h2>Non ancora attivo</h2><p>Nessun controllo eseguito in questa fase</p></article>
      </section>
      <section className="welcome"><div className="welcome-mark" aria-hidden="true">S</div><p className="eyebrow">DAL 2007, UNO SCRIPT ALLA VOLTA</p><h2>Il monitoring porta la tua firma.</h2><p>SMon nasce per dare una nuova casa ai tuoi controlli e un’interfaccia al modo in cui lavori.</p><div className="next"><span className="step">01</span><div><h3>La prossima tappa</h3><p>Collegare un host, eseguire un plugin Nagios e leggere il primo risultato.</p></div></div></section>
      <footer><span>SMon · Fondamenta del progetto</span><span>{updated ? `Ultima verifica ${updated}` : 'Verifica in corso…'}</span></footer>
    </main>
  </div>;
}

createRoot(document.getElementById('root')!).render(<StrictMode><App /></StrictMode>);
