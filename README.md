## Tailscale & the Domain Controller Locator Process

### Problem Statement
Large Windows environments often have multiple Domain Controllers in different
regions of the planet. Client Windows systems use the DC Locator Process to find
the nearest Domain Controller.

The subnet-to-site mapping in Active Directory would, in a traditional network,
consist mostly of CIDR ranges:

| Name | Site | Location | Type | Description
|------|------|----------|------|------------
|192.168.1.0/16|LHR-DC|London|Subnet|
|192.168.2.0/16|ATL-DC|Atlanta|Subnet|
|192.168.3.0/16|TOK-DC|Tokyo|Subnet|

This works for corporate buildings and for traditional concentrator-based
VPNs, where all VPN clients should use the Domain Controller closest to the
concentrator. It does not work well for a Mesh VPN.

### Overview of this tool
This is a CLI tool to run on a Tailscale-connected client which can see all
of the other connected clients. If ACLs block visibility between most Users,
it will need to run on an administrative node which does have visibility to
the other devices.

The CLI tool takes a set of command-line arguments to map Tailscale
[DERP Relay locations](https://tailscale.com/blog/how-tailscale-works/#encrypted-tcp-relays-derp)
to the nearest Active Director Comain Controller. For example in a
set of Domain Controllers in Atlanta, London, and Tokyo, one might use:
```
ts-AD-sites.exe --map="nyc:ATL-DC" --map="ord:ATL-DC" --map="dfw:ATL-DC" \
    --map="sin:TOK-DC" --map="syd:TOK-DC" --map="tok:TOK-DC" \
    -–map-default="LHR-DC"
```

iIt will then output a series of PowerShell commands to populate the Active
Directory `Sites` > `Subnets` table:
```
Get-ADReplicationSubnet -Filter "Location -like '*Tailscale'" | Remove-ADReplicationSubnet
New-ADReplicationSubnet -Name '100.100.101.101/32' -Location 'Tailscale' -Site TOK-DC
New-ADReplicationSubnet -Name 'fd7a:115c:a1e0:ab12:4843:cd96:6264:6565/128' -Location 'Tailscale' -Site TOK-DC
New-ADReplicationSubnet -Name '100.100.101.102/32' -Location 'Tailscale' -Site ATL-DC
New-ADReplicationSubnet -Name 'fd7a:115c:a1e0:ab12:4843:cd96:6264:6566/128' -Location 'Tailscale' -Site ATL-DC
New-ADReplicationSubnet -Name '100.100.101.103/32' -Location 'Tailscale' -Site TOK-DC
New-ADReplicationSubnet -Name 'fd7a:115c:a1e0:ab12:4843:cd96:6264:6567/128' -Location 'Tailscale' -Site TOK-DC
New-ADReplicationSubnet -Name '100.100.101.104/32' -Location 'Tailscale' -Site LHR-DC
New-ADReplicationSubnet -Name 'fd7a:115c:a1e0:ab12:4843:cd96:6264:6568/128' -Location 'Tailscale' -Site LHR-DC
```

Note that the tool makes no immediate changes in the configuration of Active Directory,
and indeed has no permissions to do so. It just outputs PowerShell commands to be
examined and run by a Domain Admin.

| Name | Site | Location | Type | Description
|------|------|----------|------|------------
|100.100.101.101/32|TOK-DC|Tailscale|Subnet|
|fd7a:115c:a1e0:ab12:4843:cd96:6264:6565/128|TOK-DC|Tailscale|Subnet|
|100.100.101.102/32|ATL-DC|Tailscale|Subnet|
|fd7a:115c:a1e0:ab12:4843:cd96:6264:6566/128|ATL-DC|Tailscale|Subnet|
|100.100.101.103/32|TOK-DC|Tailscale|Subnet|
|fd7a:115c:a1e0:ab12:4843:cd96:6264:6567/128|TOK-DC|Tailscale|Subnet|
|100.100.101.104/32|LHR-DC|Tailscale|Subnet|
|fd7a:115c:a1e0:ab12:4843:cd96:6264:6568/128|LHR-DC|Tailscale|Subnet|
