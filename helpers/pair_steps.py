#! /usr/local/bin/python

import sys
import os
from splinter import Browser
import time

ROOT_ADDRESS = os.environ['RCROOTADDRESS']
USER_NAME = os.environ['RCTESTUSER']
PASSWORD = os.environ['RCTESTPASSWORD']

def get_claim_code_from_server():
    time.sleep(5)
    browser = Browser('phantomjs', service_args=['--ignore-ssl-errors=true'])
    browser.visit(ROOT_ADDRESS + "/merchant-login")
    browser.fill_form({"email": USER_NAME, "password": PASSWORD})
    browser.find_by_id("loginButton")[0].click()
    time.sleep(1)
    browser.visit(ROOT_ADDRESS + "/api-tokens")
    code = get_code_from_page(browser, "thiscodewillneverbevalid")
    gopath = os.environ['GOPATH']
    tempath = gopath + "/temp"
    if not os.path.exists(tempath):
        os.makedirs(tempath)
    write_code_to_file(code, tempath + "/retrievecode.txt")
    print(code)
    time.sleep(10)
    browser.reload()
    code = get_code_from_page(browser, code)
    write_code_to_file(code, tempath + "/paircode.txt")
    print(code)
    time.sleep(10)
    browser.reload()
    code = get_code_from_page(browser, code)
    write_code_to_file(code, tempath + "/invoicecode.txt")
    print(code)
    return code

def get_code_from_page(browser, code):
    browser.find_by_css(".token-access-new-button").find_by_css(".btn").find_by_css(".icon-plus")[0].click()
    browser.find_by_id("token-new-form").find_by_css(".btn")[0].click()
    browser.reload()
    newcode = browser.find_by_css(".token-claimcode")[0].html
    return newcode

def write_code_to_file(code, fname):
    f = open(fname, 'w')
    f.write(code)
    f.close()

get_claim_code_from_server()
print("done")
