## Issues
### UDP buffer size
If you have such error message: `sys_conn.go:62: failed to sufficiently increase receive buffer size (was: 208 kiB, wanted: 7168 kiB, got: 416 kiB). See https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes for details.`. Follow the instruction to increase the UDP receive buffer size.
```bash
sudo sysctl -w net.core.rmem_max=7500000
sudo sysctl -w net.core.wmem_max=7500000
sudo sysctl -p
```