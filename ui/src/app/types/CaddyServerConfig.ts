export type CaddyConfig = {
  apps: {
    http: HttpConfig;
  }
}

export type HttpConfig = {
  servers: {
    [key: string]: ServerConfig
  }
}

export type ServerConfig = {
  listen: string[];
  routes: ServerRoute[]
}

export type ServerRoute = {
  handle: Handler[]
  match: Match[]
}

export type Handler = {
  handler: string
  routes?: ServerRoute[]
}

export type Match = {
  path?: string[]
  method?: string[]
  host?: string[]
}