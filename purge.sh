#!/bin/bash
# remove all block files and EC, RO, PO key

rm -rv BLOCK
rm -rv internal/role/EC/*
rm -rv internal/role/RO/*
rm -rv internal/role/PO/*
rm -rv test/levelDB

