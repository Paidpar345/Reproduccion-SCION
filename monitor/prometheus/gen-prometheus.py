#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Genera monitor/prometheus/prometheus.yml con etiquetas por nodo.

A diferencia del prometheus.yml estático (instance = IP:puerto), este generador
añade etiquetas node/isd/as/role a cada target, lo que permite filtrar en Grafana
por nodo o por ISD. Sin dependencias externas: usa la convención del testbed
(4 carpetas ISD x 5 AS, IP = 10.100.0.<dir><as>).

Uso:
    python monitor/prometheus/gen-prometheus.py > monitor/prometheus/prometheus.yml
"""
# carpeta ISDk  ->  identificador SCION real
ISD = {1: 16, 2: 17, 3: 18, 4: 19}
GATEWAYS = [("endhost-as15", "10.100.0.115"), ("endhost-as35", "10.100.0.135")]


def job(name, port, role):
    out = ["  - job_name: %s" % name, "    metrics_path: /metrics",
           "    static_configs:"]
    for d in (1, 2, 3, 4):
        for a in (1, 2, 3, 4, 5):
            ip = "10.100.0.%d%d" % (d, a)
            out.append("      - targets: ['%s:%d']" % (ip, port))
            out.append("        labels: {node: scion%d%d, isd: '%d', as: '%d%d', role: %s}"
                       % (d, a, ISD[d], d, a, role))
    return "\n".join(out)


def main():
    print("global:\n  scrape_interval: 5s\n  evaluation_interval: 15s\n")
    print("scrape_configs:")
    print(job("scion-cs", 30452, "control"))
    print(job("scion-br", 30442, "border-router"))
    print("  - job_name: scion-gateway\n    metrics_path: /metrics\n    static_configs:")
    for name, ip in GATEWAYS:
        print("      - targets: ['%s:30456']" % ip)
        print("        labels: {node: %s, role: gateway}" % name)
    print("  - job_name: prometheus\n    static_configs:\n      - targets: ['localhost:9090']")


if __name__ == "__main__":
    main()
