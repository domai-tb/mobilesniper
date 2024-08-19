# mobilesniper

A 5G mobile network penetration testing tool that 

## Current Status 

**1. Implemented a "enum services" command**

The command performs a multi-threaded port enumeration with Nmap and report the results to the commandline. Essentially it is just a wrapper around Nmap with optimized commandline operators. In our Cocus Campus Lab, the initially Nmap TCP scan over all ports and full network ranges (e.g. NGC Control & Data plane, UE pool etc.) took several hours. The optimized wrapper tooks around ~20 minute in a single network range (CIDR /24). 

**2. Implemented a "enum nf" command**

After found possible network function (and other services) via the "enum services" command (or normal nmap), the "enum nf" command should detect which HTTP(S) service in the network is a NF. For that, I tried to match the given OpenAPI definitions (by 3GPP) to the reachable paths on the detected web service. In our Cocus Campus Lab that turned out as a not really suitable / accurate method. 

**3. Currently implementing a "analyze pcap" command**

To increase the accurancy of the NF detection, we want to use a offline / static analyzation of a pre-captured PCAP-file. That should lead to a way more accurate NF discovery, but is limited by the NFs contacted within in the PCAP file. In a engagement, we could capture the initially 5G connection and then analyze them. 

## Roadmap

**1. Implement "analyze pcap" Command**
   - **Objective**: Extend the "analyze pcap" and improve detection accuracy.
   - **Steps**:
     - Implement protocol-specific analysis features, such as detection of specific NF interactions or security misconfigurations.
     - Incorporate these parsers into the "analyze pcap" command for a more comprehensive analysis.

**2. Develop a "vuln scan" Command for NF Vulnerability Detection**
   - **Objective**: Implement a command that identifies vulnerabilities in detected NFs by scanning for known CVEs, misconfigurations, and insecure deployments. Propaly by using a well know vulnerability scanner like Nessus or Nuclei.
   - **Steps**:
     - Integrate with public vulnerability databases (e.g., NVD) for up-to-date CVE information.
     - Implement scanning techniques to detect known vulnerabilities in NFs.
     - Provide detailed vulnerability reports with remediation suggestions.

**3. Add "traffic simulate" Command for Active NF Probing**
   - **Objective**: Enable active testing of NFs by simulating 5G network traffic, allowing for deeper analysis of NF behavior and potential security weaknesses.
   - **Steps**:
     - Develop traffic generation modules that simulate common 5G interactions (e.g., authentication, session management).
     - Allow customization of traffic parameters to test specific scenarios.
     - Analyze NF responses to simulated traffic to identify potential security issues.

## Requirements

These tools must be installed before succesfully running all MobileSniper commands:

- [nmap](https://nmap.org)