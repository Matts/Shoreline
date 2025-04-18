'use server';

import {HttpConfig, ServerConfig} from "@/app/types/CaddyServerConfig";

export async function getServers() {
  const response: HttpConfig = await fetch('http://localhost:2019/config/apps/http').then((configs) => {
    return configs.json()
  }).catch((er) => {
    console.log(er);

    return null;
  });
  if (response == null) {
    return null;
  }

  return response;
}

export async function updateServerConfig(serverName: string, newConfig: ServerConfig) {
  console.log(newConfig)
  const response = await fetch(`http://localhost:2019/config/apps/http/servers/${serverName}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(newConfig),
  }).catch((er) => {
    console.log(er);

    return er;
  });

  if (!response.ok) {
    throw new Error('Failed to update server configuration');
  }

  return "updated"
}