# retchecking

An individual project at mercari summer internship 2020.

メルカリサマーインターン2020での成果物です。

This project is inspired by [DQNEO](https://twitter.com/DQNEO)'s idea, who is an extreamely talented engineer at mercari inc.

このプロジェクトの大元のアイデアはメルカリのエンジニアである[DQNEO](https://twitter.com/DQNEO)氏によるものです。

## Description

This program checks whether `return` statement is used immediately after certain functions are called in the given package.

パッケージの中の特定の関数について、その関数を呼び出した後、すぐにreturnしているかを検証するプログラムです。

## Install

`go get -u github.com/BOBO1997/retchecking`

## Usage

1. Add the `types.Object` of the function you want to check, into the variable `funcSetObj`.
2. Put the `TOBETESTED.go` file to be checked into `testdata/src/TOBETESTED/` directory and run `go test`.

Or you can use `go vet` command as followinig.

`go vet -vettool=$(which retcheking) TOBETESTED.go`

## Notes

The current source code is a prototype and is not so efficient and not covering all the situation.
THis tool uses the AST and TYPES information, but it seems that AST itself cannot detect whole possible cases.
I am now changing the technical relaization of the core part from AST to SSA.
