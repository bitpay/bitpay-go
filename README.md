# BitPay Library for Go 
[![](https://travis-ci.org/bitpay/bitpay-go.svg?branch=master)](http://travis-ci.org/bitpay/bitpay-go)

Powerful, flexible, lightweight interface to the BitPay Bitcoin Payment Gateway API.

## [Getting Started &raquo;](http://dev.bitpay.com/guides/go.html)

Code documentation is available on [godoc](http://godoc.org/github.com/bitpay/bitpay-go)
## API Documentation

API Documentation is available on the [BitPay site](https://bitpay.com/api).

## Running the Tests

The reference project is at https://github.com/bitpay/bitpay-go-cli. You will need a working go installation to follow these instructions.

In order to run the tests, follow these steps:

1. Clone the repository

    `git clone https://github.com/bitpay/bitpay-go-cli.git`

1. Set the $GOPATH and $PATH variables

    `source helpers/enviro.sh`

1. Set your test api url (such as https://test.bitpay.com) and your username and password.

    `source helpers/set_constants.sh <url> <username> <password>` 
    
    `python helpers/pair_steps.py`
    
   The python script should retrieve three pairing codes from the server and store them in three files in a `temp` directory in the main project directory, `temp/retrievecode.txt`, `temp/paircode.txt`, and `temp/invoicecode.txt`. If this does not go smoothly, you can manually add the pairing codes to those files by visiting (https://test.bitpay.com/dashboard/merchant/api-tokens) and creating three tokens, saving each pairing code into one of the files in temp.

1. For reasons that are not entirely clear, we need to delete all of the required files and re-import them.
  
    `rm -rf src/github.com src/golang`

    `go get -u -t github.com/bitpay/bitpay-go/client`

1. We are now ready to run the tests.
  
  `ginkgo -r src/github.com/bitpay/`
 
## Found a bug?
Let us know! Send a pull request or a patch. Questions? Ask! We're here to help. We will respond to all filed issues.

**BitPay Support:**

* [BitPay Labs](https://labs.bitpay.com/c/libraries/python)
  * Post a question in our discussion forums
* [GitHub Issues](https://github.com/bitpay/bitpay-python/issues)
  * Open an issue if you are having issues with this library
* [Support](https://support.bitpay.com)
  * BitPay merchant support documentation

Sometimes a download can become corrupted for various reasons.  However, you can verify that the release package you downloaded is correct by checking the md5 checksum "fingerprint" of your download against the md5 checksum value shown on the Releases page.  Even the smallest change in the downloaded release package will cause a different value to be shown!
  * If you are using Windows, you can download a checksum verifier tool and instructions directly from Microsoft here: http://www.microsoft.com/en-us/download/details.aspx?id=11533
  * If you are using Linux or OS X, you already have the software installed on your system.
    * On Linux systems use the md5sum program.  For example:
      * md5sum filename
    * On OS X use the md5 program.  For example:
      * md5 filename
