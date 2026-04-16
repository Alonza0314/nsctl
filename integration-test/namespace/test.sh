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

empty_list() {
    output=$(../nsctl ns list)
    expect=$(cat empty_list.txt)

    diff "Empty list" "$output" "$expect"
}

create_test-1() {
    output=$(../nsctl ns create test-1)
    expect=$(cat create_test-1.txt)

    diff "Create test-1" "$output" "$expect"
}

list_test-1() {
    output=$(../nsctl ns list)
    expect=$(cat list_test-1.txt)

    diff "List test-1" "$output" "$expect"
}

create_test-2() {
    output=$(../nsctl ns create test-2)
    expect=$(cat create_test-2.txt)

    diff "Create test-2" "$output" "$expect"
}

list_test-1-2() {
    output=$(../nsctl ns list)
    expect=$(cat list_test-1-2.txt)

    diff "List test-1-2" "$output" "$expect"
}

delete_test-1() {
    output=$(../nsctl ns delete test-1)
    expect=$(cat delete_test-1.txt)

    diff "Delete test-1" "$output" "$expect"
}

delete_test-1_error() {
    output=$(../nsctl ns delete test-1)
    expect=$(cat delete_test-1_error.txt)

    diff "Delete test-1_error" "$output" "$expect"
}

create_test-2_error() {
    output=$(../nsctl ns create test-2)
    expect=$(cat create_test-2_error.txt)

    diff "Create test-2_error" "$output" "$expect"
}

list_test-2() {
    output=$(../nsctl ns list)
    expect=$(cat list_test-2.txt)

    diff "List test-2" "$output" "$expect"
}

delete_test-2() {
    output=$(../nsctl ns delete test-2)
    expect=$(cat delete_test-2.txt)

    diff "Delete test-2" "$output" "$expect"
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