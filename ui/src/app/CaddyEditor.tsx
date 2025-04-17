"use client";
import {useEffect, useRef, useState, useTransition} from "react";
import {getServers, updateServerConfig} from "@/app/caddy-actions";
import {HttpConfig} from "@/app/types/CaddyServerConfig";

export default function CaddyEditor() {
  const [data, setData] = useState<HttpConfig | null>(null);
  const [isPending, startTransition] = useTransition();

  const [selectedServer, setSelectedServer] = useState<string | null>(null);
  const editorRef = useRef<HTMLTextAreaElement|null>(null);

  useEffect(() => {
    startTransition(async () => {
      const res = await getServers();
      setData(res);
    });
  }, []);

  const updateServer = () => {
    if(!editorRef.current) return;

    const updatedConfig = editorRef.current.value;
    const serverName = selectedServer;

    updateServerConfig(serverName!, JSON.parse(updatedConfig));
  }

  return (
    <div>
      {data && (
        <div className={"flex h-full w-full"}>
          <ul className={"w-24"}>
            {Object.entries(data.servers).map(([serverName, serverConfig]) => (
              <li key={serverName}>
                <button onClick={() => setSelectedServer(serverName)}>
                  {serverName}
                </button>
              </li>
            ))}
          </ul>
          {selectedServer && (
            <div className={"w-full h-full"}>
              <h3>{selectedServer} Configuration</h3>
              <textarea key={selectedServer} className={"w-full h-full"} ref={editorRef} defaultValue={JSON.stringify(data.servers[selectedServer], null, 2)}></textarea>
              <button onClick={updateServer}>Update Server</button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}