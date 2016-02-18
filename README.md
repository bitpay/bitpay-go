# BitPay Library for Go
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/bitpay/bitpay-go/master/LICENSE)
[![Travis Build](https://img.shields.io/travis/bitpay/bitpay-go.svg?style=flat-square)](http://travis-ci.org/bitpay/bitpay-go)


Powerful, flexible, lightweight interface to the BitPay Bitcoin Payment Gateway API.

## [Getting Started &raquo;](https://github.com/bitpay/bitpay-go/blob/master/GUIDE.md)

Code documentation is available on [godoc](http://godoc.org/github.com/bitpay/bitpay-go).
## API Documentation

API Documentation is available on the [BitPay site](https://bitpay.com/api).

## Running the Tests

In order to run the tests, follow these steps:

1. Set the $GOPATH and $PATH variables

1. Install the dependencies

 ```bash
$ go get github.com/btcsuite/btcutil
$ go get github.com/gorilla/mux
$ go get github.com/onsi/ginkgo
$ go get golang.org/x/crypto
```

1. Clone the repository

    `git clone https://github.com/bitpay/bitpay-go.git`

    Into src/github.com/bitpay/bitpay-go/

1. Set the environment variables `BITPAYAPI` & `BITPAYPEM` to "https://test.bitpay.com" and a valid PEM value.

    This is slightly tricky, the PEM file has to already be paired with a merchant token on your bitpay account. To do this it is probably best to use the [bitpay test helper](https://github.com/bitpay/bitpay-test-helper).

1. You will also need a paid invoice on the server. Set the environment variable `INVOICEID` to the id of a paid invoice on the server.

1. We are now ready to run the tests.

  `ginkgo -r src/github.com/bitpay/`

## Found a bug?
Let us know! Send a pull request or a patch. Questions? Ask! We're here to help. We will respond to all filed issues.

**BitPay Support:**

* [GitHub Issues](https://github.com/bitpay/bitpay-python/issues)
  * Open an issue if you are having issues with this library
* [Support](https://help.bitpay.com)
  * BitPay merchant support documentation

Sometimes a download can become corrupted for various reasons.  However, you can verify that the release package you downloaded is correct by checking the md5 checksum "fingerprint" of your download against the md5 checksum value shown on the Releases page.  Even the smallest change in the downloaded release package will cause a different value to be shown!
  * If you are using Windows, you can download a checksum verifier tool and instructions directly from Microsoft here: http://www.microsoft.com/en-us/download/details.aspx?id=11533
  * If you are using Linux or OS X, you already have the software installed on your system.
    * On Linux systems use the md5sum program.  For example:
      * md5sum filename
    * On OS X use the md5 program.  For example:
      * md5 filename
