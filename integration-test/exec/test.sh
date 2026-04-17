#!/bin/bash

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

setup_two_namespace() {
    ../nsctl ns create test-1
    ../nsctl ns create test-2
    ../nsctl ns list
}

cleanup_two_namespace() {
    ../nsctl ns delete test-1
    ../nsctl ns delete test-2
}

connect_two_namespace() {
    ../nsctl net connect test-1 test-2
}

disconnect_two_namespace() {
    ../nsctl net disconnect test-1 test-2
}

set_ip() {
    ../nsctl net set-ip test-1 test-1-test-2 10.0.0.1/24
    ../nsctl net set-ip test-2 test-2-test-1 10.0.0.2/24
}

if_up() {
    ../nsctl net up test-1 test-1-test-2
    ../nsctl net up test-2 test-2-test-1
    echo
}

if_down() {
    echo
    ../nsctl net down test-1 test-1-test-2
    ../nsctl net down test-2 test-2-test-1
}

exec_ping_test() {
    output=$(../nsctl exec test-1 -- ping -I test-1-test-2 -c 3 10.0.0.2)

    assert_match "Exec ping header" "$output" '^PING 10\.0\.0\.2 \(10\.0\.0\.2\) from 10\.0\.0\.1 test-1-test-2: 56\(84\) bytes of data\.$'
    assert_count "Exec ping reply count" "$output" '^64 bytes from 10\.0\.0\.2: icmp_seq=[1-3] ttl=64 time=[0-9.]+ ms$' 3
    assert_match "Exec ping packet loss" "$output" '^3 packets transmitted, 3 received, 0% packet loss, time [0-9]+ms$'
    assert_match "Exec ping rtt stats" "$output" '^rtt min/avg/max/mdev = [0-9.]+/[0-9.]+/[0-9.]+/[0-9.]+ ms$'
}

main() {
    setup_two_namespace
    connect_two_namespace
    set_ip
    if_up

    exec_ping_test

    if_down
    disconnect_two_namespace
    cleanup_two_namespace
}

main "$@"