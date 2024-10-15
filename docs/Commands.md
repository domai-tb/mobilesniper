# Commands

- **`enum services`:** Enumerates ports on a given network or IP to identify potential 5G networking functions. This command scans open ports to detect services that could be linked to specific 5G functions.
  
- **`enum nf`:** Enumerates the given network or IP on HTTP(S) ports and checks if 5G network function endpoints match the responses on those ports. The results are presented as a percentage of matching HTTP endpoints, helping identify 5G services based on known endpoint patterns.

- **`analyze pcap`:** Analyzes a given PCAP file to detect HTTP endpoints that have been called during network activity. If any of the called endpoints match predefined network function definitions, the tool flags them as potential 5G networking functions.

- **`scan nuclei`:** Runs a Nuclei scan against the given network or IP, utilizing templates to check for known vulnerabilities or misconfigurations in the 5G network.

- **`scan nessus`:** Executes a Nessus scan on the specified network or IP, performing a comprehensive vulnerability assessment to uncover potential security weaknesses in the 5G infrastructure.