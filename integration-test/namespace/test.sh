#!/bin/bash

empty_list() {
    output=$(../nsctl ns list)
    expect=$(cat empty_list.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] Empty list test failed"
        echo "Expected:"
        echo "$expected"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] Empty list test passed"
}

create_test-1() {
    output=$(../nsctl ns create test-1)
    expect=$(cat create_test-1.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] Create test-1 failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] Create test-1 passed"
}

list_test-1() {
    output=$(../nsctl ns list)
    expect=$(cat list_test-1.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] List test-1 failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] List test-1 passed"
}

create_test-2() {
    output=$(../nsctl ns create test-2)
    expect=$(cat create_test-2.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] Create test-2 failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] Create test-2 passed"
}

list_test-1-2() {
    output=$(../nsctl ns list)
    expect=$(cat list_test-1-2.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] List test-1-2 failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] List test-1-2 passed"
}

delete_test-1() {
    output=$(../nsctl ns delete test-1)
    expect=$(cat delete_test-1.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] Delete test-1 failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] Delete test-1 passed"
}

delete_test-1_error() {
    output=$(../nsctl ns delete test-1)
    expect=$(cat delete_test-1_error.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] Delete test-1_error failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] Delete test-1_error passed"
}

create_test-2_error() {
    output=$(../nsctl ns create test-2)
    expect=$(cat create_test-2_error.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] Create test-2_error failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] Create test-2_error passed"
}

list_test-2() {
    output=$(../nsctl ns list)
    expect=$(cat list_test-2.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] List test-2 failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] List test-2 passed"
}

delete_test-2() {
    output=$(../nsctl ns delete test-2)
    expect=$(cat delete_test-2.txt)

    if [ "$output" != "$expect" ]; then
        echo "[-] Delete test-2 failed"
        echo "Expected:"
        echo "$expect"
        echo "Got:"
        echo "$output"
        exit 1
    fi

    echo "[+] Delete test-2 passed"
}

main() {
    empty_list
    create_test-1
    list_test-1
    create_test-2
    list_test-1-2
    delete_test-1
    delete_test-1_error
    create_test-2_error
    list_test-2
    delete_test-2
    empty_list
}

main "$@"