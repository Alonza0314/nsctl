#!/bin/bash

diff() {
    if [ "$2" != "$3" ]; then
        echo "[-][FAIL] $1 test"
        echo "Expected:"
        echo "$3"
        echo "Got:"
        echo "$2"
        exit 1
    fi

    echo "[+][PASS] $1 test"
}

setup_two_namespace() {
    ../nsctl ns create test-1
    ../nsctl ns create test-2
    ../nsctl ns list
    echo
}

cleanup_two_namespace() {
    echo
    ../nsctl ns delete test-1
    ../nsctl ns delete test-2
}

list_test_1_lo() {
    output=$(../nsctl net list test-1)
    expect=$(cat list_test_1_lo.txt)

    diff "List test 1 lo" "$output" "$expect"
}

connect_test() {
    output=$(../nsctl net connect test-1 test-2)
    expect=$(cat connect_test.txt)

    diff "Connect test" "$output" "$expect"
}

connect_test_error() {
    output=$(../nsctl net connect test-1 test-2)
    expect=$(cat connect_test_error.txt)

    diff "Connect test error" "$output" "$expect"
}

disconnect_test() {
    output=$(../nsctl net disconnect test-1 test-2)
    expect=$(cat disconnect_test.txt)

    diff "Disconnect test" "$output" "$expect"
}

main() {
    setup_two_namespace

    list_test_1_lo
    connect_test
    connect_test_error
    disconnect_test
    list_test_1_lo

    cleanup_two_namespace
}

main "$@"