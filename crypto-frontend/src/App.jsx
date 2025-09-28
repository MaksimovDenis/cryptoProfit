import React, { useEffect, useState, useRef } from "react";
import "./App.css";

export default function App() {
  const [tab, setTab] = useState("profit"); // "profit" или "price"
  const [prices, setPrices] = useState([]);
  const [profit, setProfit] = useState({ profits: [], revenue: 0 });
  const [query, setQuery] = useState("");

  const wsRef = useRef(null);

  useEffect(() => {
    // Выбираем URL в зависимости от таба
    const url = tab === "profit" ? "ws://localhost:8081/ws/profit" : "ws://localhost:8081/ws/price";

    // Создаем WS соединение только один раз для текущей вкладки
    wsRef.current = new WebSocket(url);

    wsRef.current.onopen = () => console.log(`WebSocket connected to ${url}`);
    wsRef.current.onclose = () => console.log("WebSocket closed");
    wsRef.current.onerror = (e) => console.error("WebSocket error:", e);

    wsRef.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (tab === "profit") setProfit(data);
        else setPrices(data);
      } catch (err) {
        console.error("Failed to parse WS message:", err);
      }
    };

    // Закрываем WS при размонтировании компонента или смене таба
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [tab]); // пересоздается только при смене таба

  const renderPrices = () => {
    const filtered = prices.filter(r =>
      r.symbol.toLowerCase().includes(query.toLowerCase())
    );

    return (
      <table className="prices-table">
        <thead>
          <tr>
            <th>Тикер</th>
            <th>Цена</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map(r => (
            <tr key={r.symbol}>
              <td>{r.symbol}</td>
              <td>${r.price.toFixed(4)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  };

  const renderProfit = () => {
    const filtered = profit.profits.filter(r =>
      r.ticker.toLowerCase().includes(query.toLowerCase())
    );

    return (
      <>
        <div className="summary">
          Общий доход: <strong>{profit.revenue.toFixed(2)} USDT</strong>
        </div>
        <table>
          <thead>
            <tr>
              <th>Тикер</th>
              <th>Баланс $</th>
              <th>Профит</th>
              <th>Профит %</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map(r => (
              <tr key={r.ticker}>
                <td>{r.ticker}</td>
                <td>{r.balance.toFixed(2)}</td>
                <td className={r.profit > 0 ? "pos" : r.profit < 0 ? "neg" : ""}>
                  {r.profit.toFixed(2)}
                </td>
                <td className={r.profit_percent > 0 ? "pos" : r.profit_percent < 0 ? "neg" : ""}>
                  {r.profit_percent.toFixed(2)}%
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </>
    );
  };

  return (
    <div className="container">
      <h1>Crypto Dashboard</h1>

      <div className="tabs">
        <div
          onClick={() => setTab("profit")}
          className={"tab " + (tab === "profit" ? "active" : "")}
        >
          Профит
        </div>
        <div
          onClick={() => setTab("price")}
          className={"tab " + (tab === "price" ? "active" : "")}
        >
          Цены
        </div>
      </div>

      <div className="controls">
        <input
          className="search"
          placeholder="Фильтр по тикеру…"
          value={query}
          onChange={e => setQuery(e.target.value)}
        />
      </div>

      {tab === "profit" ? renderProfit() : renderPrices()}
    </div>
  );
}