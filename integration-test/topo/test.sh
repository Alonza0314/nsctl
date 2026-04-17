#!/bin/bash

diff() {
    if [ "$2" != "$3" ]; then
        echo "[-][FAIL] $1"
        echo "Expected:"
        echo "$3"
        echo "Got:"
        echo "$2"
        exit 1
    fi

    echo "[+][PASS] $1"
}

assert_match() {
    if ! printf '%s\n' "$2" | grep -Eq "$3"; then
        echo "[-][FAIL] $1"
        echo "Expected pattern:"
        echo "$3"
        echo "Got:"
        echo "$2"
        exit 1
    fi

    echo "[+][PASS] $1"
}

assert_count() {
    count=$(printf '%s\n' "$2" | grep -Ec "$3")
    if [ "$count" -ne "$4" ]; then
        echo "[-][FAIL] $1"
        echo "Expected count: $4"
        echo "Pattern:"
        echo "$3"
        echo "Got count: $count"
        echo "Output:"
        echo "$2"
        exit 1
    fi

    echo "[+][PASS] $1"
}

topo_apply_cycle() {
    output=$(../nsctl topo apply topo_template_deps_cycle.yaml)
    expect=$(cat topo_apply_cycle.txt)

    diff "Topo apply cycle" "$output" "$expect"
}

topo_apply() {
    output=$(../nsctl topo apply topo_template.yaml)
    expect=$(cat topo_apply.txt)

    diff "Topo apply" "$output" "$expect"
}

exec_ping_test() {
    output=$(../nsctl exec test-1 -- ping -I br-test-1 -c 3 10.0.0.2)

    assert_match "Exec ping header" "$output" '^PING 10\.0\.0\.2 \(10\.0\.0\.2\) from 10\.0\.0\.1 br-test-1: 56\(84\) bytes of data\.$'
    assert_count "Exec ping reply count" "$output" '^64 bytes from 10\.0\.0\.2: icmp_seq=[1-3] ttl=64 time=[0-9.]+ ms$' 3
    assert_match "Exec ping packet loss" "$output" '^3 packets transmitted, 3 received, 0% packet loss, time [0-9]+ms$'
    assert_match "Exec ping rtt stats" "$output" '^rtt min/avg/max/mdev = [0-9.]+/[0-9.]+/[0-9.]+/[0-9.]+ ms$'
}

topo_delete() {
    output=$(../nsctl topo delete topo_template.yaml)
    expect=$(cat topo_delete.txt)

    diff "Topo delete" "$output" "$expect"
}

main() {
    topo_apply_cycle
    topo_apply
    exec_ping_test
    topo_delete
}

main "$@"