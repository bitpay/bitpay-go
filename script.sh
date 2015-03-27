#!/bin/bash
export PAIRINGCODE=$(python helpers/pair_steps.py)
echo $PAIRINGCODE
export INVOICEPAIR=$(python helpers/pair_steps.py)
echo $INVOICEPAIR
export RETRIEVEPAIR=$(python helpers/pair_steps.py)
echo $RETRIEVEPAIR
