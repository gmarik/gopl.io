package main

// Exercise 3.13: Write const declarations for KB, MB, up through YB as compactly
// as you can

// 1 KB		= 1000 bytes = 1000^n [n=1]
// 1 KiB	= 1024 bytes = 2^(10*n) [n=1]
// 1 GB 	= 1000^n = 1000 ^ 3 = 1 000 000 000 (= 976 562.2 GiB)
// 1 GiB	= 2^(10 * 3) = 2^30 = 1 073 741 824

const (
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB
	EB = 1000 * PB
	ZB = 1000 * EB
	YB = 1000 * ZB
)
