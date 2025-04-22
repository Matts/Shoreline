"use client";

import {useEffect, useState} from 'react';

const LogStream = () => {
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    const eventSource = new EventSource('/api/sse');

    // Handle incoming messages
    eventSource.onmessage = (event) => {
      const newLog = JSON.parse(event.data);
      setLogs((prevLogs) => [...prevLogs, `${newLog.timestamp}: ${newLog.message}`]);
    };

    // Handle errors
    eventSource.onerror = (error) => {
      console.error('Error with SSE connection:', error);
      eventSource.close(); // Close the connection on error
    };

    // Cleanup on unmount
    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div>
      <h3>Logs:</h3>
      <ul>
        {logs.map((log, idx) => (
          <li key={idx}>{log}</li>
        ))}
      </ul>
    </div>
  );
};

export default LogStream;