# Gen IV Distribution Server
A simple server for distributing WonderCards to GenIV Pokemon games

## Setup
The server is useless without redirecting requests to it

### Requirements:
- Nginx compiled with support for SSLv3, RC4-MD5/RC4-SHA, and 1024-bit RSA keys
- A DNS server, such as Bind9
- [nds-constrain't](https://github.com/KaeruTeam/nds-constraint) compadible certificate
- [WC4/PCD WonderCard](https://github.com/projectpokemon/EventsGallery) files to serve

### Steps:
1) Generate nds-constrain't compadible certificates
2) Configure your DNS server to point `conntest.nintendowifi.net`, `nas.nintendowifi.net`, and `dls1.nintendowifi.net` to your Nginx server
3) Configure Nginx to point:
    - `conntest.nintendowifi.net/` on port 80 -> `http://127.0.0.12:8080/runConnTest`
    - `nas.nintendowifi.net/` on port 80 -> `http://127.0.0.12:8080/nas/`
    - `nas.nintendowifi.net/` on port 443 with SSL enabled -> `http://127.0.0.12:8080/nas/`
    - `dls1.nintendowifi.net/` on port 80 -> `http://127.0.0.12:8080/dls1/`
    - `dls1.nintendowifi.net/` on port 443 with SSL enabled -> `http://127.0.0.12:8080/dls1/`
    - Use your cert/key you generated in Step 1 as the SSL cert/key
4) Start event server
5) Copy WC4/PCD files to cards/geniv


Special thanks to teams/individuals at ProjectPokemon, Custom MarioKart Wiiki, Wiimmfi, and WFCReplay. Without their hard work, this would not have been possible
