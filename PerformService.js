#!/usr/bin/osascript -l JavaScript

/*
	Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
	MIT Licence applies http://opensource.org/licenses/MIT
	Created on 2020-08-01

	Runs macOS service on contents of general pasteboard.

	Name of service is passed as $1
*/

ObjC.import('Cocoa')

function run(argv) {
	let service = argv[0]
	console.log(`performing service "${service}" ...`)
	let ret = $.NSPerformService(service, $.NSPasteboard.generalPasteboard)
	if (!ret) return `Service “${service}” failed`
}