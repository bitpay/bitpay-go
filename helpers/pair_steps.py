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
    browser.find_by_css(".token-access-new-button").find_by_css(".btn").find_by_css(".icon-plus")[0].click()
    browser.find_by_id("token-new-form").find_by_css(".btn")[0].click()
    code = browser.find_by_css(".token-claimcode")[0].html
    gopath = os.environ['GOPATH']
    tempath = gopath + "/temp"
    if not os.path.exists(tempath):
        os.makedirs(tempath)
    f = open(tempath + "/retrievecode.txt", 'w')
    f.write(code)
    f.close()
    time.sleep(10)
    browser.reload()
    browser.find_by_css(".token-access-new-button").find_by_css(".btn").find_by_css(".icon-plus")[0].click()
    browser.find_by_id("token-new-form").find_by_css(".btn")[0].click()
    code = browser.find_by_css(".token-claimcode")[0].html
    gopath = os.environ['GOPATH']
    tempath = gopath + "/temp"
    if not os.path.exists(tempath):
        os.makedirs(tempath)
    f = open(tempath + "/paircode.txt", 'w')
    f.write(code)
    f.close()
    time.sleep(10)
    browser.reload()
    browser.find_by_css(".token-access-new-button").find_by_css(".btn").find_by_css(".icon-plus")[0].click()
    browser.find_by_id("token-new-form").find_by_css(".btn")[0].click()
    code = browser.find_by_css(".token-claimcode")[0].html
    gopath = os.environ['GOPATH']
    tempath = gopath + "/temp"
    if not os.path.exists(tempath):
        os.makedirs(tempath)
    f = open(tempath + "/invoicecode.txt", 'w')
    f.write(code)
    f.close()
    return code

code = get_claim_code_from_server()
sys.stdout.write(code)
