"use client";

import {useEffect, useState} from "react";
import getHistory from "@/app/history-actions";

export default function HistoryViewer() {
  const [history, setHistory] = useState<string | null>(null);

  useEffect(() => {
    getHistory().then(res => {
      setHistory(res);
    });
  })
  return (
    <div>
      <h1>History Viewer</h1>
      <p className={"whitespace-pre-wrap"}>{history}</p>
    </div>
  )
}