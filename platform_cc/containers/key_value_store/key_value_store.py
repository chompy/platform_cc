#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function
import sys
import os
import json
import argparse
import hashlib

DEFAULT_DATA_STORE_PATH = "/data"
INDEX_FILENAME = ".index"
HASH_KEY = "V1aKkEZbPCguWR4IVrgU5Ln59xBU4mIJSYD7ezQo8H##412%"

def _keyPath(key):
    keyHash = hashlib.sha256(
        HASH_KEY + key
    ).hexdigest()
    return os.path.join(
        DEFAULT_DATA_STORE_PATH,
        keyHash
    )

def listIndex():
    indexPath = os.path.join(
        DEFAULT_DATA_STORE_PATH,
        INDEX_FILENAME
    )
    if not os.path.isfile(indexPath):
        return []
    with open(indexPath, "r") as f:
        return json.load(f)

def updateIndex(key, action = "add"):
    indexPath = os.path.join(
        DEFAULT_DATA_STORE_PATH,
        INDEX_FILENAME
    )
    index = listIndex()
    if key not in index and action == "add":
        index.append(key)
    if key in index and action == "delete":
        index.remove(key)
    with open(indexPath, "w") as f:
        f.write(json.dumps(index))

def get(key):
    """ Get value of key. """
    keyPath = _keyPath(key)
    if not os.path.isfile(keyPath):
        return ""
    with open(keyPath, "r") as f:
        return f.read()

def set(key, value):
    """ Set value of key. """
    keyPath = _keyPath(key)
    with open(keyPath, "w") as f:
        f.write(value)
    updateIndex(key, "add")

def delete(key):
    """ Delete key. """
    keyPath = _keyPath(key)
    if not os.path.isfile(keyPath):
        return
    os.remove(keyPath)
    updateIndex(key, "delete")

def list():
    """ List all key/value pairs """
    output = {}
    for key in listIndex():
        output[key] = get(key)
    return output

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="Manage key/value store.")
    parser.add_argument(
        "action",
        type=str,
        help="Action to perform. (get, set, delete, list)"
    )
    parser.add_argument(
        "-k", "--key",
        type=str,
        help="Key to perform action against."
    )
    parser.add_argument(
        "-v", "--value",
        type=str,
        help="Value to store."
    )
    args = parser.parse_args()

    # LIST COMMAND
    if args.action == "list":
        print(json.dumps(list()))
        sys.exit()

    # for all other commands a key is required
    if not args.key:
        sys.exit("Key not provided.")

    # GET COMMAND
    if args.action == "get":
        print(get(args.key))

    # DELETE COMMAND
    elif args.action == "delete":
        delete(args.key)

    # SET COMMAND
    elif args.action == "set":
        set(args.key, args.value)