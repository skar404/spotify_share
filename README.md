# Spotify share 
Telegram bot in share link play in you spotify

## Spotify link

- User apps: https://www.spotify.com/account/apps/
- App dashboard: https://developer.spotify.com/dashboard/applications/
- OAuth: https://developer.spotify.com/documentation/general/guides/authorization-guide/
- Get the User's Currently Playing Track: https://developer.spotify.com/documentation/web-api/reference/player/get-the-users-currently-playing-track/

## Schema architecture app 

<!-- generated by mermaid compile action - START -->
![~mermaid diagram 1~](/.resources/README-md-1.png)
<details>
  <summary>Mermaid markup</summary>

```mermaid
flowchart TB
    client[client]
    tg[Telegram] 

    spotify[Spotify]
    
    subgraph pn["Private network"]
    bot[BotApp]
    db[(MongoDB)]
    end
    tr[Treafic]
    
    client --> tg
    tg -- "TLS Web hook" --> tr
    tr --> bot
    db <--> bot

    spotify -- "TLS OAuth" --> tr
    bot -- "HTTP API" --> spotify
    
    bot -- "HTTP API (Bot Token)" --> tg 
```

</details>
<!-- generated by mermaid compile action - END -->