# BruteDirBuster
Golang implementation of brute force version of OWASP's DirBuster. I took the idea (and directory list) initially from [1], 
but decided to implement in Go because of its nice concurrency features.

The purpose of this program is to find hidden directories leading from a URL. It will write successful GETs to file (and print
in verbose mode). This is purely educational / experimental code. It has the potential to DoS a small server so please use carefully.

#OPTIONS / FLAGS
  
    - URL. The URL from which to launch the search.
    - FNAME. Path to file of strings that will be appended to the URL.
    - V. Verbose Mode, Prints 200 and 401 codes to the screen.
    - TOR. Make requests through Tor (or not).

#TODO:

	# Add Http Authentication.
	# Add Cookie Authentication.

- Tested on OSX and Ubuntu




[1] https://github.com/NoobieDog/Dir-Xcan
